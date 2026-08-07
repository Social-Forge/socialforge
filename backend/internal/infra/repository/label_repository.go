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

type LabelRepository interface {
	BaseRepository
	Create(ctx context.Context, label *entity.Label) error
	List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entity.Label, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Label, error)
	Update(ctx context.Context, label *entity.Label) (*entity.Label, error)
	Delete(ctx context.Context, id uuid.UUID) error
	// Conversation <-> label attachment.
	Attach(ctx context.Context, cl *entity.ConversationLabel) error
	Detach(ctx context.Context, conversationID, labelID uuid.UUID) error
	ListForConversation(ctx context.Context, conversationID uuid.UUID) ([]*entity.Label, error)
}

type labelRepository struct {
	*baseRepository
}

func NewLabelRepository(db *pgxpool.Pool) LabelRepository {
	return &labelRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

var ErrLabelNameTaken = errors.New("label name already exists in this tenant")

func (r *labelRepository) Create(ctx context.Context, label *entity.Label) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if label.ID == uuid.Nil {
		label.ID = uuid.New()
	}
	if label.Color == "" {
		label.Color = "#64748b"
	}
	query := `
		INSERT INTO labels (id, tenant_id, name, color)
		VALUES ($1,$2,$3,$4)
		RETURNING id, created_at, updated_at`
	err := r.q(subCtx).QueryRow(subCtx, query, label.ID, label.TenantID, label.Name, label.Color).
		Scan(&label.ID, &label.CreatedAt, &label.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrLabelNameTaken
		}
		return fmt.Errorf("failed to create label: %w", err)
	}
	return nil
}

func (r *labelRepository) List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entity.Label, int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var total int64
	if err := r.q(subCtx).QueryRow(subCtx, `SELECT COUNT(*) FROM labels WHERE tenant_id = $1`, tenantID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count labels: %w", err)
	}

	var labels []*entity.Label
	if err := pgxscan.Select(subCtx, r.q(subCtx), &labels,
		`SELECT * FROM labels WHERE tenant_id = $1 ORDER BY name ASC LIMIT $2 OFFSET $3`, tenantID, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("failed to list labels: %w", err)
	}
	return labels, total, nil
}

func (r *labelRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Label, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var label entity.Label
	if err := pgxscan.Get(subCtx, r.q(subCtx), &label, `SELECT * FROM labels WHERE id = $1`, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("label not found")
		}
		return nil, fmt.Errorf("failed to find label: %w", err)
	}
	return &label, nil
}

func (r *labelRepository) Update(ctx context.Context, label *entity.Label) (*entity.Label, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `UPDATE labels SET name = $1, color = $2 WHERE id = $3
		RETURNING id, tenant_id, name, color, created_at, updated_at`
	var out entity.Label
	err := r.q(subCtx).QueryRow(subCtx, query, label.Name, label.Color, label.ID).
		Scan(&out.ID, &out.TenantID, &out.Name, &out.Color, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrLabelNameTaken
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("label not found")
		}
		return nil, fmt.Errorf("failed to update label: %w", err)
	}
	return &out, nil
}

func (r *labelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tag, err := r.q(subCtx).Exec(subCtx, `DELETE FROM labels WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete label: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("label not found")
	}
	return nil
}

func (r *labelRepository) Attach(ctx context.Context, cl *entity.ConversationLabel) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if cl.ID == uuid.Nil {
		cl.ID = uuid.New()
	}
	query := `
		INSERT INTO conversation_labels (id, tenant_id, conversation_id, label_id)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (conversation_id, label_id) DO NOTHING`
	if _, err := r.q(subCtx).Exec(subCtx, query, cl.ID, cl.TenantID, cl.ConversationID, cl.LabelID); err != nil {
		return fmt.Errorf("failed to attach label: %w", err)
	}
	return nil
}

func (r *labelRepository) Detach(ctx context.Context, conversationID, labelID uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	_, err := r.q(subCtx).Exec(subCtx,
		`DELETE FROM conversation_labels WHERE conversation_id = $1 AND label_id = $2`, conversationID, labelID)
	if err != nil {
		return fmt.Errorf("failed to detach label: %w", err)
	}
	return nil
}

func (r *labelRepository) ListForConversation(ctx context.Context, conversationID uuid.UUID) ([]*entity.Label, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT l.* FROM labels l
		JOIN conversation_labels cl ON cl.label_id = l.id
		WHERE cl.conversation_id = $1
		ORDER BY l.name ASC`
	var labels []*entity.Label
	if err := pgxscan.Select(subCtx, r.q(subCtx), &labels, query, conversationID); err != nil {
		return nil, fmt.Errorf("failed to list conversation labels: %w", err)
	}
	return labels, nil
}
