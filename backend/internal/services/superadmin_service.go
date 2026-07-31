package services

import (
	"context"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/repository"
	"time"

	"go.uber.org/zap"
)

// SuperadminService provides cross-tenant platform administration. It operates
// outside tenant RLS scope (the tenants table is not RLS-guarded).
type SuperadminService struct {
	tenantRepo repository.TenantRepository
	logger     *zap.Logger
}

func NewSuperadminService(tenantRepo repository.TenantRepository, logger *zap.Logger) *SuperadminService {
	return &SuperadminService{tenantRepo: tenantRepo, logger: logger}
}

func (s *SuperadminService) ListTenants(ctx context.Context, opts *repository.ListOptions) ([]*entity.Tenant, int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if opts == nil {
		opts = repository.NewListOptions()
	}
	return s.tenantRepo.Search(subCtx, opts)
}
