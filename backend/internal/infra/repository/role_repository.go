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

type RoleRepository interface {
	BaseRepository
	Create(ctx context.Context, role *entity.Role) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Role, error)
	GetByName(ctx context.Context, name string) (*entity.Role, error)
	Count(ctx context.Context, filter *Filter) (int64, error)
	Search(ctx context.Context, opts *ListOptions) ([]*entity.Role, int64, error)
	GetAll(ctx context.Context) ([]*entity.Role, error)
	Update(ctx context.Context, role *entity.Role) (*entity.Role, error)
	Delete(ctx context.Context, id uuid.UUID) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
}
type roleRepository struct {
	*baseRepository
}

func NewRoleRepository(db *pgxpool.Pool) RoleRepository {
	return &roleRepository{
		baseRepository: NewBaseRepository(db).(*baseRepository),
	}
}
func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `INSERT INTO roles (id, name, slug, description, level, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (name) DO NOTHING
	RETURNING id, created_at, updated_at`

	args := []interface{}{
		role.ID,
		role.Name,
		role.Slug,
		role.Description,
		role.Level,
		role.CreatedAt,
	}

	var roleID uuid.UUID
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(subCtx, query, args...).Scan(&roleID, &createdAt, &updatedAt)
	if err != nil {
		config.Logger.Error("Failed to create role", zap.Error(err))
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			switch pgxErr.ConstraintName {
			case "roles_slug_key":
				return fmt.Errorf("role slug already exists")
			case "roles_name_length_check":
				return fmt.Errorf("role name length must be between 3 and 20 characters")
			case "roles_slug_length_check":
				return fmt.Errorf("role slug length must be between 3 and 20 characters")
			case "roles_level_check":
				return fmt.Errorf("role level must be between 0 and 100")
			case "roles_name_check":
				return fmt.Errorf("role name must contain only letters, numbers, and underscores")
			default:
				return fmt.Errorf("unique constraint violation (%s): %w", pgxErr.ConstraintName, err)
			}
		}
		return fmt.Errorf("failed to create role: %w", err)
	}

	role.ID = roleID
	role.CreatedAt = createdAt
	role.UpdatedAt = updatedAt
	return nil
}
func (r *roleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `SELECT * FROM roles WHERE id = $1 AND deleted_at IS NULL`

	args := []interface{}{
		id,
	}

	var role entity.Role
	err := pgxscan.Get(subCtx, r.db, &role, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find role by id: %w", err)
	}
	return &role, nil
}
func (r *roleRepository) FindBySlug(ctx context.Context, slug string) (*entity.Role, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `SELECT * FROM roles WHERE slug = $1 AND deleted_at IS NULL`

	args := []interface{}{
		slug,
	}

	var role entity.Role
	err := pgxscan.Get(subCtx, r.db, &role, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find role by slug: %w", err)
	}
	return &role, nil
}
func (r *roleRepository) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `SELECT * FROM roles WHERE name = $1 AND deleted_at IS NULL`

	args := []interface{}{
		name,
	}

	var role entity.Role
	err := pgxscan.Get(subCtx, r.db, &role, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find role by name: %w", err)
	}
	return &role, nil
}

func (r *roleRepository) Count(ctx context.Context, filter *Filter) (int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	qb := r.buildBaseQuery("SELECT COUNT(*) FROM roles", filter)
	query, args := qb.Build()

	var count int64
	err := r.db.QueryRow(subCtx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count roles: %w", err)
	}
	return count, nil
}
func (r *roleRepository) GetAll(ctx context.Context) ([]*entity.Role, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `SELECT * FROM roles WHERE deleted_at IS NULL`

	var roles []*entity.Role
	err := pgxscan.Select(subCtx, r.db, &roles, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	return roles, nil
}
func (r *roleRepository) Search(ctx context.Context, opts *ListOptions) ([]*entity.Role, int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if opts == nil {
		opts = NewListOptions()
	}

	totalRows, err := r.Count(subCtx, opts.Filter)
	if err != nil {
		return nil, 0, err
	}

	qb := r.buildBaseQuery("SELECT * FROM roles", opts.Filter)

	if opts.OrderBy != "" {
		qb.OrderByField(opts.OrderBy, opts.OrderDir)
	} else {
		qb.OrderByField("level", "ASC") // Default ordering
	}
	if opts.Pagination != nil && opts.Pagination.Limit > 0 {
		qb.WithLimit(opts.Pagination.Limit)
		if opts.Pagination.Page > 1 {
			qb.WithOffset((opts.Pagination.Page - 1) * opts.Pagination.Limit)
		}
	}

	query, args := qb.Build()

	var roles []*entity.Role
	err = pgxscan.Select(subCtx, r.db, &roles, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, totalRows, nil
		}
		return nil, totalRows, fmt.Errorf("failed to list roles: %w", err)
	}
	return roles, totalRows, nil
}
func (r *roleRepository) Update(ctx context.Context, role *entity.Role) (*entity.Role, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE roles SET 
	name = $1, 
	slug = $2, 
	description = $3, 
	level = $4, 
	WHERE id = $5 AND deleted_at IS NULL
	RETURNING id, name, slug, description, level, created_at, updated_at`

	args := []interface{}{
		role.Name,
		role.Slug,
		role.Description,
		role.Level,
		role.ID,
	}

	var updateRole entity.Role
	err := r.db.QueryRow(subCtx, query, args...).Scan(
		&updateRole.ID,
		&updateRole.Name,
		&updateRole.Slug,
		&updateRole.Description,
		&updateRole.Level,
		&updateRole.CreatedAt,
		&updateRole.UpdatedAt,
	)

	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			switch pgxErr.ConstraintName {
			case "roles_name_key":
				return nil, fmt.Errorf("role name '%s' is already taken", role.Name)
			case "roles_slug_key":
				return nil, fmt.Errorf("role slug '%s' is already taken", role.Slug)
			case "roles_name_length_check":
				return nil, fmt.Errorf("role name length must be between 3 and 20 characters")
			case "roles_slug_length_check":
				return nil, fmt.Errorf("role slug length must be between 3 and 20 characters")
			case "roles_level_check":
				return nil, fmt.Errorf("role level must be between 0 and 100")
			case "roles_name_check":
				return nil, fmt.Errorf("role name must contain only letters, numbers, and underscores")
			default:
				return nil, fmt.Errorf("unique constraint violation (%s): %w", pgxErr.ConstraintName, err)
			}
		}
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return &updateRole, nil
}
func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE roles SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	args := []interface{}{
		id,
	}

	cmdTag, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("role not found or already deleted")
	}

	return nil
}
func (r *roleRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `DELETE FROM roles WHERE id = $1`

	args := []interface{}{
		id,
	}

	cmdTag, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to hard delete role: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("role not found or already deleted")
	}

	return nil
}
func (r *roleRepository) Restore(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE roles SET deleted_at = NULL WHERE id = $1 AND deleted_at IS NOT NULL`

	args := []interface{}{
		id,
	}

	cmdTag, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to restore role: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("role not found or not deleted")
	}

	return nil
}
func (r *roleRepository) buildBaseQuery(baseQuery string, filter *Filter) *QueryBuilder {
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
		qb.Where("(name ILIKE $? OR slug ILIKE $? OR description ILIKE $?)", searchPattern, searchPattern, searchPattern)
	}

	return qb
}
