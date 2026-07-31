package repository

import (
	"context"
	"fmt"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WebhookEventRepository interface {
	BaseRepository
	// TryInsert records an inbound event idempotently. inserted=false means the
	// event (channel_id, provider_event_id) was already seen -> skip processing.
	TryInsert(ctx context.Context, e *entity.WebhookEvent) (inserted bool, err error)
	MarkProcessed(ctx context.Context, id uuid.UUID) error
}

type webhookEventRepository struct {
	*baseRepository
}

func NewWebhookEventRepository(db *pgxpool.Pool) WebhookEventRepository {
	return &webhookEventRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *webhookEventRepository) TryInsert(ctx context.Context, e *entity.WebhookEvent) (bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if len(e.Payload) == 0 {
		e.Payload = []byte("{}")
	}

	// webhook_events is not tenant-RLS-scoped (keyed by channel), so it can be
	// checked before the tenant context is resolved from the channel.
	query := `
		INSERT INTO webhook_events (id, channel_id, provider, provider_event_id, event_type, payload)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (channel_id, provider_event_id) DO NOTHING`

	tag, err := r.q(subCtx).Exec(subCtx, query, e.ID, e.ChannelID, e.Provider, e.ProviderEventID, e.EventType, e.Payload)
	if err != nil {
		return false, fmt.Errorf("failed to record webhook event: %w", err)
	}
	return tag.RowsAffected() > 0, nil
}

func (r *webhookEventRepository) MarkProcessed(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	_, err := r.q(subCtx).Exec(subCtx, `UPDATE webhook_events SET processed_at = now() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to mark webhook processed: %w", err)
	}
	return nil
}
