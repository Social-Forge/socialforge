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

type AutoResponseRepository interface {
	BaseRepository
	// GetByChannel returns the channel's auto-response config, or nil if unset.
	GetByChannel(ctx context.Context, channelID uuid.UUID) (*entity.AutoResponse, error)
	// Upsert creates or updates the single auto-response for a channel.
	Upsert(ctx context.Context, ar *entity.AutoResponse) (*entity.AutoResponse, error)
}

type autoResponseRepository struct {
	*baseRepository
}

func NewAutoResponseRepository(db *pgxpool.Pool) AutoResponseRepository {
	return &autoResponseRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *autoResponseRepository) GetByChannel(ctx context.Context, channelID uuid.UUID) (*entity.AutoResponse, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var ar entity.AutoResponse
	err := pgxscan.Get(subCtx, r.q(subCtx), &ar, `SELECT * FROM auto_responses WHERE channel_id = $1`, channelID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get auto response: %w", err)
	}
	return &ar, nil
}

func (r *autoResponseRepository) Upsert(ctx context.Context, ar *entity.AutoResponse) (*entity.AutoResponse, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if ar.ID == uuid.Nil {
		ar.ID = uuid.New()
	}
	query := `
		INSERT INTO auto_responses (id, tenant_id, channel_id, is_enabled, content_type, body, media)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (channel_id) DO UPDATE SET
			is_enabled = EXCLUDED.is_enabled,
			content_type = EXCLUDED.content_type,
			body = EXCLUDED.body,
			media = EXCLUDED.media,
			updated_at = now()
		RETURNING id, tenant_id, channel_id, is_enabled, content_type, body, media, created_at, updated_at`
	var out entity.AutoResponse
	err := r.q(subCtx).QueryRow(subCtx, query, ar.ID, ar.TenantID, ar.ChannelID, ar.IsEnabled, ar.ContentType, ar.Body, ar.Media).
		Scan(&out.ID, &out.TenantID, &out.ChannelID, &out.IsEnabled, &out.ContentType, &out.Body, &out.Media, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert auto response: %w", err)
	}
	return &out, nil
}
