package repository

import (
	"context"
	"fmt"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AgentWorkingHoursRepository interface {
	BaseRepository
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.AgentWorkingHours, error)
	// ReplaceForUser atomically replaces all working-hour rows for a user.
	ReplaceForUser(ctx context.Context, tenantID, userID uuid.UUID, hours []*entity.AgentWorkingHours) error
}

type agentWorkingHoursRepository struct {
	*baseRepository
}

func NewAgentWorkingHoursRepository(db *pgxpool.Pool) AgentWorkingHoursRepository {
	return &agentWorkingHoursRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *agentWorkingHoursRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.AgentWorkingHours, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT id, tenant_id, user_id, day_of_week,
			start_time::text AS start_time, end_time::text AS end_time,
			is_active, created_at, updated_at
		FROM agent_working_hours WHERE user_id = $1
		ORDER BY day_of_week ASC, start_time ASC`
	var out []*entity.AgentWorkingHours
	if err := pgxscan.Select(subCtx, r.q(subCtx), &out, query, userID); err != nil {
		return nil, fmt.Errorf("failed to list working hours: %w", err)
	}
	return out, nil
}

func (r *agentWorkingHoursRepository) ReplaceForUser(ctx context.Context, tenantID, userID uuid.UUID, hours []*entity.AgentWorkingHours) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if _, err := r.q(subCtx).Exec(subCtx, `DELETE FROM agent_working_hours WHERE user_id = $1`, userID); err != nil {
		return fmt.Errorf("failed to clear working hours: %w", err)
	}
	for _, h := range hours {
		_, err := r.q(subCtx).Exec(subCtx,
			`INSERT INTO agent_working_hours (tenant_id, user_id, day_of_week, start_time, end_time, is_active)
			 VALUES ($1,$2,$3,$4::time,$5::time,$6)`,
			tenantID, userID, h.DayOfWeek, h.StartTime, h.EndTime, h.IsActive)
		if err != nil {
			return fmt.Errorf("failed to insert working hours: %w", err)
		}
	}
	return nil
}
