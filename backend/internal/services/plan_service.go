package services

import (
	"context"
	"fmt"
	"github/socialforge/internal/dto"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// PlanService manages the global plan catalog. Reads are public; mutations are
// superadmin-only (enforced at the route layer).
type PlanService struct {
	repo   repository.PlanRepository
	logger *zap.Logger
}

func NewPlanService(repo repository.PlanRepository, logger *zap.Logger) *PlanService {
	return &PlanService{repo: repo, logger: logger}
}

func (s *PlanService) List(ctx context.Context, activeOnly bool) ([]*entity.Plan, error) {
	return s.repo.List(ctx, activeOnly)
}

func (s *PlanService) GetByCode(ctx context.Context, code string) (*entity.Plan, error) {
	return s.repo.FindByCode(ctx, code)
}

func (s *PlanService) Create(ctx context.Context, req *dto.CreatePlanRequest) (*entity.Plan, error) {
	p := &entity.Plan{
		ID:       uuid.New(),
		Code:     req.Code,
		Name:     req.Name,
		Price:    req.Price,
		Currency: orDefault(req.Currency, "IDR"),
		Interval: orDefault(req.Interval, entity.PlanIntervalMonthly),
		Features: entity.JSONMap(req.Features),
		IsActive: req.IsActive == nil || *req.IsActive,
		Sort:     req.Sort,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PlanService) Update(ctx context.Context, id string, req *dto.UpdatePlanRequest) (*entity.Plan, error) {
	pid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid plan id: %w", err)
	}
	existing, err := s.repo.FindByID(ctx, pid)
	if err != nil {
		return nil, err
	}
	existing.Code = req.Code
	existing.Name = req.Name
	existing.Price = req.Price
	existing.Currency = orDefault(req.Currency, existing.Currency)
	existing.Interval = orDefault(req.Interval, existing.Interval)
	if req.Features != nil {
		existing.Features = entity.JSONMap(req.Features)
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	existing.Sort = req.Sort
	return s.repo.Update(ctx, existing)
}

func (s *PlanService) Delete(ctx context.Context, id string) error {
	pid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid plan id: %w", err)
	}
	return s.repo.Delete(ctx, pid)
}

func orDefault(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
