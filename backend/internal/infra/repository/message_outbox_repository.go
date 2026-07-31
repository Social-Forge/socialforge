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

type MessageOutboxRepository interface {
	BaseRepository
	Create(ctx context.Context, o *entity.MessageOutbox) error
	MarkSent(ctx context.Context, id uuid.UUID) error
	MarkFailed(ctx context.Context, id uuid.UUID, lastError string) error
	SetStatusByMessage(ctx context.Context, messageID uuid.UUID, status, lastError string) error
}

type messageOutboxRepository struct {
	*baseRepository
}

func NewMessageOutboxRepository(db *pgxpool.Pool) MessageOutboxRepository {
	return &messageOutboxRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *messageOutboxRepository) Create(ctx context.Context, o *entity.MessageOutbox) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	if o.Status == "" {
		o.Status = entity.OutboxStatusPending
	}
	if o.MaxAttempts == 0 {
		o.MaxAttempts = 3
	}
	query := `
		INSERT INTO message_outboxes (id, tenant_id, message_id, status, attempts, max_attempts)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, created_at, updated_at`
	return r.q(subCtx).QueryRow(subCtx, query, o.ID, o.TenantID, o.MessageID, o.Status, o.Attempts, o.MaxAttempts).
		Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
}

func (r *messageOutboxRepository) MarkSent(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	_, err := r.q(subCtx).Exec(subCtx,
		`UPDATE message_outboxes SET status = 'sent', attempts = attempts + 1, last_error = NULL WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to mark outbox sent: %w", err)
	}
	return nil
}

func (r *messageOutboxRepository) SetStatusByMessage(ctx context.Context, messageID uuid.UUID, status, lastError string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	_, err := r.q(subCtx).Exec(subCtx,
		`UPDATE message_outboxes SET status = $2, attempts = attempts + 1, last_error = NULLIF($3,'') WHERE message_id = $1`,
		messageID, status, lastError)
	if err != nil {
		return fmt.Errorf("failed to update outbox by message: %w", err)
	}
	return nil
}

func (r *messageOutboxRepository) MarkFailed(ctx context.Context, id uuid.UUID, lastError string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	_, err := r.q(subCtx).Exec(subCtx,
		`UPDATE message_outboxes SET status = 'failed', attempts = attempts + 1, last_error = $2 WHERE id = $1`, id, lastError)
	if err != nil {
		return fmt.Errorf("failed to mark outbox failed: %w", err)
	}
	return nil
}
