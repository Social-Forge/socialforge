package repository

import (
	"context"
	"errors"
	"fmt"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	BaseRepository
	// Create operations
	Create(ctx context.Context, user *entity.User) error
	CreateTx(ctx context.Context, tx pgx.Tx, user *entity.User) error
	CreateWithRecovery(ctx context.Context, user *entity.User) error
	// Read operations
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByEmailOrUsername(ctx context.Context, identifier string) (*entity.User, error)
	GetUserTenantWithDetailsByUserID(ctx context.Context, id uuid.UUID) (*entity.UserTenantWithDetails, error)
	GetUserTenantWithDetailsByTenantID(ctx context.Context, tenantID uuid.UUID) (*entity.UserTenantWithDetails, error)
	GetUserTenantWithDetailsWithNested(ctx context.Context, userID uuid.UUID) (*entity.UserTenantWithDetailsNested, error)
	Search(ctx context.Context, opts *ListOptions) ([]*entity.User, int64, error)
	Count(ctx context.Context, filter *Filter) (int64, error)
	// Update operations
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateTx(ctx context.Context, tx pgx.Tx, user *entity.User) (*entity.User, error)
	UpdateWithRecovery(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	UpdateLastLoginTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) error
	UpdateTwoFaSecret(ctx context.Context, id uuid.UUID, twoFaSecret *string) error
	RemoveTwoFaSecret(ctx context.Context, id uuid.UUID) error
	SetEmailVerified(ctx context.Context, id uuid.UUID, isVerified bool) error
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	UpdateAvatar(ctx context.Context, id uuid.UUID, avatarURL string) (string, error)
	// Delete operations
	Delete(ctx context.Context, id uuid.UUID) error // Soft delete
	HardDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	// Check operations
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	// Check two factor authentication
	IsTwoFaEnabled(ctx context.Context, id uuid.UUID) (bool, error)
}
type userRepository struct {
	*baseRepository
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		baseRepository: NewBaseRepository(db).(*baseRepository),
	}
}
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (
			id, email, username, password_hash, full_name, phone, avatar_url,
			status, email_verified_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(subCtx,
		query,
		user.ID,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.FullName,
		user.Phone,
		user.AvatarURL,
		user.Status,
		user.EmailVerifiedAt,
		user.CreatedAt,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			switch pgxErr.ConstraintName {
			case "users_username_key":
				return fmt.Errorf("username %s is already taken", user.Username)
			case "users_email_key":
				return fmt.Errorf("email %s is already registered", user.Email)
			case "users_name_length_check":
				return fmt.Errorf("full name %s is invalid, must be between 2 and 50 characters", user.FullName)
			case "users_username_length_check":
				return fmt.Errorf("username %s is invalid, must be between 3 and 20 characters", user.Username)
			default:
				return fmt.Errorf("unique constraint violation (%s): %w", pgxErr.ConstraintName, err)
			}
		}
		return fmt.Errorf("failed to create new user: %w", err)
	}
	return nil
}
func (r *userRepository) CreateTx(ctx context.Context, tx pgx.Tx, user *entity.User) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (
			id, email, username, password_hash, full_name, phone, avatar_url,
			status, email_verified_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	err := tx.QueryRow(subCtx,
		query,
		user.ID,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.FullName,
		user.Phone,
		user.AvatarURL,
		user.Status,
		user.EmailVerifiedAt,
		user.CreatedAt,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			switch pgxErr.ConstraintName {
			case "users_username_key":
				return fmt.Errorf("username %s is already taken", user.Username)
			case "users_email_key":
				return fmt.Errorf("email %s is already registered", user.Email)
			case "users_name_length_check":
				return fmt.Errorf("full name %s is invalid, must be between 2 and 50 characters", user.FullName)
			case "users_username_length_check":
				return fmt.Errorf("username %s is invalid, must be between 3 and 20 characters", user.Username)
			default:
				return fmt.Errorf("unique constraint violation (%s): %w", pgxErr.ConstraintName, err)
			}
		}
		return fmt.Errorf("failed to create new user: %w", err)
	}
	return nil
}
func (r *userRepository) CreateWithRecovery(ctx context.Context, user *entity.User) error {
	return r.WithTransaction(ctx, func(tx pgx.Tx) error {
		return r.CreateTx(ctx, tx, user)
	})
}
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL
	`
	var user entity.User
	err := pgxscan.Get(subCtx, r.db, &user, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	return &user, nil
}
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL
	`
	var user entity.User
	err := pgxscan.Get(subCtx, r.db, &user, query, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}
func (r *userRepository) FindByTenantID(ctx context.Context, tenantID uuid.UUID) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM users WHERE tenant_id = $1 AND deleted_at IS NULL
	`
	var user entity.User
	err := pgxscan.Get(subCtx, r.db, &user, query, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by tenant_id: %w", err)
	}
	return &user, nil
}

func (r *userRepository) FindByRoleID(ctx context.Context, roleID uuid.UUID) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM users WHERE role_id = $1 AND deleted_at IS NULL
	`
	var user entity.User
	err := pgxscan.Get(subCtx, r.db, &user, query, roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by tenant_id: %w", err)
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM users WHERE username = $1 AND deleted_at IS NULL
	`
	var user entity.User
	err := pgxscan.Get(subCtx, r.db, &user, query, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}
	return &user, nil
}
func (r *userRepository) FindByEmailOrUsername(ctx context.Context, identifier string) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM users WHERE (email = $1 OR username = $1) AND deleted_at IS NULL
	`
	var user entity.User
	err := pgxscan.Get(subCtx, r.db, &user, query, identifier)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by email or username: %w", err)
	}
	return &user, nil
}
func (r *userRepository) GetUserTenantWithDetailsByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserTenantWithDetails, error) {
	subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT 
			json_build_object(
				'user_tenant', json_build_object(
					'id', ut.id,
					'user_id', ut.user_id,
					'tenant_id', ut.tenant_id,
					'role_id', ut.role_id,
					'is_active', ut.is_active,
					'created_at', ut.created_at,
					'updated_at', ut.updated_at
				),
				'user', json_build_object(
					'id', u.id,
					'email', u.email,
					'username', u.username,
					'full_name', u.full_name,
					'phone', u.phone,
					'avatar_url', u.avatar_url,
					'two_fa_secret', u.two_fa_secret,
					'is_active', u.is_active,
					'is_verified', u.is_verified,
					'email_verified_at', u.email_verified_at,
					'last_login_at', u.last_login_at,
					'created_at', u.created_at,
					'updated_at', u.updated_at
				),
				'tenant', json_build_object(
					'id', t.id,
					'name', t.name,
					'slug', t.slug,
					'owner_id', t.owner_id,
					'logo_url', t.logo_url,
					'description', t.description,
					'max_divisions', t.max_divisions,
					'max_agents', t.max_agents,
					'max_quick_replies', t.max_quick_replies,
					'max_pages', t.max_pages,
					'max_whatsapp', t.max_whatsapp,
					'max_meta_whatsapp', t.max_meta_whatsapp,
					'max_meta_messenger', t.max_meta_messenger,
					'max_instagram', t.max_instagram,
					'max_telegram', t.max_telegram,
					'max_webchat', t.max_webchat,
					'max_linkchat', t.max_linkchat,
					'subscription_plan', t.subscription_plan,
					'subscription_status', t.subscription_status,
					'trial_ends_at', t.trial_ends_at,
					'is_active', t.is_active,
					'created_at', t.created_at,
					'updated_at', t.updated_at
				),
				'role', json_build_object(
					'id', r.id,
					'name', r.name,
					'slug', r.slug,
					'description', r.description,
					'level', r.level,
					'created_at', r.created_at,
					'updated_at', r.updated_at
				),
				'role_permissions', COALESCE(
					(
						SELECT json_agg(
							json_build_object(
								'id', rp.id,
								'role_id', rp.role_id,
								'permission_id', rp.permission_id,
								'created_at', rp.created_at,
								'updated_at', rp.updated_at,
								'role_name', r2.name,
								'role_slug', r2.slug,
								'role_level', r2.level,
								'permission_name', p.name,
								'permission_slug', p.slug,
								'permission_resource', p.resource,
								'permission_action', p.action
							)
							ORDER BY p.resource, p.action
						)
						FROM role_permissions rp
						JOIN roles r2 ON rp.role_id = r2.id AND r2.deleted_at IS NULL
						JOIN permissions p ON rp.permission_id = p.id AND p.deleted_at IS NULL
						WHERE rp.role_id = ut.role_id AND rp.deleted_at IS NULL
					),
					'[]'
				),
				'metadata', json_build_object(
					'permission_count', (
						SELECT COUNT(*) 
						FROM role_permissions rp 
						WHERE rp.role_id = ut.role_id AND rp.deleted_at IS NULL
					),
					'user_status', CASE 
						WHEN u.is_active AND ut.is_active THEN 'active'
						WHEN NOT u.is_active THEN 'user_inactive'
						WHEN NOT ut.is_active THEN 'tenant_access_inactive'
						ELSE 'unknown'
					END,
					'last_updated', GREATEST(
						ut.updated_at, 
						u.updated_at, 
						t.updated_at,
						COALESCE((SELECT MAX(updated_at) FROM role_permissions WHERE role_id = ut.role_id), ut.updated_at)
					)
				)
			) as user_tenant_data
		FROM user_tenants ut
		JOIN users u ON ut.user_id = u.id AND u.deleted_at IS NULL
		JOIN tenants t ON ut.tenant_id = t.id AND t.deleted_at IS NULL
		JOIN roles r ON ut.role_id = r.id AND r.deleted_at IS NULL
		WHERE ut.user_id = $1 AND ut.deleted_at IS NULL
		ORDER BY ut.created_at DESC
		LIMIT 1
	`

	var result struct {
		UserTenantData entity.UserTenantWithDetails `db:"user_tenant_data"`
	}

	err := pgxscan.Get(subCtx, r.db, &result, query, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user tenant not found")
		}
		return nil, fmt.Errorf("failed to get user tenant with details: %w", err)
	}

	return &result.UserTenantData, nil
}
func (r *userRepository) GetUserTenantWithDetailsByTenantID(ctx context.Context, tenantID uuid.UUID) (*entity.UserTenantWithDetails, error) {
	subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT 
			json_build_object(
				'user_tenant', json_build_object(
					'id', ut.id,
					'user_id', ut.user_id,
					'tenant_id', ut.tenant_id,
					'role_id', ut.role_id,
					'is_active', ut.is_active,
					'created_at', ut.created_at,
					'updated_at', ut.updated_at
				),
				'user', json_build_object(
					'id', u.id,
					'email', u.email,
					'username', u.username,
					'full_name', u.full_name,
					'phone', u.phone,
					'avatar_url', u.avatar_url,
					'two_fa_secret', u.two_fa_secret,
					'is_active', u.is_active,
					'is_verified', u.is_verified,
					'email_verified_at', u.email_verified_at,
					'last_login_at', u.last_login_at,
					'created_at', u.created_at,
					'updated_at', u.updated_at
				),
				'tenant', json_build_object(
					'id', t.id,
					'name', t.name,
					'slug', t.slug,
					'owner_id', t.owner_id,
					'logo_url', t.logo_url,
					'description', t.description,
					'max_divisions', t.max_divisions,
					'max_agents', t.max_agents,
					'max_quick_replies', t.max_quick_replies,
					'max_pages', t.max_pages,
					'max_whatsapp', t.max_whatsapp,
					'max_meta_whatsapp', t.max_meta_whatsapp,
					'max_meta_messenger', t.max_meta_messenger,
					'max_instagram', t.max_instagram,
					'max_telegram', t.max_telegram,
					'max_webchat', t.max_webchat,
					'max_linkchat', t.max_linkchat,
					'subscription_plan', t.subscription_plan,
					'subscription_status', t.subscription_status,
					'trial_ends_at', t.trial_ends_at,
					'is_active', t.is_active,
					'created_at', t.created_at,
					'updated_at', t.updated_at
				),
				'role', json_build_object(
					'id', r.id,
					'name', r.name,
					'slug', r.slug,
					'description', r.description,
					'level', r.level,
					'created_at', r.created_at,
					'updated_at', r.updated_at
				),
				'role_permissions', COALESCE(
					(
						SELECT json_agg(
							json_build_object(
								'id', rp.id,
								'role_id', rp.role_id,
								'permission_id', rp.permission_id,
								'created_at', rp.created_at,
								'updated_at', rp.updated_at,
								'role_name', r2.name,
								'role_slug', r2.slug,
								'role_level', r2.level,
								'permission_name', p.name,
								'permission_slug', p.slug,
								'permission_resource', p.resource,
								'permission_action', p.action
							)
							ORDER BY p.resource, p.action
						)
						FROM role_permissions rp
						JOIN roles r2 ON rp.role_id = r2.id AND r2.deleted_at IS NULL
						JOIN permissions p ON rp.permission_id = p.id AND p.deleted_at IS NULL
						WHERE rp.role_id = ut.role_id AND rp.deleted_at IS NULL
					),
					'[]'
				),
				'metadata', json_build_object(
					'permission_count', (
						SELECT COUNT(*) 
						FROM role_permissions rp 
						WHERE rp.role_id = ut.role_id AND rp.deleted_at IS NULL
					),
					'user_status', CASE 
						WHEN u.is_active AND ut.is_active THEN 'active'
						WHEN NOT u.is_active THEN 'user_inactive'
						WHEN NOT ut.is_active THEN 'tenant_access_inactive'
						ELSE 'unknown'
					END,
					'last_updated', GREATEST(
						ut.updated_at, 
						u.updated_at, 
						t.updated_at,
						COALESCE((SELECT MAX(updated_at) FROM role_permissions WHERE role_id = ut.role_id), ut.updated_at)
					)
				)
			) as user_tenant_data
		FROM user_tenants ut
		JOIN users u ON ut.user_id = u.id AND u.deleted_at IS NULL
		JOIN tenants t ON ut.tenant_id = t.id AND t.deleted_at IS NULL
		JOIN roles r ON ut.role_id = r.id AND r.deleted_at IS NULL
		WHERE ut.tenant_id = $1 AND ut.deleted_at IS NULL
		ORDER BY ut.created_at DESC
		LIMIT 1
	`

	var result struct {
		UserTenantData entity.UserTenantWithDetails `db:"user_tenant_data"`
	}

	err := pgxscan.Get(subCtx, r.db, &result, query, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user tenant not found")
		}
		return nil, fmt.Errorf("failed to get user tenant with details: %w", err)
	}

	return &result.UserTenantData, nil
}
func (r *userRepository) GetUserTenantWithDetailsWithNested(ctx context.Context, userID uuid.UUID) (*entity.UserTenantWithDetailsNested, error) {
	subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT 
			json_build_object(
				'user_tenant', json_build_object(
					'id', ut.id,
					'user_id', ut.user_id,
					'tenant_id', ut.tenant_id,
					'role_id', ut.role_id,
					'is_active', ut.is_active,
					'created_at', ut.created_at,
					'updated_at', ut.updated_at
				),
				'user', json_build_object(
					'id', u.id,
					'email', u.email,
					'username', u.username,
					'full_name', u.full_name,
					'phone', u.phone,
					'avatar_url', u.avatar_url,
					'is_active', u.is_active,
					'is_verified', u.is_verified,
					'email_verified_at', u.email_verified_at,
					'last_login_at', u.last_login_at,
					'created_at', u.created_at,
					'updated_at', u.updated_at
				),
				'tenant', json_build_object(
					'id', t.id,
					'name', t.name,
					'slug', t.slug,
					'owner_id', t.owner_id,
					'logo_url', t.logo_url,
					'description', t.description,
					'max_divisions', t.max_divisions,
					'max_agents', t.max_agents,
					'max_quick_replies', t.max_quick_replies,
					'max_pages', t.max_pages,
					'max_whatsapp', t.max_whatsapp,
					'max_meta_whatsapp', t.max_meta_whatsapp,
					'max_meta_messenger', t.max_meta_messenger,
					'max_instagram', t.max_instagram,
					'max_telegram', t.max_telegram,
					'max_webchat', t.max_webchat,
					'max_linkchat', t.max_linkchat,
					'subscription_plan', t.subscription_plan,
					'subscription_status', t.subscription_status,
					'trial_ends_at', t.trial_ends_at,
					'is_active', t.is_active,
					'created_at', t.created_at,
					'updated_at', t.updated_at
				),
				'role', json_build_object(
					'id', r.id,
					'name', r.name,
					'slug', r.slug,
					'description', r.description,
					'level', r.level,
					'created_at', r.created_at,
					'updated_at', r.updated_at
				),
				'role_permissions', COALESCE(
					(
						SELECT json_agg(
							json_build_object(
								'role_permission', json_build_object(
									'id', rp.id,
									'role_id', rp.role_id,
									'permission_id', rp.permission_id,
									'created_at', rp.created_at,
									'updated_at', rp.updated_at
								),
								'role', json_build_object(
									'id', r2.id,
									'name', r2.name,
									'slug', r2.slug,
									'description', r2.description,
									'level', r2.level,
									'created_at', r2.created_at,
									'updated_at', r2.updated_at
								),
								'permission', json_build_object(
									'id', p.id,
									'name', p.name,
									'slug', p.slug,
									'resource', p.resource,
									'action', p.action,
									'description', p.description,
									'created_at', p.created_at,
									'updated_at', p.updated_at
								)
							)
							ORDER BY p.resource, p.action
						)
						FROM role_permissions rp
						JOIN roles r2 ON rp.role_id = r2.id AND r2.deleted_at IS NULL
						JOIN permissions p ON rp.permission_id = p.id AND p.deleted_at IS NULL
						WHERE rp.role_id = ut.role_id AND rp.deleted_at IS NULL
					),
					'[]'
				),
				'metadata', json_build_object(
					'permission_count', (
						SELECT COUNT(*) 
						FROM role_permissions rp 
						WHERE rp.role_id = ut.role_id AND rp.deleted_at IS NULL
					),
					'user_status', CASE 
						WHEN u.is_active AND ut.is_active THEN 'active'
						WHEN NOT u.is_active THEN 'user_inactive'
						WHEN NOT ut.is_active THEN 'tenant_access_inactive'
						ELSE 'unknown'
					END,
					'last_updated', GREATEST(
						ut.updated_at, 
						u.updated_at, 
						t.updated_at,
						COALESCE((SELECT MAX(updated_at) FROM role_permissions WHERE role_id = ut.role_id), ut.updated_at)
					)
				)
			) as user_tenant_data
		FROM user_tenants ut
		JOIN users u ON ut.user_id = u.id AND u.deleted_at IS NULL
		JOIN tenants t ON ut.tenant_id = t.id AND t.deleted_at IS NULL
		JOIN roles r ON ut.role_id = r.id AND r.deleted_at IS NULL  -- ✅ Join roles
		WHERE ut.user_id = $1 AND ut.deleted_at IS NULL
		ORDER BY ut.created_at DESC
		LIMIT 1
	`

	var result struct {
		UserTenantData entity.UserTenantWithDetailsNested `db:"user_tenant_data"`
	}

	err := pgxscan.Get(subCtx, r.db, &result, query, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user tenant not found")
		}
		return nil, fmt.Errorf("failed to get user tenant with details: %w", err)
	}

	return &result.UserTenantData, nil
}
func (r *userRepository) Search(ctx context.Context, opts *ListOptions) ([]*entity.User, int64, error) {
	subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if opts == nil {
		opts = NewListOptions()
	}

	totalRows, err := r.Count(ctx, opts.Filter)
	if err != nil {
		return nil, 0, err
	}

	qb := r.buildBaseQuery("SELECT * FROM users", opts.Filter)

	if opts.OrderBy != "" {
		qb.OrderByField(opts.OrderBy, opts.OrderDir)
	} else {
		qb.OrderByField("created_at", "DESC")
	}
	if opts.Pagination != nil && opts.Pagination.Limit > 0 {
		qb.WithLimit(opts.Pagination.Limit)
		if opts.Pagination.Page > 1 {
			qb.WithOffset(opts.Pagination.GetOffset())
		}
	}

	query, args := qb.Build()
	var users []*entity.User
	err = pgxscan.Select(subCtx, r.db, &users, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, fmt.Errorf("no users found")
		}
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, totalRows, nil
}
func (r *userRepository) Count(ctx context.Context, filter *Filter) (int64, error) {
	subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	qb := r.buildBaseQuery("SELECT COUNT(*) FROM users", filter)
	query, args := qb.Build()

	var count int64
	err := r.db.QueryRow(subCtx, query, args...).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("no users found")
		}
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}
func (r *userRepository) GetUserByRole(ctx context.Context, role string, opts *ListOptions) ([]*entity.User, int64, error) {
	subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if opts == nil {
		opts = NewListOptions()
	}

	baseQuery := `
		SELECT 
			u.id, u.email, u.username, u.password_hash, u.full_name, 
			u.phone, u.avatar_url, u.two_fa_secret, u.is_active, u.is_verified,
			u.email_verified_at, u.last_login_at, u.created_at, u.updated_at, u.deleted_at
		FROM users u
		INNER JOIN user_tenants ut ON u.id = ut.user_id
		INNER JOIN roles r ON ut.role_id = r.id
		WHERE ut.tenant_id = $? AND r.slug = $? AND u.deleted_at IS NULL
	`

	qb := NewQueryBuilder(baseQuery)
	qb.Where("ut.tenant_id = $?", opts.Filter.TenantID)
	qb.Where("r.slug = $?", role)
	qb.Where("u.deleted_at IS NULL")
	qb.Where("ut.deleted_at IS NULL")

	if opts.Filter != nil {
		if opts.Filter.Search != "" {
			searchPattern := "%" + opts.Filter.Search + "%"
			qb.Where("(u.full_name ILIKE $? OR u.email ILIKE $? OR u.username ILIKE $?)",
				searchPattern, searchPattern, searchPattern)
		}

		if opts.Filter.IsActive != nil {
			qb.Where("u.is_active = $?", *opts.Filter.IsActive)
		}

		if opts.Filter.IsVerified != nil {
			qb.Where("u.is_verified = $?", *opts.Filter.IsVerified)
		}

		if opts.Filter.RangeDate != nil {
			var startDate time.Time
			var endDate time.Time

			if !opts.Filter.RangeDate.StartDate.IsZero() {
				startDate = opts.Filter.RangeDate.StartDate
			} else {
				startDate = time.Now().AddDate(0, 0, -7)
			}
			if !opts.Filter.RangeDate.EndDate.IsZero() {
				endDate = opts.Filter.RangeDate.EndDate
			} else {
				endDate = time.Now()
			}
			if !startDate.IsZero() || !endDate.IsZero() {
				qb.Where("created_at BETWEEN $? AND $?", startDate, endDate)
			}
		}
	}

	if opts.OrderBy != "" {
		safeOrderBy := opts.OrderBy
		if safeOrderBy == "created_at" || safeOrderBy == "updated_at" {
			safeOrderBy = "u." + safeOrderBy
		}
		qb.OrderByField(safeOrderBy, opts.OrderDir)
	} else {
		qb.OrderByField("u.created_at", "DESC")
	}

	if opts.Pagination != nil && opts.Pagination.Limit > 0 {
		qb.WithLimit(opts.Pagination.Limit)
		if opts.Pagination.Page > 1 {
			qb.WithOffset(opts.Pagination.GetOffset())
		}
	}

	query, args := qb.Build()
	var users []*entity.User
	err := pgxscan.Select(subCtx, r.db, &users, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, fmt.Errorf("no users found")
		}
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	totalRows, err := r.CountUsersByRole(ctx, *opts.Filter.TenantID, role, opts.Filter)
	if err != nil {
		return nil, 0, err
	}

	return users, totalRows, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE users SET
			email = $1,
			username = $2,
			full_name = $3,
			phone = $4,
			avatar_url = $5,
			is_active = $6,
			is_verified = $7,
			email_verified_at = $8
		WHERE id = $9 AND deleted_at IS NULL
		RETURNING id, username, email, full_name, phone, avatar_url, is_active, is_verified, email_verified_at, last_login_at, created_at, updated_at
	`
	args := []interface{}{
		user.Email,
		user.Username,
		user.FullName,
		user.Phone,
		user.AvatarURL,
		user.IsActive,
		user.IsVerified,
		user.EmailVerifiedAt,
		user.ID,
	}

	updatedUser := &entity.User{}
	err := r.db.QueryRow(
		subCtx,
		query,
		args...).Scan(
		&updatedUser.ID,
		&updatedUser.Username,
		&updatedUser.Email,
		&updatedUser.FullName,
		&updatedUser.Phone,
		&updatedUser.AvatarURL,
		&updatedUser.IsActive,
		&updatedUser.IsVerified,
		&updatedUser.EmailVerifiedAt,
		&updatedUser.LastLoginAt,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			switch pgxErr.ConstraintName {
			case "users_username_key":
				return nil, fmt.Errorf("username %s is already taken", user.Username)
			case "users_email_key":
				return nil, fmt.Errorf("email %s is already registered", user.Email)
			case "users_name_length_check":
				return nil, fmt.Errorf("full name %s is invalid, must be between 2 and 50 characters", user.FullName)
			case "users_username_length_check":
				return nil, fmt.Errorf("username %s is invalid, must be between 3 and 20 characters", user.Username)
			default:
				return nil, fmt.Errorf("unique constraint violation (%s): %w", pgxErr.ConstraintName, err)
			}
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return updatedUser, nil
}
func (r *userRepository) UpdateTx(ctx context.Context, tx pgx.Tx, user *entity.User) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE users SET
			email = $1,
			username = $2,
			full_name = $3,
			phone = $4,
			avatar_url = $5,
			is_active = $6,
			is_verified = $7,
			email_verified_at = $8
		WHERE id = $9 AND deleted_at IS NULL
		RETURNING id, username, email, full_name, phone, avatar_url, is_active, is_verified, email_verified_at, last_login_at, created_at, updated_at
	`
	args := []interface{}{
		user.Email,
		user.Username,
		user.FullName,
		user.Phone,
		user.AvatarURL,
		user.IsActive,
		user.IsVerified,
		user.EmailVerifiedAt,
		user.ID,
	}

	updatedUser := &entity.User{}
	err := tx.QueryRow(subCtx, query, args...).Scan(
		&updatedUser.ID,
		&updatedUser.Username,
		&updatedUser.Email,
		&updatedUser.FullName,
		&updatedUser.Phone,
		&updatedUser.AvatarURL,
		&updatedUser.IsActive,
		&updatedUser.IsVerified,
		&updatedUser.EmailVerifiedAt,
		&updatedUser.LastLoginAt,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			switch pgxErr.ConstraintName {
			case "users_username_key":
				return nil, fmt.Errorf("username %s is already taken", user.Username)
			case "users_email_key":
				return nil, fmt.Errorf("email %s is already registered", user.Email)
			case "users_name_length_check":
				return nil, fmt.Errorf("full name %s is invalid, must be between 2 and 50 characters", user.FullName)
			case "users_username_length_check":
				return nil, fmt.Errorf("username %s is invalid, must be between 3 and 20 characters", user.Username)
			default:
				return nil, fmt.Errorf("unique constraint violation (%s): %w", pgxErr.ConstraintName, err)
			}
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return updatedUser, nil
}
func (r *userRepository) UpdateWithRecovery(ctx context.Context, user *entity.User) (*entity.User, error) {
	var updatedUser *entity.User

	err := r.WithTransaction(ctx, func(tx pgx.Tx) error {
		var innerErr error
		updatedUser, innerErr = r.UpdateTx(ctx, tx, user) // ✅ Gunakan ctx utama
		return innerErr
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found or already updated")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}
func (r *userRepository) UpdateTwoFaSecret(ctx context.Context, id uuid.UUID, twoFaSecret *string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE users SET two_fa_secret = $1 WHERE id = $2 AND deleted_at IS NULL`
	args := []interface{}{
		twoFaSecret,
		id,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update two fa secret: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already updated")
	}
	return nil
}
func (r *userRepository) RemoveTwoFaSecret(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE users SET two_fa_secret = NULL WHERE id = $1 AND deleted_at IS NULL`
	args := []interface{}{
		id,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to remove two fa secret: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already updated")
	}
	return nil
}
func (r *userRepository) SetEmailVerified(ctx context.Context, id uuid.UUID, isVerified bool) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE users SET is_verified = $1, email_verified_at = NOW() WHERE id = $2 AND deleted_at IS NULL`
	args := []interface{}{
		isVerified,
		id,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to set email verified: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already updated")
	}
	return nil
}
func (r *userRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE users SET password_hash = $1 WHERE id = $2 AND deleted_at IS NULL`
	args := []interface{}{
		passwordHash,
		id,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already updated")
	}
	return nil
}
func (r *userRepository) UpdateAvatar(ctx context.Context, id uuid.UUID, avatarURL string) (string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE users SET avatar_url = $1 WHERE id = $2 AND deleted_at IS NULL RETURNING avatar_url`
	args := []interface{}{
		avatarURL,
		id,
	}

	var newAvatarURL string
	err := r.db.QueryRow(subCtx, query, args...).Scan(&newAvatarURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("user not found or already updated")
		}
		return "", fmt.Errorf("failed to update avatar: %w", err)
	}
	return newAvatarURL, nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE users SET
			last_login_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{
		id,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already updated")
	}
	return nil
}
func (r *userRepository) UpdateLastLoginTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE users SET
			last_login_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{
		id,
	}
	result, err := tx.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already updated")
	}
	return nil
}
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE users SET
			deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{
		id,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already deleted")
	}
	return nil
}
func (r *userRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		DELETE FROM users WHERE id = $1
	`
	args := []interface{}{
		id,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to hard delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already deleted")
	}
	return nil
}
func (r *userRepository) Restore(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE users SET
			deleted_at = NULL,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NOT NULL
	`
	args := []interface{}{
		id,
	}
	result, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to restore user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already restored")
	}
	return nil
}
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT EXISTS(
			SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL
		)
	`
	args := []interface{}{
		email,
	}
	var exists bool
	err := pgxscan.Get(subCtx, r.db, &exists, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if user exists by email: %w", err)
	}
	return exists, nil
}
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT EXISTS(
			SELECT 1 FROM users WHERE username = $1 AND deleted_at IS NULL
		)
	`
	args := []interface{}{
		username,
	}
	var exists bool
	err := pgxscan.Get(subCtx, r.db, &exists, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if user exists by username: %w", err)
	}
	return exists, nil
}
func (r *userRepository) IsTwoFaEnabled(ctx context.Context, id uuid.UUID) (bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT two_fa_secret IS NOT NULL FROM users WHERE id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{
		id,
	}
	var isEnabled bool
	err := pgxscan.Get(subCtx, r.db, &isEnabled, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if two fa is enabled: %w", err)
	}
	return isEnabled, nil
}
func (r *userRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT EXISTS(
			SELECT 1 FROM users WHERE phone = $1 AND deleted_at IS NULL
		)
	`
	args := []interface{}{
		phone,
	}
	var exists bool
	err := pgxscan.Get(subCtx, r.db, &exists, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if user exists by phone: %w", err)
	}
	return exists, nil
}

// Helpers :
func (r *userRepository) buildBaseQuery(baseQuery string, filter *Filter) *QueryBuilder {
	qb := NewQueryBuilder(baseQuery)

	if filter == nil {
		qb.Where("deleted_at IS NULL")
		return qb
	}

	if filter.IncludeDeleted != nil && *filter.IncludeDeleted {
		qb.Where("deleted_at IS NOT NULL")
	} else {
		qb.Where("deleted_at IS NULL")
	}

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		qb.Where("(email ILIKE $? OR username ILIKE $? OR full_name ILIKE $?)",
			searchPattern, searchPattern, searchPattern)
	}
	if filter.IsActive != nil {
		qb.Where("is_active = $?", *filter.IsActive)
	}
	if filter.TenantID != nil {
		qb.Where("tenant_id = $?", *filter.TenantID)
	}
	if filter.UserID != nil {
		qb.Where("id = $?", *filter.UserID)
	}
	if filter.IsVerified != nil {
		qb.Where("is_verified = $?", *filter.IsVerified)
	}
	if filter.RangeDate != nil {
		var startDate time.Time
		var endDate time.Time

		if !filter.RangeDate.StartDate.IsZero() {
			startDate = filter.RangeDate.StartDate
		} else {
			startDate = time.Now().AddDate(0, 0, -7)
		}
		if !filter.RangeDate.EndDate.IsZero() {
			endDate = filter.RangeDate.EndDate
		} else {
			endDate = time.Now()
		}
		if !startDate.IsZero() || !endDate.IsZero() {
			qb.Where("created_at BETWEEN $? AND $?", startDate, endDate)
		}
	}

	return qb
}
func (r *userRepository) CountUsersByRole(ctx context.Context, tenantID uuid.UUID, roleSlug string, filter *Filter) (int64, error) {
	baseQuery := `
		SELECT COUNT(*)
		FROM users u
		INNER JOIN user_tenants ut ON u.id = ut.user_id
		INNER JOIN roles r ON ut.role_id = r.id
		WHERE ut.tenant_id = $? AND r.slug = $? AND u.deleted_at IS NULL AND ut.deleted_at IS NULL
	`

	qb := NewQueryBuilder(baseQuery)
	qb.Where("ut.tenant_id = $?", tenantID)
	qb.Where("r.slug = $?", roleSlug)
	qb.Where("u.deleted_at IS NULL")
	qb.Where("ut.deleted_at IS NULL")

	if filter != nil {
		if filter.Search != "" {
			searchPattern := "%" + filter.Search + "%"
			qb.Where("(u.full_name ILIKE $? OR u.email ILIKE $? OR u.username ILIKE $?)",
				searchPattern, searchPattern, searchPattern)
		}

		if filter.IsActive != nil {
			qb.Where("u.is_active = $?", *filter.IsActive)
		}

		if filter.IsVerified != nil {
			qb.Where("u.is_verified = $?", *filter.IsVerified)
		}
	}

	query, args := qb.Build()

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users by role %s: %w", roleSlug, err)
	}

	return count, nil
}
