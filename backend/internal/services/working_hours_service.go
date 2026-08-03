package services

import (
	"context"
	"fmt"
	"github/socialforge/internal/dto"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/repository"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WorkingHoursService struct {
	repo   repository.AgentWorkingHoursRepository
	logger *zap.Logger
}

func NewWorkingHoursService(repo repository.AgentWorkingHoursRepository, logger *zap.Logger) *WorkingHoursService {
	return &WorkingHoursService{repo: repo, logger: logger}
}

func (s *WorkingHoursService) List(ctx context.Context, tenantID, userID string) ([]*entity.AgentWorkingHours, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id: %w", err)
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	var out []*entity.AgentWorkingHours
	err = s.repo.RunInTenantTx(repository.WithTenantID(subCtx, tid), func(txCtx context.Context) error {
		out, err = s.repo.ListByUser(txCtx, uid)
		return err
	})
	return out, err
}

func (s *WorkingHoursService) Replace(ctx context.Context, tenantID, userID string, req *dto.SetWorkingHoursRequest) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return fmt.Errorf("invalid tenant id: %w", err)
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}
	hours := make([]*entity.AgentWorkingHours, 0, len(req.Hours))
	for _, h := range req.Hours {
		hours = append(hours, &entity.AgentWorkingHours{
			DayOfWeek: h.DayOfWeek,
			StartTime: h.StartTime,
			EndTime:   h.EndTime,
			IsActive:  h.IsActive,
		})
	}
	return s.repo.RunInTenantTx(repository.WithTenantID(subCtx, tid), func(txCtx context.Context) error {
		return s.repo.ReplaceForUser(txCtx, tid, uid, hours)
	})
}
