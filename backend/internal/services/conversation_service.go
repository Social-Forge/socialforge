package services

import (
	"context"
	"fmt"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/repository"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ConversationService struct {
	conversationRepo repository.ConversationRepository
	messageRepo      repository.MessageRepository
	logger           *zap.Logger
}

func NewConversationService(
	conversationRepo repository.ConversationRepository,
	messageRepo repository.MessageRepository,
	logger *zap.Logger,
) *ConversationService {
	return &ConversationService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		logger:           logger,
	}
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
