package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	"strings"
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
	// CreateUserTenant links a user to a tenant with a role (membership row).
	CreateUserTenant(ctx context.Context, ut *entity.UserTenant) error
	CreateUserTenantTx(ctx context.Context, tx pgx.Tx, ut *entity.UserTenant) error
	// Tenant member management (supervisors/agents = user_tenants).
	ListTenantMembers(ctx context.Context, tenantID uuid.UUID, roleSlug string) ([]*entity.UserTenantWithDetails, error)
	CountTenantMembersInRoles(ctx context.Context, tenantID uuid.UUID, roleSlugs []string) (int, error)
	FindUserTenantByID(ctx context.Context, id uuid.UUID) (*entity.UserTenant, error)
	UpdateUserTenant(ctx context.Context, id, roleID uuid.UUID, isActive bool) error
	SoftDeleteUserTenant(ctx context.Context, id uuid.UUID) error
	// Read operations
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByEmailOrUsername(ctx context.Context, identifier string) (*entity.User, error)
	// GetUserTenantWithDetailsByUserID loads a user's active membership together
	// with its tenant and role for authentication/session building.
	GetUserTenantWithDetailsByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserTenantWithDetails, error)
	GetUsersByTenantAndRole(ctx context.Context, tenantID uuid.UUID, roleID uuid.UUID) ([]*entity.User, error)
	GetUsersByTenantAndRoleSlug(ctx context.Context, tenantID uuid.UUID, roleSlug string) ([]*entity.User, error)
	Search(ctx context.Context, opts *ListOptions) ([]*entity.User, int64, error)
	Count(ctx context.Context, filter *Filter) (int64, error)
	CountUsersByRole(ctx context.Context, tenantID uuid.UUID, roleSlug string, filter *Filter) (int64, error)
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
			id, full_name, email, password_hash, phone, avatar_url,
			two_fa_secret, status, email_verified_at, last_login_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(subCtx,
		query,
		user.ID,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.Phone,
		user.AvatarURL,
		user.TwoFaSecret,
		user.Status,
		user.EmailVerifiedAt,
		user.LastLoginAt,
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
			case "users_email_key":
				return fmt.Errorf("email %s is already registered", user.Email)
			case "users_name_length_check":
				return fmt.Errorf("full name %s is invalid, must be between 2 and 50 characters", user.FullName)
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
			id, full_name, email, password_hash, phone, avatar_url,
			two_fa_secret, status, email_verified_at, last_login_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	err := tx.QueryRow(subCtx,
		query,
		user.ID,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.Phone,
		user.AvatarURL,
		user.TwoFaSecret,
		user.Status,
		user.EmailVerifiedAt,
		user.LastLoginAt,
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
			case "users_email_key":
				return fmt.Errorf("email %s is already registered", user.Email)
			case "users_name_length_check":
				return fmt.Errorf("full name %s is invalid, must be between 2 and 50 characters", user.FullName)
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
		SELECT
			id, email, password_hash, full_name, phone, avatar_url, two_fa_secret,
			status, email_verified_at, last_login_at, created_at, updated_at, deleted_at
		FROM users WHERE id = $1 AND deleted_at IS NULL
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
		SELECT
			id, email, password_hash, full_name, phone, avatar_url, two_fa_secret,
			status, email_verified_at, last_login_at, created_at, updated_at, deleted_at
		FROM users WHERE email = $1 AND deleted_at IS NULL
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

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT
			id, email, password_hash, full_name, phone, avatar_url, two_fa_secret,
			status, email_verified_at, last_login_at, created_at, updated_at, deleted_at
		FROM users WHERE split_part(email, '@', 1) = $1 AND deleted_at IS NULL
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
		SELECT
			id, email, password_hash, full_name, phone, avatar_url, two_fa_secret,
			status, email_verified_at, last_login_at, created_at, updated_at, deleted_at
		FROM users
		WHERE (email = $1 OR split_part(email, '@', 1) = $1)
		  AND deleted_at IS NULL
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

// CreateUserTenant inserts a membership row linking a user to a tenant + role.
func (r *userRepository) CreateUserTenant(ctx context.Context, ut *entity.UserTenant) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	return r.insertUserTenant(subCtx, r.q(subCtx), ut)
}

// CreateUserTenantTx is the transactional variant of CreateUserTenant.
func (r *userRepository) CreateUserTenantTx(ctx context.Context, tx pgx.Tx, ut *entity.UserTenant) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	return r.insertUserTenant(subCtx, tx, ut)
}

func (r *userRepository) insertUserTenant(ctx context.Context, q Querier, ut *entity.UserTenant) error {
	if ut.ID == uuid.Nil {
		ut.ID = uuid.New()
	}
	query := `
		INSERT INTO user_tenants (id, user_id, tenant_id, role_id, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	err := q.QueryRow(ctx, query, ut.ID, ut.UserID, ut.TenantID, ut.RoleID, ut.IsActive).
		Scan(&ut.ID, &ut.CreatedAt, &ut.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user tenant: %w", err)
	}
	return nil
}

// ListTenantMembers returns the tenant's memberships (user + role). Pass an
// empty roleSlug for all roles. Uses users.status (no is_active column).
func (r *userRepository) ListTenantMembers(ctx context.Context, tenantID uuid.UUID, roleSlug string) ([]*entity.UserTenantWithDetails, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT
			u.id, u.email, u.full_name, u.phone, u.avatar_url, u.status,
			u.email_verified_at, u.last_login_at, u.created_at, u.updated_at,
			ut.id, ut.user_id, ut.tenant_id, ut.role_id, ut.is_active,
			ut.created_at, ut.updated_at,
			r.id, r.name, r.slug, r.level
		FROM user_tenants ut
		JOIN users u ON u.id = ut.user_id AND u.deleted_at IS NULL
		JOIN roles r ON r.id = ut.role_id
		WHERE ut.tenant_id = $1 AND ut.deleted_at IS NULL
			AND ($2 = '' OR r.slug = $2)
		ORDER BY r.level ASC, u.full_name ASC`

	rows, err := r.q(subCtx).Query(subCtx, query, tenantID, roleSlug)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenant members: %w", err)
	}
	defer rows.Close()

	members := make([]*entity.UserTenantWithDetails, 0)
	for rows.Next() {
		var d entity.UserTenantWithDetails
		if err := rows.Scan(
			&d.User.ID, &d.User.Email, &d.User.FullName, &d.User.Phone, &d.User.AvatarURL,
			&d.User.Status, &d.User.EmailVerifiedAt, &d.User.LastLoginAt,
			&d.User.CreatedAt, &d.User.UpdatedAt,
			&d.UserTenant.ID, &d.UserTenant.UserID, &d.UserTenant.TenantID, &d.UserTenant.RoleID,
			&d.UserTenant.IsActive, &d.UserTenant.CreatedAt, &d.UserTenant.UpdatedAt,
			&d.Role.ID, &d.Role.Name, &d.Role.Slug, &d.Role.Level,
		); err != nil {
			return nil, fmt.Errorf("failed to scan tenant member: %w", err)
		}
		members = append(members, &d)
	}
	return members, rows.Err()
}

// CountTenantMembersInRoles counts active memberships whose role slug is in the
// given set (used for plan quota, e.g. agents = supervisor + agent).
func (r *userRepository) CountTenantMembersInRoles(ctx context.Context, tenantID uuid.UUID, roleSlugs []string) (int, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT COUNT(*)
		FROM user_tenants ut
		JOIN roles r ON r.id = ut.role_id
		WHERE ut.tenant_id = $1 AND ut.deleted_at IS NULL AND ut.is_active = true
			AND r.slug = ANY($2)`

	var count int
	if err := r.q(subCtx).QueryRow(subCtx, query, tenantID, roleSlugs).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count tenant members: %w", err)
	}
	return count, nil
}

func (r *userRepository) FindUserTenantByID(ctx context.Context, id uuid.UUID) (*entity.UserTenant, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, tenant_id, role_id, is_active, created_at, updated_at, deleted_at
		FROM user_tenants WHERE id = $1 AND deleted_at IS NULL`

	var ut entity.UserTenant
	if err := pgxscan.Get(subCtx, r.q(subCtx), &ut, query, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("membership not found")
		}
		return nil, fmt.Errorf("failed to find membership: %w", err)
	}
	return &ut, nil
}

func (r *userRepository) UpdateUserTenant(ctx context.Context, id, roleID uuid.UUID, isActive bool) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE user_tenants SET role_id = $1, is_active = $2, updated_at = now()
		WHERE id = $3 AND deleted_at IS NULL`
	tag, err := r.q(subCtx).Exec(subCtx, query, roleID, isActive, id)
	if err != nil {
		return fmt.Errorf("failed to update membership: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("membership not found")
	}
	return nil
}

func (r *userRepository) SoftDeleteUserTenant(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE user_tenants SET deleted_at = now(), is_active = false, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`
	tag, err := r.q(subCtx).Exec(subCtx, query, id)
	if err != nil {
		return fmt.Errorf("failed to remove membership: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("membership not found")
	}
	return nil
}

// GetUserTenantWithDetailsByUserID loads the user's active tenant membership
// with tenant + role. Uses users.status (no is_active column) and picks the
// earliest active membership as the primary one.
func (r *userRepository) GetUserTenantWithDetailsByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserTenantWithDetails, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT
			u.id, u.email, u.password_hash, u.full_name, u.phone, u.avatar_url,
			u.two_fa_secret, u.status, u.email_verified_at, u.last_login_at,
			u.created_at, u.updated_at,
			ut.id, ut.user_id, ut.tenant_id, ut.role_id, ut.is_active,
			ut.created_at, ut.updated_at,
			t.id, t.name, t.slug, t.subscription_plan, t.subscription_status,
			t.is_active, t.trial_ends_at,
			r.id, r.name, r.slug, r.level
		FROM users u
		JOIN user_tenants ut ON ut.user_id = u.id
			AND ut.deleted_at IS NULL AND ut.is_active = true
		JOIN tenants t ON t.id = ut.tenant_id AND t.deleted_at IS NULL
		JOIN roles r ON r.id = ut.role_id
		WHERE u.id = $1 AND u.deleted_at IS NULL
		ORDER BY ut.created_at ASC
		LIMIT 1`

	var d entity.UserTenantWithDetails
	err := r.q(subCtx).QueryRow(subCtx, query, userID).Scan(
		&d.User.ID, &d.User.Email, &d.User.PasswordHash, &d.User.FullName, &d.User.Phone,
		&d.User.AvatarURL, &d.User.TwoFaSecret, &d.User.Status, &d.User.EmailVerifiedAt,
		&d.User.LastLoginAt, &d.User.CreatedAt, &d.User.UpdatedAt,
		&d.UserTenant.ID, &d.UserTenant.UserID, &d.UserTenant.TenantID, &d.UserTenant.RoleID,
		&d.UserTenant.IsActive, &d.UserTenant.CreatedAt, &d.UserTenant.UpdatedAt,
		&d.Tenant.ID, &d.Tenant.Name, &d.Tenant.Slug, &d.Tenant.SubscriptionPlan,
		&d.Tenant.SubscriptionStatus, &d.Tenant.IsActive, &d.Tenant.TrialEndsAt,
		&d.Role.ID, &d.Role.Name, &d.Role.Slug, &d.Role.Level,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no active tenant membership found for user %s: %w", userID, err)
		}
		return nil, fmt.Errorf("failed to get user tenant details: %w", err)
	}

	return &d, nil
}
func (r *userRepository) GetUsersByTenantAndRole(ctx context.Context, tenantID, roleID uuid.UUID) ([]*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
        WITH tenant_users AS (
            SELECT 
                ut.id as user_tenant_id,
                ut.user_id, ut.role_id, ut.is_active as ut_active,
                ut.created_at as ut_created, ut.updated_at as ut_updated,
                u.id, u.email, u.password_hash, u.full_name, u.phone, 
                u.avatar_url, u.two_fa_secret, u.status, u.is_active,
                u.email_verified_at, u.last_login_at,
                u.created_at, u.updated_at, u.deleted_at,
                r.id as role_id, r.name as role_name, r.slug as role_slug, r.level
            FROM user_tenants ut
            JOIN users u ON ut.user_id = u.id AND u.deleted_at IS NULL
            JOIN roles r ON ut.role_id = r.id
            WHERE ut.tenant_id = $1 
              AND ut.role_id = $2
              AND ut.deleted_at IS NULL 
              AND ut.is_active = true
            ORDER BY u.full_name ASC
        ),
        division_aggregate AS (
            SELECT 
                dm.user_tenant_id,
                json_agg(
                    json_build_object(
                        'id', d.id,
                        'name', d.name,
                        'slug', d.slug,
                        'routing_type', d.routing_type,
                        'is_active', d.is_active,
                        'link_url', d.link_url
                    )
                    ORDER BY d.name ASC
                ) as divisions
            FROM division_members dm
            JOIN divisions d ON dm.division_id = d.id AND d.deleted_at IS NULL
            WHERE dm.deleted_at IS NULL AND dm.is_active = true
            GROUP BY dm.user_tenant_id
        )
        SELECT 
            json_agg(
                json_build_object(
                    'id', tu.id,
                    'email', tu.email,
                    'full_name', tu.full_name,
                    'phone', tu.phone,
                    'avatar_url', tu.avatar_url,
                    'status', tu.status,
                    'is_active', tu.is_active,
                    'email_verified_at', tu.email_verified_at,
                    'last_login_at', tu.last_login_at,
                    'created_at', tu.created_at,
                    'updated_at', tu.updated_at,
                    'role', json_build_object(
                        'id', tu.role_id,
                        'name', tu.role_name,
                        'slug', tu.role_slug,
                        'level', tu.level
                    ),
                    'divisions', COALESCE(da.divisions, '[]'::json)
                )
                ORDER BY tu.full_name ASC
            ) as users
        FROM tenant_users tu
        LEFT JOIN division_aggregate da ON da.user_tenant_id = tu.user_tenant_id
    `

	var result struct {
		Users json.RawMessage `db:"users"`
	}

	err := pgxscan.Get(subCtx, r.db, &result, query, tenantID, roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.User{}, nil
		}
		return nil, fmt.Errorf("failed to get users by tenant and role: %w", err)
	}

	var users []*entity.User
	if err := json.Unmarshal(result.Users, &users); err != nil {
		return nil, fmt.Errorf("failed to parse users: %w", err)
	}

	return users, nil
}
func (r *userRepository) GetUsersByTenantAndRoleSlug(ctx context.Context, tenantID uuid.UUID, roleSlug string) ([]*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
        WITH tenant_users AS (
            SELECT 
                ut.id as user_tenant_id,
                ut.user_id, ut.role_id, ut.is_active as ut_active,
                ut.created_at as ut_created, ut.updated_at as ut_updated,
                u.id, u.email, u.password_hash, u.full_name, u.phone, 
                u.avatar_url, u.two_fa_secret, u.status, u.is_active,
                u.email_verified_at, u.last_login_at,
                u.created_at, u.updated_at, u.deleted_at,
                r.id as role_id, r.name as role_name, r.slug as role_slug, r.level
            FROM user_tenants ut
            JOIN users u ON ut.user_id = u.id AND u.deleted_at IS NULL
            JOIN roles r ON ut.role_id = r.id
            WHERE ut.tenant_id = $1 
              AND r.slug = $2
              AND ut.deleted_at IS NULL 
              AND ut.is_active = true
            ORDER BY u.full_name ASC
        ),
        division_aggregate AS (
            SELECT 
                dm.user_tenant_id,
                json_agg(
                    json_build_object(
                        'id', d.id,
                        'name', d.name,
                        'slug', d.slug,
                        'routing_type', d.routing_type,
                        'is_active', d.is_active,
                        'link_url', d.link_url
                    )
                    ORDER BY d.name ASC
                ) as divisions
            FROM division_members dm
            JOIN divisions d ON dm.division_id = d.id AND d.deleted_at IS NULL
            WHERE dm.deleted_at IS NULL AND dm.is_active = true
            GROUP BY dm.user_tenant_id
        )
        SELECT 
            json_agg(
                json_build_object(
                    'id', tu.id,
                    'email', tu.email,
                    'full_name', tu.full_name,
                    'phone', tu.phone,
                    'avatar_url', tu.avatar_url,
                    'status', tu.status,
                    'is_active', tu.is_active,
                    'email_verified_at', tu.email_verified_at,
                    'last_login_at', tu.last_login_at,
                    'created_at', tu.created_at,
                    'updated_at', tu.updated_at,
                    'role', json_build_object(
                        'id', tu.role_id,
                        'name', tu.role_name,
                        'slug', tu.role_slug,
                        'level', tu.level
                    ),
                    'divisions', COALESCE(da.divisions, '[]'::json)
                )
                ORDER BY tu.full_name ASC
            ) as users
        FROM tenant_users tu
        LEFT JOIN division_aggregate da ON da.user_tenant_id = tu.user_tenant_id
    `

	var result struct {
		Users json.RawMessage `db:"users"`
	}

	err := pgxscan.Get(subCtx, r.db, &result, query, tenantID, roleSlug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.User{}, nil
		}
		return nil, fmt.Errorf("failed to get users by tenant and role slug: %w", err)
	}

	var users []*entity.User
	if err := json.Unmarshal(result.Users, &users); err != nil {
		return nil, fmt.Errorf("failed to parse users: %w", err)
	}

	return users, nil
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
func (r *userRepository) CountUsersByRole(ctx context.Context, tenantID uuid.UUID, roleSlug string, filter *Filter) (int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
        SELECT COUNT(DISTINCT u.id)
        FROM users u
        JOIN user_tenants ut ON u.id = ut.user_id
        JOIN roles r ON ut.role_id = r.id
        WHERE ut.tenant_id = $1 
          AND r.slug = $2 
          AND ut.deleted_at IS NULL 
          AND ut.is_active = true
          AND u.deleted_at IS NULL
    `

	args := []interface{}{tenantID, roleSlug}
	argCount := 3

	if filter != nil {
		if filter.Search != "" {
			searchPattern := "%" + filter.Search + "%"
			query += fmt.Sprintf(" AND (u.full_name ILIKE $%d OR u.email ILIKE $%d)", argCount, argCount+1)
			args = append(args, searchPattern, searchPattern)
			argCount += 2
		}
		if filter.IsActive != nil {
			query += fmt.Sprintf(" AND u.is_active = $%d", argCount)
			args = append(args, *filter.IsActive)
			argCount++
		}
	}

	var count int64
	err := r.db.QueryRow(subCtx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users by role: %w", err)
	}
	return count, nil
}
func (r *userRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
    UPDATE users SET
        email = $1,
        full_name = $2,
        phone = $3,
        avatar_url = $4,
        status = $5,
        email_verified_at = $6,
        last_login_at = $7
    WHERE id = $8 AND deleted_at IS NULL
    RETURNING id, email, full_name, phone, avatar_url, two_fa_secret, status, email_verified_at, last_login_at, created_at, updated_at, deleted_at
    `

	args := []interface{}{
		user.Email,
		user.FullName,
		user.Phone,
		user.AvatarURL,
		user.Status,
		user.EmailVerifiedAt,
		user.LastLoginAt,
		user.ID,
	}

	updatedUser := &entity.User{}
	err := r.db.QueryRow(
		subCtx,
		query,
		args...).Scan(
		&updatedUser.ID,
		&updatedUser.Email,
		&updatedUser.FullName,
		&updatedUser.Phone,
		&updatedUser.AvatarURL,
		&updatedUser.TwoFaSecret,
		&updatedUser.Status,
		&updatedUser.EmailVerifiedAt,
		&updatedUser.LastLoginAt,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
		&updatedUser.DeletedAt,
	)

	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			switch pgxErr.ConstraintName {
			case "users_email_key":
				return nil, fmt.Errorf("email %s is already registered", user.Email)
			case "users_name_length_check":
				return nil, fmt.Errorf("full name %s is invalid, must be between 2 and 50 characters", user.FullName)
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
        full_name = $2,
        phone = $3,
        avatar_url = $4,
        status = $5,
        email_verified_at = $6,
        last_login_at = $7
    WHERE id = $8 AND deleted_at IS NULL
    RETURNING id, email, full_name, phone, avatar_url, two_fa_secret, status, email_verified_at, last_login_at, created_at, updated_at, deleted_at
    `
	args := []interface{}{
		user.Email,
		user.FullName,
		user.Phone,
		user.AvatarURL,
		user.Status,
		user.EmailVerifiedAt,
		user.LastLoginAt,
		user.ID,
	}

	updatedUser := &entity.User{}
	err := tx.QueryRow(subCtx, query, args...).Scan(
		&updatedUser.ID,
		&updatedUser.Email,
		&updatedUser.FullName,
		&updatedUser.Phone,
		&updatedUser.AvatarURL,
		&updatedUser.TwoFaSecret,
		&updatedUser.Status,
		&updatedUser.EmailVerifiedAt,
		&updatedUser.LastLoginAt,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
		&updatedUser.DeletedAt,
	)

	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			switch pgxErr.ConstraintName {
			case "users_email_key":
				return nil, fmt.Errorf("email %s is already registered", user.Email)
			case "users_name_length_check":
				return nil, fmt.Errorf("full name %s is invalid, must be between 2 and 50 characters", user.FullName)
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

	query := `
		UPDATE users
		SET email_verified_at = NOW(),
			status = CASE WHEN $1 THEN 'active' ELSE status END
		WHERE id = $2 AND deleted_at IS NULL
	`
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
			SELECT 1 FROM users WHERE split_part(email, '@', 1) = $1 AND deleted_at IS NULL
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

		qb.Where(`(
            -- User fields
            email ILIKE $? 
            OR split_part(email, '@', 1) ILIKE $? 
            OR full_name ILIKE $?
            OR phone ILIKE $?
            -- Tenant fields (through user_tenants)
            OR EXISTS (
                SELECT 1 FROM user_tenants ut
                JOIN tenants t ON ut.tenant_id = t.id
                WHERE ut.user_id = users.id 
                  AND ut.deleted_at IS NULL
                  AND (t.name ILIKE $? OR t.slug ILIKE $?)
            )
            -- Division fields (through division_members)
            OR EXISTS (
                SELECT 1 FROM user_tenants ut
                JOIN division_members dm ON dm.user_tenant_id = ut.id
                JOIN divisions d ON dm.division_id = d.id
                WHERE ut.user_id = users.id 
                  AND ut.deleted_at IS NULL
                  AND dm.deleted_at IS NULL
                  AND (d.name ILIKE $? OR d.slug ILIKE $?)
            )
        )`,
			searchPattern, // email
			searchPattern, // username (split email)
			searchPattern, // full_name
			searchPattern, // phone
			searchPattern, // tenant name
			searchPattern, // tenant slug
			searchPattern, // division name
			searchPattern, // division slug
		)
	}
	if filter.Status != "" {
		qb.Where("status = $?", filter.Status)
	}
	if filter.TenantID != nil {
		qb.Where("EXISTS (SELECT 1 FROM user_tenants ut WHERE ut.user_id = users.id AND ut.tenant_id = $? AND ut.deleted_at IS NULL)", *filter.TenantID)
	}
	if filter.UserID != nil {
		qb.Where("id = $?", *filter.UserID)
	}
	if filter.IsVerified != nil {
		if *filter.IsVerified {
			qb.Where("email_verified_at IS NOT NULL")
		} else {
			qb.Where("email_verified_at IS NULL")
		}
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

	if filter.Extra != nil {
		if divisionID, ok := filter.Extra["division_id"].(uuid.UUID); ok {
			qb.Where(`EXISTS (
                SELECT 1 FROM user_tenants ut
                JOIN division_members dm ON dm.user_tenant_id = ut.id
                WHERE ut.user_id = users.id 
                  AND dm.division_id = $?
                  AND dm.deleted_at IS NULL
                  AND dm.is_active = true
            )`, divisionID)
		}

		// Filter by multiple division IDs
		if divisionIDs, ok := filter.Extra["division_ids"].([]uuid.UUID); ok && len(divisionIDs) > 0 {
			placeholders := make([]string, len(divisionIDs))
			args := make([]interface{}, len(divisionIDs))
			for i, id := range divisionIDs {
				placeholders[i] = "$?"
				args[i] = id
			}
			qb.Where(fmt.Sprintf(`EXISTS (
                SELECT 1 FROM user_tenants ut
                JOIN division_members dm ON dm.user_tenant_id = ut.id
                WHERE ut.user_id = users.id 
                  AND dm.division_id IN (%s)
                  AND dm.deleted_at IS NULL
                  AND dm.is_active = true
            )`, strings.Join(placeholders, ", ")), args...)
		}

		if roleID, ok := filter.Extra["role_id"].(uuid.UUID); ok {
			qb.Where(`EXISTS (
                SELECT 1 FROM user_tenants ut
                WHERE ut.user_id = users.id 
                  AND ut.role_id = $?
                  AND ut.deleted_at IS NULL
                  AND ut.is_active = true
            )`, roleID)
		}

		// Filter by role slug
		if roleSlug, ok := filter.Extra["role_slug"].(string); ok {
			qb.Where(`EXISTS (
                SELECT 1 FROM user_tenants ut
                JOIN roles r ON ut.role_id = r.id
                WHERE ut.user_id = users.id 
                  AND r.slug = $?
                  AND ut.deleted_at IS NULL
                  AND ut.is_active = true
            )`, roleSlug)
		}

		// Filter users with NO division
		if noDivision, ok := filter.Extra["no_division"].(bool); ok && noDivision {
			qb.Where(`NOT EXISTS (
                SELECT 1 FROM user_tenants ut
                JOIN division_members dm ON dm.user_tenant_id = ut.id
                WHERE ut.user_id = users.id 
                  AND dm.deleted_at IS NULL
                  AND dm.is_active = true
            )`)
		}

		// Filter by multiple tenant IDs
		if tenantIDs, ok := filter.Extra["tenant_ids"].([]uuid.UUID); ok && len(tenantIDs) > 0 {
			placeholders := make([]string, len(tenantIDs))
			args := make([]interface{}, len(tenantIDs))
			for i, id := range tenantIDs {
				placeholders[i] = "$?"
				args[i] = id
			}
			qb.Where(fmt.Sprintf(`EXISTS (
								SELECT 1 FROM user_tenants ut
								WHERE ut.user_id = users.id 
									AND ut.tenant_id IN (%s)
									AND ut.deleted_at IS NULL
						)`, strings.Join(placeholders, ", ")), args...)
		}
	}

	return qb
}
