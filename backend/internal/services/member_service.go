package services

import (
	"context"
	"fmt"
	"github/socialforge/internal/dto"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/repository"
	"github/socialforge/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// agentRoleSlugs are the roles that count against the tenant's agent quota
// (spec: "10 Agent CS termasuk supervisor").
var agentRoleSlugs = []string{entity.RoleSupervisor, entity.RoleAgent}

type MemberService struct {
	userRepo   repository.UserRepository
	roleRepo   repository.RoleRepository
	tenantRepo repository.TenantRepository
	logger     *zap.Logger
}

func NewMemberService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	tenantRepo repository.TenantRepository,
	logger *zap.Logger,
) *MemberService {
	return &MemberService{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		tenantRepo: tenantRepo,
		logger:     logger,
	}
}

func (s *MemberService) List(ctx context.Context, tenantID, roleSlug string) ([]*entity.UserTenantWithDetails, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id: %w", err)
	}
	return s.userRepo.ListTenantMembers(subCtx, tid, roleSlug)
}

func (s *MemberService) Create(ctx context.Context, tenantID string, req *dto.CreateMemberRequest) (*entity.UserTenantWithDetails, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id: %w", err)
	}

	if !utils.IsStrongPassword(req.Password) {
		return nil, dto.ErrWeakPassword
	}

	// Plan quota: supervisors + agents must stay within tenant.MaxAgents.
	tenant, err := s.tenantRepo.FindByID(subCtx, tid)
	if err != nil {
		return nil, fmt.Errorf("tenant not found: %w", err)
	}
	count, err := s.userRepo.CountTenantMembersInRoles(subCtx, tid, agentRoleSlugs)
	if err != nil {
		return nil, err
	}
	if count >= tenant.MaxAgents {
		return nil, fmt.Errorf("agent quota reached (max %d) for plan %s", tenant.MaxAgents, tenant.SubscriptionPlan)
	}

	if exists, _ := s.userRepo.ExistsByEmail(subCtx, req.Email); exists {
		return nil, dto.ErrEmailAlreadyExists
	}

	role, err := s.roleRepo.GetByName(subCtx, req.Role)
	if err != nil {
		return nil, fmt.Errorf("role %q not found: %w", req.Role, err)
	}

	hash, err := utils.GeneratePasswordHash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now()
	user := &entity.User{
		ID:              uuid.New(),
		Email:           req.Email,
		FullName:        req.FullName,
		Phone:           entity.NewNullString(req.Phone),
		PasswordHash:    hash,
		Status:          entity.UserStatusActive,
		EmailVerifiedAt: entity.NewNullTime(now), // owner-created accounts are pre-verified
		CreatedAt:       now,
	}
	membership := &entity.UserTenant{
		UserID:   user.ID,
		TenantID: tid,
		RoleID:   role.ID,
		IsActive: true,
	}

	if err := s.userRepo.WithTransaction(subCtx, func(tx pgx.Tx) error {
		if err := s.userRepo.CreateTx(subCtx, tx, user); err != nil {
			return err
		}
		return s.userRepo.CreateUserTenantTx(subCtx, tx, membership)
	}); err != nil {
		s.logger.Error("Failed to create member", zap.String("tenant_id", tenantID), zap.Error(err))
		return nil, fmt.Errorf("failed to create member: %w", err)
	}

	return &entity.UserTenantWithDetails{
		User:       *user,
		UserTenant: *membership,
		Tenant:     *tenant,
		Role:       *role,
	}, nil
}

func (s *MemberService) Update(ctx context.Context, tenantID, userTenantID string, req *dto.UpdateMemberRequest) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	utid, err := uuid.Parse(userTenantID)
	if err != nil {
		return fmt.Errorf("invalid membership id: %w", err)
	}
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return fmt.Errorf("invalid tenant id: %w", err)
	}

	membership, err := s.userRepo.FindUserTenantByID(subCtx, utid)
	if err != nil {
		return err
	}
	if membership.TenantID != tid {
		return fmt.Errorf("membership does not belong to this tenant")
	}

	role, err := s.roleRepo.GetByName(subCtx, req.Role)
	if err != nil {
		return fmt.Errorf("role %q not found: %w", req.Role, err)
	}

	isActive := membership.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	return s.userRepo.UpdateUserTenant(subCtx, utid, role.ID, isActive)
}

func (s *MemberService) Delete(ctx context.Context, tenantID, userTenantID string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	utid, err := uuid.Parse(userTenantID)
	if err != nil {
		return fmt.Errorf("invalid membership id: %w", err)
	}
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return fmt.Errorf("invalid tenant id: %w", err)
	}

	membership, err := s.userRepo.FindUserTenantByID(subCtx, utid)
	if err != nil {
		return err
	}
	if membership.TenantID != tid {
		return fmt.Errorf("membership does not belong to this tenant")
	}
	// Never allow removing an owner membership through member management.
	if membership.RoleID != uuid.Nil {
		if role, err := s.roleRepo.FindByID(subCtx, membership.RoleID); err == nil && role.Name == entity.RoleTenantOwner {
			return fmt.Errorf("cannot remove the tenant owner")
		}
	}
	return s.userRepo.SoftDeleteUserTenant(subCtx, utid)
}
