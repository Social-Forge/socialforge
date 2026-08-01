package services

import (
	"context"
	"fmt"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/centrifugo"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/repository"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ConversationService struct {
	conversationRepo repository.ConversationRepository
	messageRepo      repository.MessageRepository
	centrifugo       *centrifugo.CentrifugoClient
	logger           *zap.Logger
}

func NewConversationService(
	conversationRepo repository.ConversationRepository,
	messageRepo repository.MessageRepository,
	centrifugoClient *centrifugo.CentrifugoClient,
	logger *zap.Logger,
) *ConversationService {
	return &ConversationService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		centrifugo:       centrifugoClient,
		logger:           logger,
	}
}

// mutate runs a conversation update inside the tenant tx and broadcasts the
// change to the conversation channel.
func (s *ConversationService) mutate(ctx context.Context, tenantID, conversationID string, event string, fn func(txCtx context.Context, cid uuid.UUID) error) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return fmt.Errorf("invalid tenant id: %w", err)
	}
	cid, err := uuid.Parse(conversationID)
	if err != nil {
		return fmt.Errorf("invalid conversation id: %w", err)
	}

	if err := s.conversationRepo.RunInTenantTx(repository.WithTenantID(subCtx, tid), func(txCtx context.Context) error {
		return fn(txCtx, cid)
	}); err != nil {
		return err
	}

	if s.centrifugo != nil {
		_ = s.centrifugo.BroadcastConversationUpdate(subCtx, cid.String(), map[string]interface{}{"type": event})
		_ = s.centrifugo.PublishToTenant(subCtx, tid.String(), map[string]interface{}{
			"type":            "conversation_updated",
			"event":           event,
			"conversation_id": cid.String(),
		})
	}
	return nil
}

// Assign assigns the conversation to a specific agent (owner/supervisor action).
func (s *ConversationService) Assign(ctx context.Context, tenantID, conversationID, agentID string) error {
	aid, err := uuid.Parse(agentID)
	if err != nil {
		return fmt.Errorf("invalid agent id: %w", err)
	}
	return s.mutate(ctx, tenantID, conversationID, "assigned", func(txCtx context.Context, cid uuid.UUID) error {
		return s.conversationRepo.Assign(txCtx, cid, &aid)
	})
}

// Unassign releases the conversation back to the unassigned pool.
func (s *ConversationService) Unassign(ctx context.Context, tenantID, conversationID string) error {
	return s.mutate(ctx, tenantID, conversationID, "unassigned", func(txCtx context.Context, cid uuid.UUID) error {
		return s.conversationRepo.Assign(txCtx, cid, nil)
	})
}

// Complete marks the conversation completed; Reopen sets it back to open/unassigned.
func (s *ConversationService) Complete(ctx context.Context, tenantID, conversationID string) error {
	return s.mutate(ctx, tenantID, conversationID, "completed", func(txCtx context.Context, cid uuid.UUID) error {
		return s.conversationRepo.UpdateStatus(txCtx, cid, entity.ConversationStatusCompleted)
	})
}

func (s *ConversationService) Reopen(ctx context.Context, tenantID, conversationID string) error {
	return s.mutate(ctx, tenantID, conversationID, "reopened", func(txCtx context.Context, cid uuid.UUID) error {
		return s.conversationRepo.UpdateStatus(txCtx, cid, entity.ConversationStatusOpen)
	})
}

func (s *ConversationService) SetPinned(ctx context.Context, tenantID, conversationID string, pinned bool) error {
	return s.mutate(ctx, tenantID, conversationID, "pinned", func(txCtx context.Context, cid uuid.UUID) error {
		return s.conversationRepo.SetPinned(txCtx, cid, pinned)
	})
}

func (s *ConversationService) SetArchived(ctx context.Context, tenantID, conversationID string, archived bool) error {
	return s.mutate(ctx, tenantID, conversationID, "archived", func(txCtx context.Context, cid uuid.UUID) error {
		return s.conversationRepo.SetArchived(txCtx, cid, archived)
	})
}

func (s *ConversationService) MarkRead(ctx context.Context, tenantID, conversationID string) error {
	return s.mutate(ctx, tenantID, conversationID, "read", func(txCtx context.Context, cid uuid.UUID) error {
		return s.conversationRepo.MarkRead(txCtx, cid)
	})
}

func (s *ConversationService) List(ctx context.Context, tenantID, status, agentID string) ([]*entity.Convertation, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id: %w", err)
	}
	var agentPtr *uuid.UUID
	if agentID != "" {
		if a, err := uuid.Parse(agentID); err == nil {
			agentPtr = &a
		}
	}

	var convs []*entity.Convertation
	err = s.conversationRepo.RunInTenantTx(repository.WithTenantID(subCtx, tid), func(txCtx context.Context) error {
		convs, err = s.conversationRepo.List(txCtx, tid, status, agentPtr)
		return err
	})
	if err != nil {
		return nil, err
	}
	return convs, nil
}

func (s *ConversationService) ListMessages(ctx context.Context, tenantID, conversationID string, limit int) ([]*entity.Message, error) {
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

	var messages []*entity.Message
	err = s.messageRepo.RunInTenantTx(repository.WithTenantID(subCtx, tid), func(txCtx context.Context) error {
		messages, err = s.messageRepo.ListByConversation(txCtx, cid, limit)
		return err
	})
	if err != nil {
		return nil, err
	}
	return messages, nil
}
