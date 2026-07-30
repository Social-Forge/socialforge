package repository

import (
	"context"
	"errors"
	"fmt"
	"github/socialforge/internal/entity"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type schemaAwareUserRepository struct {
	UserRepository
}

func (r *schemaAwareUserRepository) GetUserTenantWithDetailsByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserTenantWithDetails, error) {
	subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	user, err := r.UserRepository.FindByID(subCtx, userID)
	if err != nil {
		return nil, err
	}

	var tenant entity.Tenant
	if err := pgxscan.Get(subCtx, r.UserRepository.(*userRepository).db, &tenant, `SELECT * FROM tenants WHERE id = $1 AND deleted_at IS NULL`, user.TenantID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("tenant not found")
		}
		return nil, fmt.Errorf("failed to load tenant: %w", err)
	}

	var role entity.Role
	if err := pgxscan.Get(subCtx, r.UserRepository.(*userRepository).db, &role, `SELECT * FROM roles WHERE id = $1 AND deleted_at IS NULL`, user.RoleID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to load role: %w", err)
	}

	var rolePermissions []entity.RolePermissionSummary
	permQuery := `
		SELECT
			r.name AS role_name,
			p.name AS permission_name,
			p.resource AS permission_resource,
			p.action AS permission_action
		FROM role_permissions rp
		JOIN permissions p ON rp.permission_id = p.id AND p.deleted_at IS NULL
		JOIN roles r ON rp.role_id = r.id AND r.deleted_at IS NULL
		WHERE rp.role_id = $1 AND rp.deleted_at IS NULL
		ORDER BY p.resource, p.action
	`
	if err := pgxscan.Select(subCtx, r.UserRepository.(*userRepository).db, &rolePermissions, permQuery, user.RoleID); err != nil {
		return nil, fmt.Errorf("failed to load role permissions: %w", err)
	}

	return &entity.UserTenantWithDetails{
		User: *user,
		UserTenant: entity.UserTenant{
			ID:        user.ID,
			UserID:    user.ID,
			TenantID:  user.TenantID,
			RoleID:    user.RoleID,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Tenant:          tenant,
		Role:            role,
		RolePermissions: rolePermissions,
		Metadata: map[string]any{
			"permission_count": len(rolePermissions),
			"user_status":      user.Status,
		},
	}, nil
}

func (r *schemaAwareUserRepository) GetUserTenantWithDetailsByTenantID(ctx context.Context, tenantID uuid.UUID) (*entity.UserTenantWithDetails, error) {
	subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var userID uuid.UUID
	err := r.UserRepository.(*userRepository).db.QueryRow(subCtx, `
		SELECT id
		FROM users
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`, tenantID).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user tenant not found")
		}
		return nil, fmt.Errorf("failed to find user for tenant: %w", err)
	}

	return r.GetUserTenantWithDetailsByUserID(ctx, userID)
}

func (r *schemaAwareUserRepository) GetUserTenantWithDetailsWithNested(ctx context.Context, userID uuid.UUID) (*entity.UserTenantWithDetailsNested, error) {
	details, err := r.GetUserTenantWithDetailsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	nested := entity.UserTenantWithDetailsNested(*details)
	return &nested, nil
}

func (r *schemaAwareUserRepository) GetUserByRole(ctx context.Context, role string, opts *ListOptions) ([]*entity.User, int64, error) {
	subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if opts == nil {
		opts = NewListOptions()
	}

	totalRows, err := r.CountUsersByRole(ctx, *opts.Filter.TenantID, role, opts.Filter)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			u.id, u.tenant_id, u.role_id, split_part(u.email, '@', 1) AS username,
			u.email, u.password_hash, u.full_name, u.phone, u.avatar_url, u.two_fa_secret,
			u.status, (u.email_verified_at IS NOT NULL) AS is_verified,
			u.email_verified_at, u.last_login_at, u.created_at, u.updated_at, u.deleted_at
		FROM users u
		JOIN roles r ON u.role_id = r.id AND r.deleted_at IS NULL
		WHERE u.tenant_id = $1 AND r.slug = $2 AND u.deleted_at IS NULL
	`
	args := []interface{}{*opts.Filter.TenantID, role}
	var users []*entity.User
	if err := pgxscan.Select(subCtx, r.UserRepository.(*userRepository).db, &users, query, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, nil
		}
		return nil, 0, fmt.Errorf("failed to list users by role: %w", err)
	}

	return users, totalRows, nil
}

func (r *schemaAwareUserRepository) SetEmailVerified(ctx context.Context, id uuid.UUID, isVerified bool) error {
	subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE users
		SET email_verified_at = NOW(),
			status = CASE WHEN $1 THEN 'active' ELSE status END
		WHERE id = $2 AND deleted_at IS NULL
	`
	cmdTag, err := r.UserRepository.(*userRepository).db.Exec(subCtx, query, isVerified, id)
	if err != nil {
		return fmt.Errorf("failed to set email verified: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already updated")
	}
	return nil
}

func (r *schemaAwareUserRepository) CountUsersByRole(ctx context.Context, tenantID uuid.UUID, roleSlug string, filter *Filter) (int64, error) {
	subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT COUNT(*)
		FROM users u
		JOIN roles r ON u.role_id = r.id AND r.deleted_at IS NULL
		WHERE u.tenant_id = $1
		  AND r.slug = $2
		  AND u.deleted_at IS NULL
	`
	var count int64
	if err := r.UserRepository.(*userRepository).db.QueryRow(subCtx, query, tenantID, roleSlug).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count users by role: %w", err)
	}
	return count, nil
}

func init() {
	_ = zap.NewNop()
}
