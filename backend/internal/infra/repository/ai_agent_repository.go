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

type AIAgentRepository interface {
	BaseRepository
	Create(ctx context.Context, agent *entity.AIAgent) error
	List(ctx context.Context, tenantID uuid.UUID) ([]*entity.AIAgent, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.AIAgent, error)
	Update(ctx context.Context, agent *entity.AIAgent) (*entity.AIAgent, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context, tenantID uuid.UUID) (int, error)
}

type aiAgentRepository struct {
	*baseRepository
}

func NewAIAgentRepository(db *pgxpool.Pool) AIAgentRepository {
	return &aiAgentRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *aiAgentRepository) Create(ctx context.Context, a *entity.AIAgent) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	query := `
		INSERT INTO ai_agents (id, tenant_id, name, provider, model, system_prompt,
			persona, safety, guardrails, temperature, max_tokens, auto_reply_enabled, working_hours, is_active)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		RETURNING id, created_at, updated_at`
	err := r.q(subCtx).QueryRow(subCtx, query,
		a.ID, a.TenantID, a.Name, a.Provider, a.Model, a.SystemPrompt,
		a.Persona, a.Safety, a.Guardrails, a.Temperature, a.MaxTokens, a.AutoReplyEnabled, a.WorkingHours, a.IsActive,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create ai agent: %w", err)
	}
	return nil
}

func (r *aiAgentRepository) List(ctx context.Context, tenantID uuid.UUID) ([]*entity.AIAgent, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var out []*entity.AIAgent
	if err := pgxscan.Select(subCtx, r.q(subCtx), &out,
		`SELECT * FROM ai_agents WHERE tenant_id = $1 ORDER BY created_at DESC`, tenantID); err != nil {
		return nil, fmt.Errorf("failed to list ai agents: %w", err)
	}
	return out, nil
}

func (r *aiAgentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.AIAgent, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var a entity.AIAgent
	if err := pgxscan.Get(subCtx, r.q(subCtx), &a, `SELECT * FROM ai_agents WHERE id = $1`, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ai agent not found")
		}
		return nil, fmt.Errorf("failed to find ai agent: %w", err)
	}
	return &a, nil
}

func (r *aiAgentRepository) Update(ctx context.Context, a *entity.AIAgent) (*entity.AIAgent, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE ai_agents SET name=$1, provider=$2, model=$3, system_prompt=$4,
			persona=$5, safety=$6, guardrails=$7, temperature=$8, max_tokens=$9,
			auto_reply_enabled=$10, working_hours=$11, is_active=$12
		WHERE id=$13
		RETURNING id, tenant_id, name, provider, model, system_prompt, persona, safety,
			guardrails, temperature, max_tokens, auto_reply_enabled, working_hours, is_active, created_at, updated_at`
	var out entity.AIAgent
	err := r.q(subCtx).QueryRow(subCtx, query,
		a.Name, a.Provider, a.Model, a.SystemPrompt, a.Persona, a.Safety, a.Guardrails,
		a.Temperature, a.MaxTokens, a.AutoReplyEnabled, a.WorkingHours, a.IsActive, a.ID,
	).Scan(&out.ID, &out.TenantID, &out.Name, &out.Provider, &out.Model, &out.SystemPrompt,
		&out.Persona, &out.Safety, &out.Guardrails, &out.Temperature, &out.MaxTokens,
		&out.AutoReplyEnabled, &out.WorkingHours, &out.IsActive, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ai agent not found")
		}
		return nil, fmt.Errorf("failed to update ai agent: %w", err)
	}
	return &out, nil
}

func (r *aiAgentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tag, err := r.q(subCtx).Exec(subCtx, `DELETE FROM ai_agents WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete ai agent: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("ai agent not found")
	}
	return nil
}

func (r *aiAgentRepository) Count(ctx context.Context, tenantID uuid.UUID) (int, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var count int
	if err := r.q(subCtx).QueryRow(subCtx, `SELECT COUNT(*) FROM ai_agents WHERE tenant_id = $1`, tenantID).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count ai agents: %w", err)
	}
	return count, nil
}
