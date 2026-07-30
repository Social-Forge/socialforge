package repository

import (
	"context"
	"errors"
	"fmt"
	"github/socialforge/config"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type TenantRepository interface {
	BaseRepository
	Create(ctx context.Context, tenant *entity.Tenant) error
	CreateTx(ctx context.Context, tx pgx.Tx, tenant *entity.Tenant) error
	CreateWithRecovery(ctx context.Context, tenant *entity.Tenant) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Tenant, error)
	FindByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*entity.Tenant, error)
	FindByUserTenantID(ctx context.Context, userTenantID uuid.UUID) (*entity.Tenant, error)
	Search(ctx context.Context, opts *ListOptions) ([]*entity.Tenant, int64, error)
	Count(ctx context.Context, filter *Filter) (int64, error)
	Update(ctx context.Context, tenant *entity.Tenant) (*entity.Tenant, error)
	UpdateTx(ctx context.Context, tx pgx.Tx, tenant *entity.Tenant) (*entity.Tenant, error)
	UpdateWithRecovery(ctx context.Context, tenant *entity.Tenant) (*entity.Tenant, error)
	UpdateLogo(ctx context.Context, tenantID uuid.UUID, logoURL string) (string, error)
	Delete(ctx context.Context, id uuid.UUID) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
	GetAllowedTenantIDs(ctx context.Context) ([]uuid.UUID, error)
	IsAllowedTenant(ctx context.Context, tenantID uuid.UUID) (bool, error)
}

type tenantRepository struct {
	*baseRepository
}

func NewTenantRepository(db *pgxpool.Pool) TenantRepository {
	return &tenantRepository{
		baseRepository: NewBaseRepository(db).(*baseRepository),
	}
}
func (r *tenantRepository) Create(ctx context.Context, tenant *entity.Tenant) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		INSERT INTO tenants (
			id, name, slug, owner_id, logo_url, description,
			max_divisions, max_agents, max_quick_replies, max_pages,
			max_whatsapp, max_meta_whatsapp, max_meta_messenger, max_instagram, max_telegram, max_webchat, max_linkchat,
			subscription_plan, subscription_status, trial_ends_at, is_active, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
		) ON CONFLICT (slug) DO NOTHING RETURNING id, created_at, updated_at`

	args := []interface{}{
		tenant.ID,
		tenant.Name,
		tenant.Slug,
		tenant.OwnerID,
		tenant.LogoURL,
		tenant.Description,
		tenant.MaxDivisions,
		tenant.MaxAgents,
		tenant.MaxQuickReplies,
		tenant.MaxPages,
		tenant.MaxWhatsApp,
		tenant.MaxMetaWhatsApp,
		tenant.MaxMetaMessenger,
		tenant.MaxInstagram,
		tenant.MaxTelegram,
		tenant.MaxWebChat,
		tenant.MaxLinkChat,
		tenant.SubscriptionPlan,
		tenant.SubscriptionStatus,
		tenant.TrialEndsAt,
		tenant.IsActive,
		tenant.CreatedAt,
	}

	err := r.db.QueryRow(
		subCtx,
		query,
		args...).Scan(
		&tenant.ID,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			switch pgxErr.ConstraintName {
			case "tenants_slug_key":
				return fmt.Errorf("tenant slug '%s' is already taken", tenant.Slug)
			default:
				return fmt.Errorf("unique constraint violation (%s): %w", pgxErr.ConstraintName, err)
			}
		}
		return fmt.Errorf("failed to create tenant: %w", err)
	}
	return nil
}
func (r *tenantRepository) CreateTx(ctx context.Context, tx pgx.Tx, tenant *entity.Tenant) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		INSERT INTO tenants (
			id, name, slug, owner_id, logo_url, description,
			max_divisions, max_agents, max_quick_replies, max_pages,
			max_whatsapp, max_meta_whatsapp, max_meta_messenger, max_instagram, max_telegram, max_webchat, max_linkchat,
			subscription_plan, subscription_status, trial_ends_at, is_active, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
		) ON CONFLICT (slug) DO NOTHING RETURNING id, created_at, updated_at`

	args := []interface{}{
		tenant.ID,
		tenant.Name,
		tenant.Slug,
		tenant.OwnerID,
		tenant.LogoURL,
		tenant.Description,
		tenant.MaxDivisions,
		tenant.MaxAgents,
		tenant.MaxQuickReplies,
		tenant.MaxPages,
		tenant.MaxWhatsApp,
		tenant.MaxMetaWhatsApp,
		tenant.MaxMetaMessenger,
		tenant.MaxInstagram,
		tenant.MaxTelegram,
		tenant.MaxWebChat,
		tenant.MaxLinkChat,
		tenant.SubscriptionPlan,
		tenant.SubscriptionStatus,
		tenant.TrialEndsAt,
		tenant.IsActive,
		tenant.CreatedAt,
	}
	err := tx.QueryRow(subCtx, query, args...).Scan(
		&tenant.ID,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			switch pgxErr.ConstraintName {
			case "tenants_slug_key":
				return fmt.Errorf("tenant slug '%s' is already taken", tenant.Slug)
			default:
				return fmt.Errorf("unique constraint violation (%s): %w", pgxErr.ConstraintName, err)
			}
		}
		return fmt.Errorf("failed to create tenant: %w", err)
	}
	return nil
}
func (r *tenantRepository) CreateWithRecovery(ctx context.Context, tenant *entity.Tenant) error {
	return r.WithTransaction(ctx, func(tx pgx.Tx) error {
		return r.CreateTx(ctx, tx, tenant)
	})
}
func (r *tenantRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM tenants WHERE id = $1 AND deleted_at IS NULL
	`
	var tenant entity.Tenant
	err := pgxscan.Get(subCtx, r.db, &tenant, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("tenant not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find tenant by id: %w", err)
	}
	return &tenant, nil
}
func (r *tenantRepository) FindBySlug(ctx context.Context, slug string) (*entity.Tenant, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if slug == "" {
		return nil, fmt.Errorf("slug is required")
	}
	query := `
		SELECT * FROM tenants WHERE slug = $1 AND deleted_at IS NULL
	`
	var tenant entity.Tenant
	err := pgxscan.Get(subCtx, r.db, &tenant, query, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("tenant not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find tenant by slug: %w", err)
	}
	return &tenant, nil
}
func (r *tenantRepository) FindByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*entity.Tenant, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if ownerID == uuid.Nil {
		return nil, fmt.Errorf("owner id is required")
	}

	query := `
		SELECT * FROM tenants WHERE owner_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC
	`
	var tenants []*entity.Tenant
	err := pgxscan.Select(subCtx, r.db, &tenants, query, ownerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("tenant not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find tenant by owner id: %w", err)
	}
	return tenants, nil
}
func (r *tenantRepository) FindByUserTenantID(ctx context.Context, userTenantID uuid.UUID) (*entity.Tenant, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if userTenantID == uuid.Nil {
		return nil, fmt.Errorf("user tenant id is required")
	}

	query := `
		SELECT * FROM tenants WHERE user_tenant_id = $1 AND deleted_at IS NULL
	`
	var tenant entity.Tenant
	err := pgxscan.Get(subCtx, r.db, &tenant, query, userTenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("tenant not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find tenant by user tenant id: %w", err)
	}
	return &tenant, nil
}

func (r *tenantRepository) Search(ctx context.Context, opts *ListOptions) ([]*entity.Tenant, int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if opts == nil {
		opts = &ListOptions{}
	}

	totalRows, err := r.Count(ctx, opts.Filter)
	if err != nil {
		return nil, 0, err
	}

	qb := r.buildBaseQuery("SELECT * FROM tenants", opts.Filter)

	// Add ordering & pagination
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

	var tenants []*entity.Tenant
	err = pgxscan.Select(subCtx, r.db, &tenants, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, fmt.Errorf("tenant not found: %w", err)
		}
		return nil, 0, fmt.Errorf("failed to get list tenant: %w", err)
	}
	return tenants, totalRows, nil
}
func (r *tenantRepository) Count(ctx context.Context, filter *Filter) (int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	qb := r.buildBaseQuery("SELECT COUNT(*) FROM tenants", filter)
	query, args := qb.Build()

	var count int64
	err := r.db.QueryRow(subCtx, query, args...).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("tenant not found: %w", err)
		}
		return 0, fmt.Errorf("failed to count tenant: %w", err)
	}
	return count, nil
}
func (r *tenantRepository) Update(ctx context.Context, tenant *entity.Tenant) (*entity.Tenant, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE tenants SET
			name = $1,
			slug = $2,
			logo_url = $3,
			description = $4,
			subscription_plan = $5,
			subscription_status = $6,
			trial_ends_at = $7,
			is_active = $8
		WHERE id = $9 AND deleted_at IS NULL
		RETURNING id, slug, name, owner_id, logo_url, description, 
		max_divisions, max_agents, max_quick_replies, max_pages, max_whatsapp,
		max_meta_whatsapp, max_meta_messenger, max_instagram, max_telegram,
		max_webchat, max_linkchat, subscription_plan, subscription_status, trial_ends_at,
		is_active, created_at, updated_at
	`
	args := []interface{}{
		tenant.Name,
		tenant.Slug,
		tenant.LogoURL,
		tenant.Description,
		tenant.SubscriptionPlan,
		tenant.SubscriptionStatus,
		tenant.TrialEndsAt,
		tenant.IsActive,
		tenant.ID,
	}

	updateTenant := &entity.Tenant{}
	err := r.db.QueryRow(subCtx, query, args...).Scan(
		&updateTenant.ID,
		&updateTenant.Slug,
		&updateTenant.Name,
		&updateTenant.OwnerID,
		&updateTenant.LogoURL,
		&updateTenant.Description,
		&updateTenant.MaxDivisions,
		&updateTenant.MaxAgents,
		&updateTenant.MaxQuickReplies,
		&updateTenant.MaxPages,
		&updateTenant.MaxWhatsApp,
		&updateTenant.MaxMetaWhatsApp,
		&updateTenant.MaxMetaMessenger,
		&updateTenant.MaxInstagram,
		&updateTenant.MaxTelegram,
		&updateTenant.MaxWebChat,
		&updateTenant.MaxLinkChat,
		&updateTenant.SubscriptionPlan,
		&updateTenant.SubscriptionStatus,
		&updateTenant.TrialEndsAt,
		&updateTenant.IsActive,
		&updateTenant.CreatedAt,
		&updateTenant.UpdatedAt,
	)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			switch pgxErr.ConstraintName {
			case "tenants_slug_key":
				return nil, fmt.Errorf("tenant slug '%s' is already taken", tenant.Slug)
			default:
				return nil, fmt.Errorf("unique constraint violation (%s): %w", pgxErr.ConstraintName, err)
			}
		}
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}
	return updateTenant, nil
}
func (r *tenantRepository) UpdateTx(ctx context.Context, tx pgx.Tx, tenant *entity.Tenant) (*entity.Tenant, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE tenants SET
			name = $1,
			slug = $2,
			logo_url = $3,
			description = $4,
			subscription_plan = $5,
			subscription_status = $6,
			trial_ends_at = $7,
			is_active = $8
		WHERE id = $9 AND deleted_at IS NULL
		RETURNING id, slug, name, owner_id, logo_url, description, 
		max_divisions, max_agents, max_quick_replies, max_pages, max_whatsapp,
		max_meta_whatsapp, max_meta_messenger, max_instagram, max_telegram,
		max_webchat, max_linkchat, subscription_plan, subscription_status, trial_ends_at,
		is_active, created_at, updated_at
	`
	args := []interface{}{
		tenant.Name,
		tenant.Slug,
		tenant.LogoURL,
		tenant.Description,
		tenant.SubscriptionPlan,
		tenant.SubscriptionStatus,
		tenant.TrialEndsAt,
		tenant.IsActive,
		tenant.ID,
	}

	updateTenant := &entity.Tenant{}
	err := tx.QueryRow(subCtx, query, args...).Scan(
		&updateTenant.ID,
		&updateTenant.Slug,
		&updateTenant.Name,
		&updateTenant.OwnerID,
		&updateTenant.LogoURL,
		&updateTenant.Description,
		&updateTenant.MaxDivisions,
		&updateTenant.MaxAgents,
		&updateTenant.MaxQuickReplies,
		&updateTenant.MaxPages,
		&updateTenant.MaxWhatsApp,
		&updateTenant.MaxMetaWhatsApp,
		&updateTenant.MaxMetaMessenger,
		&updateTenant.MaxInstagram,
		&updateTenant.MaxTelegram,
		&updateTenant.MaxWebChat,
		&updateTenant.MaxLinkChat,
		&updateTenant.SubscriptionPlan,
		&updateTenant.SubscriptionStatus,
		&updateTenant.TrialEndsAt,
		&updateTenant.IsActive,
		&updateTenant.CreatedAt,
		&updateTenant.UpdatedAt,
	)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			switch pgxErr.ConstraintName {
			case "tenants_slug_key":
				return nil, fmt.Errorf("tenant slug '%s' is already taken", tenant.Slug)
			default:
				return nil, fmt.Errorf("unique constraint violation (%s): %w", pgxErr.ConstraintName, err)
			}
		}
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}
	return updateTenant, nil
}
func (r *tenantRepository) UpdateWithRecovery(ctx context.Context, tenant *entity.Tenant) (*entity.Tenant, error) {
	var updateTenant *entity.Tenant

	err := r.WithTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		updateTenant, err = r.UpdateTx(ctx, tx, tenant)
		return err
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("tenant not found or already updated")
		}
		return nil, fmt.Errorf("failed to update tenant with recovery: %w", err)
	}
	return updateTenant, nil
}
func (r *tenantRepository) UpdateLogo(ctx context.Context, tenantID uuid.UUID, logoURL string) (string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE tenants SET
			logo_url = $1
		WHERE id = $2 AND deleted_at IS NULL AND is_active = true
		RETURNING logo_url
	`
	args := []interface{}{
		logoURL,
		tenantID,
	}

	var newLogoURL string
	err := r.db.QueryRow(subCtx, query, args...).Scan(&newLogoURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("tenant not found or not active")
		}
		return "", fmt.Errorf("failed to update logo: %w", err)
	}

	return newLogoURL, nil
}
func (r *tenantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE tenants SET
			deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{
		id,
	}

	cmdTag, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("tenant not found or already deleted")
	}

	return nil
}
func (r *tenantRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		DELETE FROM tenants WHERE id = $1
	`
	args := []interface{}{
		id,
	}

	cmdTag, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to hard delete tenant: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("tenant not found or already deleted")
	}

	return nil
}
func (r *tenantRepository) Restore(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE tenants SET
			deleted_at = NULL
		WHERE id = $1 AND deleted_at IS NOT NULL
	`
	args := []interface{}{
		id,
	}

	cmdTag, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to restore tenant: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("tenant not found or already restored")
	}

	return nil
}
func (r *tenantRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT EXISTS(
			SELECT 1 FROM tenants WHERE slug = $1 AND deleted_at IS NULL
		)
	`
	args := []interface{}{
		slug,
	}

	var exists bool
	err := r.db.QueryRow(subCtx, query, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check tenant existence by slug: %w", err)
	}

	return exists, nil
}
func (r *tenantRepository) buildBaseQuery(baseQuery string, filter *Filter) *QueryBuilder {
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
		qb.Where("(name ILIKE $? OR slug ILIKE $?)", searchPattern, searchPattern)
	}
	if filter.Status != "" {
		qb.Where("subscription_status = $?", filter.Status)
	}

	return qb
}
func (r *tenantRepository) GetAllowedTenantIDs(ctx context.Context) ([]uuid.UUID, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT id FROM tenants WHERE subscription_status = 'active' AND deleted_at IS NULL
	`

	var tenantIDs []uuid.UUID
	err := pgxscan.Select(subCtx, r.db, &tenantIDs, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []uuid.UUID{}, nil
		}
		if errors.Is(err, subCtx.Err()) {
			if errors.Is(err, context.DeadlineExceeded) {
				config.Logger.Error("context deadline exceeded", zap.Error(err))
			} else if errors.Is(err, context.Canceled) {
				config.Logger.Error("context canceled", zap.Error(err))
			}
		}
		return nil, fmt.Errorf("failed to get allowed tenant IDs: %w", err)
	}
	return tenantIDs, nil
}
func (r *tenantRepository) IsAllowedTenant(ctx context.Context, tenantID uuid.UUID) (bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT EXISTS(
			SELECT 1 FROM tenants WHERE id = $1 AND subscription_status = 'active' AND trial_ends_at > NOW() AND deleted_at IS NULL
		)
	`
	args := []interface{}{
		tenantID,
	}

	var exists bool
	err := r.db.QueryRow(subCtx, query, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check tenant allowance: %w", err)
	}

	return exists, nil
}
