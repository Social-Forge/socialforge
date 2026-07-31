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

type ChannelRepository interface {
	BaseRepository
	Create(ctx context.Context, channel *entity.Channel) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Channel, error)
	List(ctx context.Context, tenantID uuid.UUID, chType string, divisionID *uuid.UUID) ([]*entity.Channel, error)
	CountByType(ctx context.Context, tenantID uuid.UUID, chType string) (int, error)
	Update(ctx context.Context, channel *entity.Channel) (*entity.Channel, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type channelRepository struct {
	*baseRepository
}

func NewChannelRepository(db *pgxpool.Pool) ChannelRepository {
	return &channelRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *channelRepository) Create(ctx context.Context, channel *entity.Channel) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		INSERT INTO channels (
			id, tenant_id, division_id, ai_agent_id, type, name, status,
			external_id, waha_engine, waha_session_name, webhook_secret,
			credentials, settings
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id, created_at, updated_at`

	err := r.q(subCtx).QueryRow(subCtx, query,
		channel.ID, channel.TenantID, channel.DivisionID, channel.AIAgentID,
		channel.Type, channel.Name, channel.Status, channel.ExternalID,
		channel.WahaEngine, channel.WahaSessionName, channel.WebhookSecret,
		channel.Credentials, channel.Settings,
	).Scan(&channel.ID, &channel.CreatedAt, &channel.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	return nil
}

func (r *channelRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Channel, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `SELECT * FROM channels WHERE id = $1`
	var channel entity.Channel
	if err := pgxscan.Get(subCtx, r.q(subCtx), &channel, query, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("channel not found")
		}
		return nil, fmt.Errorf("failed to find channel: %w", err)
	}
	return &channel, nil
}

func (r *channelRepository) List(ctx context.Context, tenantID uuid.UUID, chType string, divisionID *uuid.UUID) ([]*entity.Channel, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM channels
		WHERE tenant_id = $1
			AND ($2 = '' OR type = $2)
			AND ($3::uuid IS NULL OR division_id = $3)
		ORDER BY created_at DESC`

	var channels []*entity.Channel
	if err := pgxscan.Select(subCtx, r.q(subCtx), &channels, query, tenantID, chType, divisionID); err != nil {
		return nil, fmt.Errorf("failed to list channels: %w", err)
	}
	return channels, nil
}

func (r *channelRepository) CountByType(ctx context.Context, tenantID uuid.UUID, chType string) (int, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `SELECT COUNT(*) FROM channels WHERE tenant_id = $1 AND type = $2`
	var count int
	if err := r.q(subCtx).QueryRow(subCtx, query, tenantID, chType).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count channels: %w", err)
	}
	return count, nil
}

func (r *channelRepository) Update(ctx context.Context, channel *entity.Channel) (*entity.Channel, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE channels SET
			division_id = $1, ai_agent_id = $2, name = $3, status = $4,
			external_id = $5, waha_engine = $6, waha_session_name = $7,
			credentials = $8, settings = $9
		WHERE id = $10
		RETURNING id, tenant_id, division_id, ai_agent_id, type, name, status,
			external_id, waha_engine, waha_session_name, webhook_secret,
			credentials, settings, created_at, updated_at`

	var updated entity.Channel
	err := r.q(subCtx).QueryRow(subCtx, query,
		channel.DivisionID, channel.AIAgentID, channel.Name, channel.Status,
		channel.ExternalID, channel.WahaEngine, channel.WahaSessionName,
		channel.Credentials, channel.Settings, channel.ID,
	).Scan(
		&updated.ID, &updated.TenantID, &updated.DivisionID, &updated.AIAgentID,
		&updated.Type, &updated.Name, &updated.Status, &updated.ExternalID,
		&updated.WahaEngine, &updated.WahaSessionName, &updated.WebhookSecret,
		&updated.Credentials, &updated.Settings, &updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("channel not found")
		}
		return nil, fmt.Errorf("failed to update channel: %w", err)
	}
	return &updated, nil
}

func (r *channelRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tag, err := r.q(subCtx).Exec(subCtx, `UPDATE channels SET status = $1 WHERE id = $2`, status, id)
	if err != nil {
		return fmt.Errorf("failed to update channel status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("channel not found")
	}
	return nil
}

func (r *channelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tag, err := r.q(subCtx).Exec(subCtx, `DELETE FROM channels WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete channel: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("channel not found")
	}
	return nil
}
