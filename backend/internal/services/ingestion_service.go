package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/centrifugo"
	"github/socialforge/internal/infra/channels"
	"github/socialforge/internal/infra/channels/telegram"
	"github/socialforge/internal/infra/channels/waha"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/rabbitmq"
	"github/socialforge/internal/infra/repository"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// providerForChannelType maps a channel type to its expected webhook provider.
func providerForChannelType(t string) string {
	switch t {
	case entity.ChannelTypeWhatsAppWaha:
		return channels.ProviderWAHA
	case entity.ChannelTypeTelegram:
		return channels.ProviderTelegram
	case entity.ChannelTypeWhatsAppMeta:
		return channels.ProviderMetaWA
	case entity.ChannelTypeMessenger:
		return channels.ProviderMessenger
	case entity.ChannelTypeInstagram:
		return channels.ProviderInstagram
	default:
		return ""
	}
}

// InboundJob is the queue payload carrying a raw provider event to the worker.
type InboundJob struct {
	Provider  string          `json:"provider"`
	ChannelID string          `json:"channel_id"`
	Body      json.RawMessage `json:"body"`
}

type IngestionService struct {
	channelRepo      repository.ChannelRepository
	contactRepo      repository.ContactRepository
	conversationRepo repository.ConversationRepository
	messageRepo      repository.MessageRepository
	webhookEventRepo repository.WebhookEventRepository
	autoResponseRepo repository.AutoResponseRepository
	outbound         *OutboundService
	centrifugo       *centrifugo.CentrifugoClient
	rabbit           *rabbitmq.Client
	logger           *zap.Logger
}

func NewIngestionService(
	channelRepo repository.ChannelRepository,
	contactRepo repository.ContactRepository,
	conversationRepo repository.ConversationRepository,
	messageRepo repository.MessageRepository,
	webhookEventRepo repository.WebhookEventRepository,
	autoResponseRepo repository.AutoResponseRepository,
	outbound *OutboundService,
	centrifugoClient *centrifugo.CentrifugoClient,
	rabbit *rabbitmq.Client,
	logger *zap.Logger,
) *IngestionService {
	return &IngestionService{
		channelRepo:      channelRepo,
		contactRepo:      contactRepo,
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		webhookEventRepo: webhookEventRepo,
		autoResponseRepo: autoResponseRepo,
		outbound:         outbound,
		centrifugo:       centrifugoClient,
		rabbit:           rabbit,
		logger:           logger,
	}
}

// VerifyAndEnqueue is called by the public webhook handler: it authenticates the
// event (fast) and enqueues it for async processing, so providers get an
// immediate 200 and heavy work happens in the worker.
func (s *IngestionService) VerifyAndEnqueue(ctx context.Context, provider, channelIDStr string, headers map[string]string, body []byte) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 10*time.Second)
	defer cancel()

	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		return fmt.Errorf("invalid channel id: %w", err)
	}
	channel, err := s.channelRepo.FindByID(subCtx, channelID)
	if err != nil {
		return fmt.Errorf("channel not found: %w", err)
	}
	if want := providerForChannelType(channel.Type); want != provider {
		return fmt.Errorf("provider %q does not match channel type %q", provider, channel.Type)
	}
	secret := ""
	if channel.WebhookSecret.Valid {
		secret = channel.WebhookSecret.String
	}
	if err := verifySignature(provider, secret, headers, body); err != nil {
		return err
	}

	data, err := json.Marshal(InboundJob{Provider: provider, ChannelID: channelIDStr, Body: body})
	if err != nil {
		return err
	}
	// If the broker is unavailable, fall back to synchronous processing so no
	// message is lost.
	if s.rabbit == nil {
		return s.ProcessInbound(subCtx, data)
	}
	if err := s.rabbit.Publish(subCtx, rabbitmq.QueueIngestInbound, data); err != nil {
		s.logger.Warn("enqueue failed, processing inline", zap.Error(err))
		return s.ProcessInbound(subCtx, data)
	}
	return nil
}

// ProcessInbound is the worker-side consumer: it decodes the job and runs the
// persistence + realtime pipeline. Idempotent via webhook_events + message
// provider_message_id dedup, so redelivery is safe.
func (s *IngestionService) ProcessInbound(ctx context.Context, jobBody []byte) error {
	var job InboundJob
	if err := json.Unmarshal(jobBody, &job); err != nil {
		return fmt.Errorf("invalid inbound job: %w", err)
	}
	channelID, err := uuid.Parse(job.ChannelID)
	if err != nil {
		return fmt.Errorf("invalid channel id in job: %w", err)
	}
	channel, err := s.channelRepo.FindByID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("channel not found: %w", err)
	}
	return s.processMessage(ctx, channel, job.Provider, channelID, []byte(job.Body))
}

// processMessage runs the persistence + realtime pipeline for one verified,
// channel-resolved inbound event.
func (s *IngestionService) processMessage(ctx context.Context, channel *entity.Channel, provider string, channelID uuid.UUID, body []byte) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 20*time.Second)
	defer cancel()

	env, err := normalize(provider, body)
	if err != nil {
		return err
	}
	if env.ProviderEventID == "" {
		env.ProviderEventID = uuid.NewString() // fallback: always dedup on something
	}

	event := &entity.WebhookEvent{
		ChannelID:       channelID,
		Provider:        provider,
		ProviderEventID: env.ProviderEventID,
		EventType:       entity.NewNullString(env.EventType),
		Payload:         body,
	}

	// Non-message events (calls, receipts, session status): dedup + light
	// handling only. Best-effort — losing one is harmless.
	if !env.IsPersistableMessage() {
		inserted, err := s.webhookEventRepo.TryInsert(subCtx, event)
		if err != nil {
			return err
		}
		if inserted {
			if env.Kind == channels.KindCall {
				// Auto-reject + auto-response handled by the WAHA adapter (Fase 2E).
				s.logger.Info("incoming call event", zap.String("channel_id", channelID.String()), zap.String("from", env.Contact.ExternalID))
			}
			_ = s.webhookEventRepo.MarkProcessed(subCtx, event.ID)
		}
		return nil
	}

	ts := env.Timestamp
	if ts.IsZero() {
		ts = time.Now()
	}

	var conv *entity.Convertation
	var msg *entity.Message
	var duplicate, msgInserted, contactCreated bool

	// Dedup + persistence in ONE tenant transaction: if anything fails the whole
	// unit (including the webhook_events row) rolls back, so a provider retry
	// reprocesses cleanly.
	tctx := repository.WithTenantID(subCtx, channel.TenantID)
	err = s.contactRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		inserted, err := s.webhookEventRepo.TryInsert(txCtx, event)
		if err != nil {
			return err
		}
		if !inserted {
			duplicate = true
			return nil
		}

		contact, created, err := s.contactRepo.FindOrCreate(txCtx, &entity.Contact{
			TenantID:    channel.TenantID,
			ChannelID:   channel.ID,
			ExternalID:  env.Contact.ExternalID,
			DisplayName: env.Contact.DisplayName,
			AvatarURL:   entity.NewNullString(env.Contact.AvatarURL),
		})
		if err != nil {
			return err
		}
		contactCreated = created

		conv, err = s.conversationRepo.FindOrCreateOpen(txCtx, channel.TenantID, channel.ID, contact.ID)
		if err != nil {
			return err
		}

		// Auto-assign unassigned conversations to the least-loaded agent in the
		// channel's division (equal routing). No agent -> stays unassigned.
		if conv.AssignedAgentID == nil {
			if agentID, err := s.conversationRepo.PickAgentForDivision(txCtx, channel.TenantID, channel.DivisionID); err == nil && agentID != nil {
				if err := s.conversationRepo.Assign(txCtx, conv.ID, agentID); err == nil {
					conv.AssignedAgentID = agentID
					conv.Status = entity.ConversationStatusOpen
				}
			}
		}

		msg = &entity.Message{
			TenantID:          channel.TenantID,
			ConversationID:    conv.ID,
			Direction:         entity.MessageDirectionIn,
			SenderType:        entity.SenderTypeContact,
			ContentType:       env.ContentType,
			Body:              entity.NewNullString(env.Text),
			Media:             mediaToMap(env.Media),
			ProviderMessageID: entity.NewNullString(env.ProviderMsgID),
			Status:            entity.MessageStatusDelivered,
		}
		msgInserted, err = s.messageRepo.Create(txCtx, msg)
		if err != nil {
			return err
		}
		if msgInserted {
			if err := s.conversationRepo.TouchInbound(txCtx, conv.ID, ts); err != nil {
				return err
			}
		}
		return s.webhookEventRepo.MarkProcessed(txCtx, event.ID)
	})
	if err != nil {
		return fmt.Errorf("failed to persist inbound message: %w", err)
	}
	if duplicate {
		return nil
	}

	if msgInserted {
		s.publishRealtime(subCtx, channel, conv, msg)
	}

	// First-reply auto-response: only for a brand-new contact, and only when no
	// AI agent is active on the channel (AI overrides the static auto-response).
	if contactCreated && msgInserted && channel.AIAgentID == nil {
		s.maybeAutoRespond(subCtx, channel, conv.ID)
	}
	return nil
}

// IngestWebchat handles an inbound message from a webchat widget visitor. The
// visitorID acts as the contact external id. Reuses the same persistence +
// realtime + auto-assign + auto-response pipeline as provider webhooks.
func (s *IngestionService) IngestWebchat(ctx context.Context, channelIDStr, visitorID, name, text string) (string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 20*time.Second)
	defer cancel()

	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		return "", fmt.Errorf("invalid channel id: %w", err)
	}
	channel, err := s.channelRepo.FindByID(subCtx, channelID)
	if err != nil {
		return "", fmt.Errorf("channel not found: %w", err)
	}
	if channel.Type != entity.ChannelTypeWebchat {
		return "", fmt.Errorf("channel is not a webchat channel")
	}
	if name == "" {
		name = "Website Visitor"
	}

	var conv *entity.Convertation
	var msg *entity.Message
	var contactCreated, msgInserted bool

	tctx := repository.WithTenantID(subCtx, channel.TenantID)
	err = s.contactRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		contact, created, err := s.contactRepo.FindOrCreate(txCtx, &entity.Contact{
			TenantID:    channel.TenantID,
			ChannelID:   channel.ID,
			ExternalID:  visitorID,
			DisplayName: name,
		})
		if err != nil {
			return err
		}
		contactCreated = created

		conv, err = s.conversationRepo.FindOrCreateOpen(txCtx, channel.TenantID, channel.ID, contact.ID)
		if err != nil {
			return err
		}
		if conv.AssignedAgentID == nil {
			if agentID, err := s.conversationRepo.PickAgentForDivision(txCtx, channel.TenantID, channel.DivisionID); err == nil && agentID != nil {
				_ = s.conversationRepo.Assign(txCtx, conv.ID, agentID)
			}
		}

		msg = &entity.Message{
			TenantID:       channel.TenantID,
			ConversationID: conv.ID,
			Direction:      entity.MessageDirectionIn,
			SenderType:     entity.SenderTypeContact,
			ContentType:    entity.ContentTypeText,
			Body:           entity.NewNullString(text),
			Status:         entity.MessageStatusDelivered,
		}
		msgInserted, err = s.messageRepo.Create(txCtx, msg)
		if err != nil {
			return err
		}
		if msgInserted {
			return s.conversationRepo.TouchInbound(txCtx, conv.ID, time.Now())
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to ingest webchat message: %w", err)
	}

	if msgInserted {
		s.publishRealtime(subCtx, channel, conv, msg)
	}
	if contactCreated && msgInserted && channel.AIAgentID == nil {
		s.maybeAutoRespond(subCtx, channel, conv.ID)
	}
	return conv.ID.String(), nil
}

func (s *IngestionService) maybeAutoRespond(ctx context.Context, channel *entity.Channel, conversationID uuid.UUID) {
	if s.autoResponseRepo == nil || s.outbound == nil {
		return
	}
	var ar *entity.AutoResponse
	err := s.autoResponseRepo.RunInTenantTx(repository.WithTenantID(ctx, channel.TenantID), func(txCtx context.Context) error {
		var err error
		ar, err = s.autoResponseRepo.GetByChannel(txCtx, channel.ID)
		return err
	})
	if err != nil || ar == nil || !ar.IsEnabled || !ar.Body.Valid || ar.Body.String == "" {
		return
	}
	if _, err := s.outbound.SendSystemText(ctx, channel.TenantID, conversationID, ar.Body.String); err != nil {
		s.logger.Warn("failed to send auto-response", zap.Error(err))
	}
}

func (s *IngestionService) publishRealtime(ctx context.Context, channel *entity.Channel, conv *entity.Convertation, msg *entity.Message) {
	if s.centrifugo == nil {
		return
	}
	payload := map[string]interface{}{
		"id":              msg.ID.String(),
		"conversation_id": conv.ID.String(),
		"channel_id":      channel.ID.String(),
		"direction":       msg.Direction,
		"content_type":    msg.ContentType,
		"body":            msg.Body.String,
		"created_at":      msg.CreatedAt,
	}
	if err := s.centrifugo.BroadcastNewMessage(ctx, conv.ID.String(), payload); err != nil {
		s.logger.Warn("failed to broadcast new message", zap.Error(err))
	}
	if err := s.centrifugo.PublishToTenant(ctx, channel.TenantID.String(), map[string]interface{}{
		"type":            "conversation_updated",
		"conversation_id": conv.ID.String(),
		"channel_id":      channel.ID.String(),
		"last_message_at": msg.CreatedAt,
		"preview":         msg.Body.String,
	}); err != nil {
		s.logger.Warn("failed to publish tenant update", zap.Error(err))
	}
}

func normalize(provider string, body []byte) (*channels.Envelope, error) {
	switch provider {
	case channels.ProviderWAHA:
		return waha.Normalize(body)
	case channels.ProviderTelegram:
		return telegram.Normalize(body)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

func mediaToMap(m *channels.Media) map[string]interface{} {
	if m == nil {
		return nil
	}
	b, _ := json.Marshal(m)
	var out map[string]interface{}
	_ = json.Unmarshal(b, &out)
	return out
}

// verifySignature validates the webhook authenticity using the channel's
// webhook_secret. Providers differ: Telegram sends a secret-token header, WAHA
// signs the body with HMAC-SHA256. A plain X-Webhook-Secret header is also
// accepted (useful for local testing and simple setups).
func verifySignature(provider, secret string, headers map[string]string, body []byte) error {
	if secret == "" {
		return errors.New("channel has no webhook secret configured")
	}
	if hmacEqual(headers["X-Webhook-Secret"], secret) {
		return nil
	}
	switch provider {
	case channels.ProviderTelegram:
		if hmacEqual(headers["X-Telegram-Bot-Api-Secret-Token"], secret) {
			return nil
		}
	case channels.ProviderWAHA:
		if sig := headers["X-Webhook-Hmac"]; sig != "" {
			mac := hmac.New(sha256.New, []byte(secret))
			mac.Write(body)
			if hmacEqual(sig, hex.EncodeToString(mac.Sum(nil))) {
				return nil
			}
		}
	}
	return errors.New("invalid webhook signature")
}

func hmacEqual(a, b string) bool {
	return a != "" && hmac.Equal([]byte(a), []byte(b))
}
