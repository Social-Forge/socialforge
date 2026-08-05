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
	"github.com/jackc/pgx/v5/pgxpool"
)

// PlanRepository manages the global plan catalog (no tenant scoping / RLS).
type PlanRepository interface {
	BaseRepository
	List(ctx context.Context, activeOnly bool) ([]*entity.Plan, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Plan, error)
	FindByCode(ctx context.Context, code string) (*entity.Plan, error)
	Create(ctx context.Context, p *entity.Plan) error
	Update(ctx context.Context, p *entity.Plan) (*entity.Plan, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type planRepository struct {
	*baseRepository
}

func NewPlanRepository(db *pgxpool.Pool) PlanRepository {
	return &planRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *planRepository) List(ctx context.Context, activeOnly bool) ([]*entity.Plan, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	query := `SELECT * FROM plans`
	if activeOnly {
		query += ` WHERE is_active = TRUE`
	}
	query += ` ORDER BY sort ASC, price ASC`
	var out []*entity.Plan
	if err := pgxscan.Select(subCtx, r.q(subCtx), &out, query); err != nil {
		return nil, fmt.Errorf("failed to list plans: %w", err)
	}
	return out, nil
}

func (r *planRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Plan, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var p entity.Plan
	if err := pgxscan.Get(subCtx, r.q(subCtx), &p, `SELECT * FROM plans WHERE id = $1`, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("plan not found")
		}
		return nil, fmt.Errorf("failed to find plan: %w", err)
	}
	return &p, nil
}

func (r *planRepository) FindByCode(ctx context.Context, code string) (*entity.Plan, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var p entity.Plan
	if err := pgxscan.Get(subCtx, r.q(subCtx), &p, `SELECT * FROM plans WHERE code = $1`, code); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("plan not found")
		}
		return nil, fmt.Errorf("failed to find plan: %w", err)
	}
	return &p, nil
}

func (r *planRepository) Create(ctx context.Context, p *entity.Plan) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	query := `
		INSERT INTO plans (id, code, name, price, currency, interval, features, is_active, sort)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at, updated_at`
	if err := r.q(subCtx).QueryRow(subCtx, query,
		p.ID, p.Code, p.Name, p.Price, p.Currency, p.Interval, p.Features, p.IsActive, p.Sort,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return fmt.Errorf("failed to create plan: %w", err)
	}
	return nil
}

func (r *planRepository) Update(ctx context.Context, p *entity.Plan) (*entity.Plan, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var out entity.Plan
	query := `
		UPDATE plans SET code=$1, name=$2, price=$3, currency=$4, interval=$5, features=$6, is_active=$7, sort=$8
		WHERE id=$9
		RETURNING *`
	if err := pgxscan.Get(subCtx, r.q(subCtx), &out, query,
		p.Code, p.Name, p.Price, p.Currency, p.Interval, p.Features, p.IsActive, p.Sort, p.ID,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("plan not found")
		}
		return nil, fmt.Errorf("failed to update plan: %w", err)
	}
	return &out, nil
}

func (r *planRepository) Delete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	tag, err := r.q(subCtx).Exec(subCtx, `DELETE FROM plans WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete plan: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("plan not found")
	}
	return nil
}
