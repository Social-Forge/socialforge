package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/centrifugo"
	"github/socialforge/internal/infra/channels"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/rabbitmq"
	"github/socialforge/internal/infra/repository"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// DispatchJob is the queue payload for delivering an outbound message.
type DispatchJob struct {
	MessageID string `json:"message_id"`
	TenantID  string `json:"tenant_id"`
}

type OutboundService struct {
	conversationRepo repository.ConversationRepository
	messageRepo      repository.MessageRepository
	outboxRepo       repository.MessageOutboxRepository
	channelRepo      repository.ChannelRepository
	contactRepo      repository.ContactRepository
	centrifugo       *centrifugo.CentrifugoClient
	rabbit           *rabbitmq.Client
	senders          map[string]channels.Sender
	logger           *zap.Logger
}

func NewOutboundService(
	conversationRepo repository.ConversationRepository,
	messageRepo repository.MessageRepository,
	outboxRepo repository.MessageOutboxRepository,
	channelRepo repository.ChannelRepository,
	contactRepo repository.ContactRepository,
	centrifugoClient *centrifugo.CentrifugoClient,
	rabbit *rabbitmq.Client,
	senders map[string]channels.Sender,
	logger *zap.Logger,
) *OutboundService {
	return &OutboundService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		outboxRepo:       outboxRepo,
		channelRepo:      channelRepo,
		contactRepo:      contactRepo,
		centrifugo:       centrifugoClient,
		rabbit:           rabbit,
		senders:          senders,
		logger:           logger,
	}
}

// SendText persists an outbound agent text message, records it in the outbox,
// publishes it to the conversation in realtime (optimistic), and enqueues it
// for provider delivery by the worker.
func (s *OutboundService) SendText(ctx context.Context, tenantID, conversationID, agentUserID, text string) (*entity.Message, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id: %w", err)
	}
	cid, err := uuid.Parse(conversationID)
	if err != nil {
		return nil, fmt.Errorf("invalid conversation id: %w", err)
	}
	var senderID *uuid.UUID
	if agentUserID != "" {
		if aid, err := uuid.Parse(agentUserID); err == nil {
			senderID = &aid
		}
	}

	msg := &entity.Message{
		TenantID:       tid,
		ConversationID: cid,
		SenderID:       senderID,
		Direction:      entity.MessageDirectionOut,
		SenderType:     entity.SenderTypeAgent,
		ContentType:    entity.ContentTypeText,
		Body:           entity.NewNullString(text),
		Status:         entity.MessageStatusPending,
	}

	tctx := repository.WithTenantID(subCtx, tid)
	err = s.conversationRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		conv, err := s.conversationRepo.FindByID(txCtx, cid)
		if err != nil {
			return err
		}
		if _, err := s.messageRepo.Create(txCtx, msg); err != nil {
			return err
		}
		if err := s.outboxRepo.Create(txCtx, &entity.MessageOutbox{TenantID: tid, MessageID: msg.ID}); err != nil {
			return err
		}
		_ = conv
		return s.conversationRepo.TouchOutbound(txCtx, cid, time.Now())
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Optimistic realtime: show the message in the conversation immediately.
	if s.centrifugo != nil {
		_ = s.centrifugo.BroadcastNewMessage(subCtx, cid.String(), map[string]interface{}{
			"id":              msg.ID.String(),
			"conversation_id": cid.String(),
			"direction":       msg.Direction,
			"content_type":    msg.ContentType,
			"body":            text,
			"status":          msg.Status,
			"created_at":      msg.CreatedAt,
		})
	}

	// Enqueue for provider delivery.
	if s.rabbit != nil {
		data, _ := json.Marshal(DispatchJob{MessageID: msg.ID.String(), TenantID: tenantID})
		if err := s.rabbit.Publish(subCtx, rabbitmq.QueueDispatchOutbound, data); err != nil {
			s.logger.Warn("failed to enqueue dispatch", zap.Error(err))
		}
	}
	return msg, nil
}

// ProcessDispatch is the worker consumer for outbound delivery: it resolves the
// channel + recipient and sends via the provider adapter, then records the
// delivery outcome. Send failures are recorded (status=failed) without requeue;
// the outbox row retains attempt state for a future retry sweep.
func (s *OutboundService) ProcessDispatch(ctx context.Context, jobBody []byte) error {
	var job DispatchJob
	if err := json.Unmarshal(jobBody, &job); err != nil {
		return fmt.Errorf("invalid dispatch job: %w", err)
	}
	tid, err := uuid.Parse(job.TenantID)
	if err != nil {
		return fmt.Errorf("invalid tenant id: %w", err)
	}
	mid, err := uuid.Parse(job.MessageID)
	if err != nil {
		return fmt.Errorf("invalid message id: %w", err)
	}

	tctx := repository.WithTenantID(ctx, tid)
	var msg *entity.Message
	var channel *entity.Channel
	var contact *entity.Contact
	err = s.messageRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		var err error
		if msg, err = s.messageRepo.FindByID(txCtx, mid); err != nil {
			return err
		}
		conv, err := s.conversationRepo.FindByID(txCtx, msg.ConversationID)
		if err != nil {
			return err
		}
		if channel, err = s.channelRepo.FindByID(txCtx, conv.ChannelID); err != nil {
			return err
		}
		contact, err = s.contactRepo.FindByID(txCtx, conv.ContactID)
		return err
	})
	if err != nil {
		return err
	}

	sender := s.senders[channel.Type]
	if sender == nil {
		s.recordOutcome(ctx, tctx, mid, msg.ConversationID, entity.MessageStatusFailed, "no sender configured for "+channel.Type)
		return nil
	}

	providerID, sendErr := sender.SendText(ctx, channel, contact.ExternalID, msg.Body.String)
	if sendErr != nil {
		s.logger.Warn("provider send failed", zap.String("channel_type", channel.Type), zap.Error(sendErr))
		s.recordOutcome(ctx, tctx, mid, msg.ConversationID, entity.MessageStatusFailed, sendErr.Error())
		return nil
	}
	s.logger.Info("message delivered", zap.String("channel_type", channel.Type), zap.String("provider_message_id", providerID))
	s.recordOutcome(ctx, tctx, mid, msg.ConversationID, entity.MessageStatusSent, "")
	return nil
}

func (s *OutboundService) recordOutcome(ctx context.Context, tctx context.Context, messageID, conversationID uuid.UUID, status, errMsg string) {
	outboxStatus := entity.OutboxStatusSent
	if status == entity.MessageStatusFailed {
		outboxStatus = entity.OutboxStatusFailed
	}
	err := s.messageRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		if err := s.messageRepo.UpdateStatus(txCtx, messageID, status, errMsg); err != nil {
			return err
		}
		return s.outboxRepo.SetStatusByMessage(txCtx, messageID, outboxStatus, errMsg)
	})
	if err != nil {
		s.logger.Error("failed to record dispatch outcome", zap.Error(err))
	}
	if s.centrifugo != nil && conversationID != uuid.Nil {
		_ = s.centrifugo.BroadcastConversationUpdate(ctx, conversationID.String(), map[string]interface{}{
			"type":       "message_status",
			"message_id": messageID.String(),
			"status":     status,
		})
	}
}
