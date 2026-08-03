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

type DivisionRepository interface {
	BaseRepository
	Create(ctx context.Context, division *entity.Division) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Division, error)
	FindBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*entity.Division, error)
	// FindByLinkURL resolves a division by its public link token (linkchat).
	FindByLinkURL(ctx context.Context, linkURL string) (*entity.Division, error)
	GetUserDivisions(ctx context.Context, userTenantID uuid.UUID) ([]*entity.Division, error)
	GetDivisionMembers(ctx context.Context, divisionID uuid.UUID) ([]*entity.DivisionMember, error)
	AddMember(ctx context.Context, member *entity.DivisionMember) error
	RemoveMember(ctx context.Context, userTenantID, divisionID uuid.UUID) error
	Count(ctx context.Context, filter *Filter) (int64, error)
	Search(ctx context.Context, opts *ListOptions) ([]*entity.Division, int64, error)
	Update(ctx context.Context, division *entity.Division) (*entity.Division, error)
	Delete(ctx context.Context, id uuid.UUID) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	SetActiveDeactive(ctx context.Context, id uuid.UUID, isActive bool) error
	ExistsBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (bool, error)
}

type divisionRepository struct {
	*baseRepository
}

func NewDivisionRepository(db *pgxpool.Pool) DivisionRepository {
	return &divisionRepository{
		baseRepository: NewBaseRepository(db).(*baseRepository),
	}
}
func (r *divisionRepository) Create(ctx context.Context, division *entity.Division) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `INSERT INTO divisions (
			id, tenant_id, name, slug, description, routing_type, routing_config,
			is_active, link_url
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, created_at, updated_at`

	args := []interface{}{
		division.ID, division.TenantID, division.Name, division.Slug, division.Description,
		division.RoutingType, division.RoutingConfig, division.IsActive, division.LinkURL,
	}

	err := r.q(subCtx).QueryRow(subCtx, query, args...).Scan(&division.ID, &division.CreatedAt, &division.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "chk_division_routing_type":
				return fmt.Errorf("division with routing type %s already exists for tenant %s: %w", division.RoutingType, division.TenantID, err)
			default:
				return fmt.Errorf("failed to create division: %w", err)
			}
		}
		return fmt.Errorf("failed to create division: %w", err)
	}
	return nil

}
func (r *divisionRepository) Update(ctx context.Context, division *entity.Division) (*entity.Division, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE divisions SET
			name = $1, description = $2, routing_type = $3, routing_config = $4,
			is_active = $5, link_url = $6
		WHERE id = $7 AND tenant_id = $8 AND deleted_at IS NULL
		RETURNING id, tenant_id, name, slug, description, 
			routing_type, routing_config, is_active, link_url, created_at, updated_at`

	args := []interface{}{
		division.Name, division.Description, division.RoutingType, division.RoutingConfig,
		division.IsActive, division.LinkURL, division.ID, division.TenantID,
	}

	var updateDivision entity.Division
	err := r.q(subCtx).QueryRow(subCtx, query, args...).Scan(
		&updateDivision.ID,
		&updateDivision.TenantID,
		&updateDivision.Name,
		&updateDivision.Slug,
		&updateDivision.Description,
		&updateDivision.RoutingType,
		&updateDivision.RoutingConfig,
		&updateDivision.IsActive,
		&updateDivision.LinkURL,
		&updateDivision.CreatedAt,
		&updateDivision.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "chk_division_slug_tenant_id":
				return nil, fmt.Errorf("division with slug %s already exists for tenant %s: %w", division.Slug, division.TenantID, err)
			case "chk_division_routing_type":
				return nil, fmt.Errorf("division with routing type %s already exists for tenant %s: %w", division.RoutingType, division.TenantID, err)
			default:
				return nil, fmt.Errorf("failed to create division: %w", err)
			}
		}
		return nil, fmt.Errorf("failed to update division: %w", err)
	}
	return division, nil

}
func (r *divisionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Division, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `SELECT * FROM divisions WHERE id = $1 AND is_active = true AND deleted_at IS NULL`

	var division entity.Division
	err := pgxscan.Get(subCtx, r.q(subCtx), &division, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("division with id %s not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to find division by id: %w", err)
	}
	return &division, nil
}
func (r *divisionRepository) FindBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*entity.Division, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `SELECT * FROM divisions WHERE tenant_id = $1 AND slug = $2 AND is_active = true AND deleted_at IS NULL`

	var division entity.Division
	err := pgxscan.Get(subCtx, r.q(subCtx), &division, query, tenantID, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("division with slug %s for tenant %s not found: %w", slug, tenantID, err)
		}
		return nil, fmt.Errorf("failed to find division by slug: %w", err)
	}
	return &division, nil
}
func (r *divisionRepository) FindByLinkURL(ctx context.Context, linkURL string) (*entity.Division, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	// Public lookup by unique link token (works as superuser in dev; prod needs
	// a SECURITY DEFINER function since divisions is RLS-forced).
	query := `SELECT * FROM divisions WHERE link_url = $1 AND is_active = true AND deleted_at IS NULL`
	var division entity.Division
	if err := pgxscan.Get(subCtx, r.q(subCtx), &division, query, linkURL); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("link not found")
		}
		return nil, fmt.Errorf("failed to resolve link: %w", err)
	}
	return &division, nil
}

func (r *divisionRepository) GetUserDivisions(ctx context.Context, userTenantID uuid.UUID) ([]*entity.Division, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
        SELECT
            d.id, d.tenant_id, d.name, d.slug, d.description,
            d.routing_type, d.routing_config, d.is_active,
            d.link_url, d.created_at, d.updated_at, d.deleted_at
        FROM divisions d
        JOIN division_members dm ON d.id = dm.division_id
        WHERE dm.user_tenant_id = $1
          AND dm.deleted_at IS NULL
          AND dm.is_active = true
          AND d.deleted_at IS NULL
          AND d.is_active = true
        ORDER BY d.name ASC
    `

	var divisions []*entity.Division
	err := pgxscan.Select(subCtx, r.q(subCtx), &divisions, query, userTenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.Division{}, nil
		}
		return nil, err
	}

	return divisions, nil
}
func (r *divisionRepository) GetDivisionMembers(ctx context.Context, divisionID uuid.UUID) ([]*entity.DivisionMember, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
        SELECT dm.id, dm.user_tenant_id, dm.division_id, dm.is_active,
               dm.joined_at, dm.created_at, dm.updated_at, dm.deleted_at
        FROM division_members dm
        JOIN user_tenants ut ON dm.user_tenant_id = ut.id
        JOIN users u ON ut.user_id = u.id
        WHERE dm.division_id = $1
          AND dm.deleted_at IS NULL
          AND dm.is_active = true
        ORDER BY u.full_name ASC
    `

	var members []*entity.DivisionMember
	err := pgxscan.Select(subCtx, r.q(subCtx), &members, query, divisionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.DivisionMember{}, nil
		}
		return nil, fmt.Errorf("failed to get division members: %w", err)
	}
	return members, nil
}
func (r *divisionRepository) AddMember(ctx context.Context, member *entity.DivisionMember) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
        INSERT INTO division_members (id, user_tenant_id, division_id, is_active, joined_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (user_tenant_id, division_id) DO UPDATE SET
            is_active = EXCLUDED.is_active,
            deleted_at = NULL,
            updated_at = NOW()
        RETURNING id, created_at, updated_at
    `

	err := r.q(subCtx).QueryRow(subCtx, query,
		member.ID,
		member.UserTenantID,
		member.DivisionID,
		member.IsActive,
		member.JoinedAt,
	).Scan(&member.ID, &member.CreatedAt, &member.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to add member to division: %w", err)
	}
	return nil
}
func (r *divisionRepository) RemoveMember(ctx context.Context, userTenantID, divisionID uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
        UPDATE division_members SET
            deleted_at = NOW(),
            is_active = false,
            updated_at = NOW()
        WHERE user_tenant_id = $1 AND division_id = $2 AND deleted_at IS NULL
    `

	result, err := r.q(subCtx).Exec(subCtx, query, userTenantID, divisionID)
	if err != nil {
		return fmt.Errorf("failed to remove member from division: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("member not found in division")
	}
	return nil
}
func (r *divisionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE divisions SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	args := []interface{}{
		id,
	}

	cmdTag, err := r.q(subCtx).Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete division: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("division with id %s not found or already deleted", id)
	}
	return nil
}
func (r *divisionRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `DELETE FROM divisions WHERE id = $1 AND deleted_at IS NULL`

	args := []interface{}{
		id,
	}

	cmdTag, err := r.q(subCtx).Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to hard delete division: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("division with id %s not found or already deleted", id)
	}
	return nil
}
func (r *divisionRepository) Restore(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE divisions SET deleted_at = NULL WHERE id = $1 AND deleted_at IS NOT NULL`

	args := []interface{}{
		id,
	}

	cmdTag, err := r.q(subCtx).Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to restore division: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("division with id %s not found or already restored", id)
	}
	return nil
}
func (r *divisionRepository) SetActiveDeactive(ctx context.Context, id uuid.UUID, isActive bool) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE divisions SET is_active = $1 WHERE id = $2 AND deleted_at IS NULL`

	args := []interface{}{
		isActive,
		id,
	}

	cmdTag, err := r.q(subCtx).Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to set active/deactive division: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("division with id %s not found or already %v", id, isActive)
	}
	return nil
}
func (r *divisionRepository) Count(ctx context.Context, filter *Filter) (int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	qb := r.buildBaseQuery("SELECT COUNT(*) FROM divisions", filter)
	query, args := qb.Build()

	var count int64
	err := r.q(subCtx).QueryRow(subCtx, query, args...).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("no divisions found: %w", err)
		}
		return 0, fmt.Errorf("failed to count division: %w", err)
	}
	return count, nil
}
func (r *divisionRepository) Search(ctx context.Context, opts *ListOptions) ([]*entity.Division, int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if opts == nil {
		opts = &ListOptions{}
	}

	totalRows, err := r.Count(ctx, opts.Filter)
	if err != nil {
		return nil, 0, err
	}

	qb := r.buildBaseQuery("SELECT * FROM divisions", opts.Filter)

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

	var divisions []*entity.Division
	err = pgxscan.Select(subCtx, r.q(subCtx), &divisions, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, nil
		}
		return nil, 0, fmt.Errorf("failed to get list division: %w", err)
	}
	return divisions, totalRows, nil
}
func (r *divisionRepository) ExistsBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `SELECT EXISTS(SELECT 1 FROM divisions WHERE tenant_id = $1 AND slug = $2 AND deleted_at IS NULL)`

	var exists bool
	err := r.q(subCtx).QueryRow(subCtx, query, tenantID, slug).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check division existence by slug: %w", err)
	}
	return exists, nil
}
func (r *divisionRepository) buildBaseQuery(baseQuery string, filter *Filter) *QueryBuilder {
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
		qb.Where(
			"(name ILIKE $? OR slug ILIKE $? OR description ILIKE $? OR routing_type ILIKE $?)",
			searchPattern, searchPattern, searchPattern, searchPattern,
		)
	}
	if filter.TenantID != nil {
		qb.Where("tenant_id = $?", *filter.TenantID)
	}
	if filter.IsActive != nil {
		qb.Where("is_active = $?", *filter.IsActive)
	}
	// Scope divisions to those the given user is a member of (via user_tenants).
	if filter.UserID != nil {
		qb.Where(`EXISTS (
        SELECT 1 FROM division_members dm
        JOIN user_tenants ut ON dm.user_tenant_id = ut.id
        WHERE dm.division_id = divisions.id
          AND ut.user_id = $?
          AND dm.deleted_at IS NULL
    )`, *filter.UserID)
	}
	if filter.Extra != nil {
		if routingType, ok := filter.Extra["routing_type"].(string); ok {
			qb.Where("routing_type = $?", routingType)
		}
	}

	return qb
}
