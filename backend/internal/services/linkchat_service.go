package services

import (
	"context"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/repository"
	"time"

	"go.uber.org/zap"
)

type LinkchatService struct {
	divisionRepo repository.DivisionRepository
	logger       *zap.Logger
}

func NewLinkchatService(divisionRepo repository.DivisionRepository, logger *zap.Logger) *LinkchatService {
	return &LinkchatService{divisionRepo: divisionRepo, logger: logger}
}

// Resolve returns the public info for a shared division link (linkchat).
func (s *LinkchatService) Resolve(ctx context.Context, token string) (map[string]interface{}, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	div, err := s.divisionRepo.FindByLinkURL(subCtx, token)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"division_id": div.ID.String(),
		"tenant_id":   div.TenantID.String(),
		"name":        div.Name,
		"description": nullStr(div.Description),
	}, nil
}

func nullStr(n entity.NullString) string {
	if n.Valid {
		return n.String
	}
	return ""
}
