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

type MessageRepository interface {
	BaseRepository
	// Create inserts a message. inserted=false means it was a duplicate
	// (same provider_message_id) and was skipped — idempotent ingestion.
	Create(ctx context.Context, message *entity.Message) (inserted bool, err error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Message, error)
	ListByConversation(ctx context.Context, conversationID uuid.UUID, limit int) ([]*entity.Message, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, errMsg string) error
}

type messageRepository struct {
	*baseRepository
}

func NewMessageRepository(db *pgxpool.Pool) MessageRepository {
	return &messageRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *messageRepository) Create(ctx context.Context, m *entity.Message) (bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if m.Status == "" {
		m.Status = entity.MessageStatusPending
	}

	query := `
		INSERT INTO messages (
			id, tenant_id, conversation_id, sender_id, direction, sender_type,
			content_type, body, media, provider_message_id, status, reply_to_id
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		ON CONFLICT (provider_message_id) DO NOTHING
		RETURNING id, created_at, updated_at`

	err := r.q(subCtx).QueryRow(subCtx, query,
		m.ID, m.TenantID, m.ConversationID, m.SenderID, m.Direction, m.SenderType,
		m.ContentType, m.Body, m.Media, m.ProviderMessageID, m.Status, m.ReplyToID,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Conflict on provider_message_id -> duplicate, skipped.
			return false, nil
		}
		return false, fmt.Errorf("failed to create message: %w", err)
	}
	return true, nil
}

func (r *messageRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Message, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var msg entity.Message
	if err := pgxscan.Get(subCtx, r.q(subCtx), &msg, `SELECT * FROM messages WHERE id = $1 AND deleted_at IS NULL`, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to find message: %w", err)
	}
	return &msg, nil
}

func (r *messageRepository) ListByConversation(ctx context.Context, conversationID uuid.UUID, limit int) ([]*entity.Message, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if limit <= 0 || limit > 200 {
		limit = 50
	}
	query := `
		SELECT * FROM messages
		WHERE conversation_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT $2`

	var messages []*entity.Message
	if err := pgxscan.Select(subCtx, r.q(subCtx), &messages, query, conversationID, limit); err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	return messages, nil
}

func (r *messageRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status, errMsg string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	_, err := r.q(subCtx).Exec(subCtx,
		`UPDATE messages SET status = $1, error = NULLIF($2, '') WHERE id = $3`, status, errMsg, id)
	if err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}
	return nil
}
