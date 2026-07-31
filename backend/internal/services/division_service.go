package services

import (
	"context"
	"fmt"
	"github/socialforge/internal/dto"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/repository"
	"github/socialforge/internal/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DivisionService struct {
	divisionRepo repository.DivisionRepository
	tenantRepo   repository.TenantRepository
	logger       *zap.Logger
}

func NewDivisionService(divisionRepo repository.DivisionRepository, tenantRepo repository.TenantRepository, logger *zap.Logger) *DivisionService {
	return &DivisionService{
		divisionRepo: divisionRepo,
		tenantRepo:   tenantRepo,
		logger:       logger,
	}
}

// tenantCtx parses the tenant id and stores it in the context so the
// RLS-bound transaction (RunInTenantTx) can apply app.current_tenant.
func (s *DivisionService) tenantCtx(ctx context.Context, tenantID string) (context.Context, uuid.UUID, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return ctx, uuid.Nil, fmt.Errorf("invalid tenant id: %w", err)
	}
	return repository.WithTenantID(ctx, tid), tid, nil
}

func (s *DivisionService) Create(ctx context.Context, tenantID string, req *dto.CreateDivisionRequest) (*entity.Division, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tctx, tid, err := s.tenantCtx(subCtx, tenantID)
	if err != nil {
		return nil, err
	}

	tenant, err := s.tenantRepo.FindByID(subCtx, tid)
	if err != nil {
		return nil, fmt.Errorf("tenant not found: %w", err)
	}

	routing := req.RoutingType
	if routing == "" {
		routing = entity.RoutingEqual
	}

	division := &entity.Division{
		ID:          uuid.New(),
		TenantID:    tid,
		Name:        req.Name,
		Description: entity.NewNullString(req.Description),
		RoutingType: routing,
		IsActive:    true,
		LinkURL:     entity.NewNullString(utils.GenerateSecureToken(10)),
	}

	err = s.divisionRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		// Plan quota: enforce max_divisions per tenant.
		current, err := s.divisionRepo.Count(txCtx, &repository.Filter{TenantID: &tid})
		if err != nil {
			return err
		}
		if current >= int64(tenant.MaxDivisions) {
			return fmt.Errorf("division quota reached (max %d) for plan %s", tenant.MaxDivisions, tenant.SubscriptionPlan)
		}

		slug := utils.GenerateSlugUnicodeV2(req.Name)
		if exists, _ := s.divisionRepo.ExistsBySlug(txCtx, tid, slug); exists {
			slug = fmt.Sprintf("%s-%s", slug, strings.ToLower(utils.GenerateSecureToken(4)))
		}
		division.Slug = slug
		return s.divisionRepo.Create(txCtx, division)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create division: %w", err)
	}
	return division, nil
}

func (s *DivisionService) List(ctx context.Context, tenantID string, opts *repository.ListOptions) ([]*entity.Division, int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tctx, tid, err := s.tenantCtx(subCtx, tenantID)
	if err != nil {
		return nil, 0, err
	}
	if opts == nil {
		opts = repository.NewListOptions()
	}
	if opts.Filter == nil {
		opts.Filter = &repository.Filter{}
	}
	opts.Filter.TenantID = &tid

	var divisions []*entity.Division
	var total int64
	err = s.divisionRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		divisions, total, err = s.divisionRepo.Search(txCtx, opts)
		return err
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list divisions: %w", err)
	}
	return divisions, total, nil
}

func (s *DivisionService) Get(ctx context.Context, tenantID, id string) (*entity.Division, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	divisionID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid division id: %w", err)
	}
	tctx, _, err := s.tenantCtx(subCtx, tenantID)
	if err != nil {
		return nil, err
	}

	var division *entity.Division
	err = s.divisionRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		division, err = s.divisionRepo.FindByID(txCtx, divisionID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return division, nil
}

func (s *DivisionService) Update(ctx context.Context, tenantID, id string, req *dto.UpdateDivisionRequest) (*entity.Division, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	divisionID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid division id: %w", err)
	}
	tctx, tid, err := s.tenantCtx(subCtx, tenantID)
	if err != nil {
		return nil, err
	}

	var updated *entity.Division
	err = s.divisionRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		existing, err := s.divisionRepo.FindByID(txCtx, divisionID)
		if err != nil {
			return err
		}
		existing.Name = req.Name
		existing.Description = entity.NewNullString(req.Description)
		if req.RoutingType != "" {
			existing.RoutingType = req.RoutingType
		}
		if req.IsActive != nil {
			existing.IsActive = *req.IsActive
		}
		existing.TenantID = tid
		updated, err = s.divisionRepo.Update(txCtx, existing)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update division: %w", err)
	}
	return updated, nil
}

func (s *DivisionService) Delete(ctx context.Context, tenantID, id string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	divisionID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid division id: %w", err)
	}
	tctx, _, err := s.tenantCtx(subCtx, tenantID)
	if err != nil {
		return err
	}
	return s.divisionRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		return s.divisionRepo.Delete(txCtx, divisionID)
	})
}

func (s *DivisionService) ListMembers(ctx context.Context, tenantID, divisionID string) ([]*entity.DivisionMember, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	did, err := uuid.Parse(divisionID)
	if err != nil {
		return nil, fmt.Errorf("invalid division id: %w", err)
	}
	tctx, _, err := s.tenantCtx(subCtx, tenantID)
	if err != nil {
		return nil, err
	}

	var members []*entity.DivisionMember
	err = s.divisionRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		members, err = s.divisionRepo.GetDivisionMembers(txCtx, did)
		return err
	})
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (s *DivisionService) AddMember(ctx context.Context, tenantID, divisionID string, req *dto.AddDivisionMemberRequest) (*entity.DivisionMember, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	did, err := uuid.Parse(divisionID)
	if err != nil {
		return nil, fmt.Errorf("invalid division id: %w", err)
	}
	tctx, _, err := s.tenantCtx(subCtx, tenantID)
	if err != nil {
		return nil, err
	}

	member := &entity.DivisionMember{
		ID:           uuid.New(),
		UserTenantID: req.UserTenantID,
		DivisionID:   did,
		IsActive:     true,
		JoinedAt:     time.Now(),
	}
	err = s.divisionRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		// Ensure the division exists within this tenant (RLS-scoped) first.
		if _, err := s.divisionRepo.FindByID(txCtx, did); err != nil {
			return err
		}
		return s.divisionRepo.AddMember(txCtx, member)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add member: %w", err)
	}
	return member, nil
}

func (s *DivisionService) RemoveMember(ctx context.Context, tenantID, divisionID, userTenantID string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	did, err := uuid.Parse(divisionID)
	if err != nil {
		return fmt.Errorf("invalid division id: %w", err)
	}
	utid, err := uuid.Parse(userTenantID)
	if err != nil {
		return fmt.Errorf("invalid user_tenant id: %w", err)
	}
	tctx, _, err := s.tenantCtx(subCtx, tenantID)
	if err != nil {
		return err
	}
	return s.divisionRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		return s.divisionRepo.RemoveMember(txCtx, utid, did)
	})
}
