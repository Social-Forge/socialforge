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

type ContactService struct {
	repo   repository.ContactRepository
	logger *zap.Logger
}

func NewContactService(repo repository.ContactRepository, logger *zap.Logger) *ContactService {
	return &ContactService{repo: repo, logger: logger}
}

// List returns a paginated, tenant-scoped page of contacts + the total count.
func (s *ContactService) List(ctx context.Context, tenantID, channelID, search string, limit, offset int) ([]*entity.Contact, int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid tenant id: %w", err)
	}
	var chID *uuid.UUID
	if channelID != "" {
		if id, err := uuid.Parse(channelID); err == nil {
			chID = &id
		}
	}

	var contacts []*entity.Contact
	var total int64
	err = s.repo.RunInTenantTx(repository.WithTenantID(subCtx, tid), func(txCtx context.Context) error {
		contacts, total, err = s.repo.List(txCtx, tid, chID, search, limit, offset)
		return err
	})
	if err != nil {
		return nil, 0, err
	}
	return contacts, total, nil
}

func (s *ContactService) SetBlocked(ctx context.Context, tenantID, id string, blocked bool) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return fmt.Errorf("invalid tenant id: %w", err)
	}
	cid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid contact id: %w", err)
	}
	return s.repo.RunInTenantTx(repository.WithTenantID(subCtx, tid), func(txCtx context.Context) error {
		return s.repo.SetBlocked(txCtx, cid, blocked)
	})
}

func (s *ContactService) Delete(ctx context.Context, tenantID, id string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return fmt.Errorf("invalid tenant id: %w", err)
	}
	cid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid contact id: %w", err)
	}
	return s.repo.RunInTenantTx(repository.WithTenantID(subCtx, tid), func(txCtx context.Context) error {
		return s.repo.Delete(txCtx, cid)
	})
}
