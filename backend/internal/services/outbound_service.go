package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/centrifugo"
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
	centrifugo       *centrifugo.CentrifugoClient
	rabbit           *rabbitmq.Client
	logger           *zap.Logger
}

func NewOutboundService(
	conversationRepo repository.ConversationRepository,
	messageRepo repository.MessageRepository,
	outboxRepo repository.MessageOutboxRepository,
	centrifugoClient *centrifugo.CentrifugoClient,
	rabbit *rabbitmq.Client,
	logger *zap.Logger,
) *OutboundService {
	return &OutboundService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		outboxRepo:       outboxRepo,
		centrifugo:       centrifugoClient,
		rabbit:           rabbit,
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

// ProcessDispatch is the worker consumer for outbound delivery. Fase 2D ships a
// stub sender (marks delivered); Fase 2E plugs in the real WAHA/Telegram send.
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

	// TODO(Fase 2E): resolve the channel + call the provider adapter to actually
	// send. For now we optimistically mark the message as sent.
	tctx := repository.WithTenantID(ctx, tid)
	var convID uuid.UUID
	err = s.messageRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		m, err := s.messageRepo.FindByID(txCtx, mid)
		if err != nil {
			return err
		}
		convID = m.ConversationID
		if err := s.messageRepo.UpdateStatus(txCtx, mid, entity.MessageStatusSent, ""); err != nil {
			return err
		}
		return s.outboxRepo.SetStatusByMessage(txCtx, mid, entity.OutboxStatusSent, "")
	})
	if err != nil {
		return err
	}

	if s.centrifugo != nil && convID != uuid.Nil {
		_ = s.centrifugo.BroadcastConversationUpdate(ctx, convID.String(), map[string]interface{}{
			"type":       "message_status",
			"message_id": mid.String(),
			"status":     entity.MessageStatusSent,
		})
	}
	return nil
}
