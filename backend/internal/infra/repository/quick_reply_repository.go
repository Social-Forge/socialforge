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

var ErrShortcutTaken = errors.New("shortcut already exists in this tenant")

type QuickReplyRepository interface {
	BaseRepository
	Create(ctx context.Context, qr *entity.QuickReply) error
	List(ctx context.Context, tenantID uuid.UUID, search string) ([]*entity.QuickReply, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.QuickReply, error)
	Update(ctx context.Context, qr *entity.QuickReply) (*entity.QuickReply, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context, tenantID uuid.UUID) (int, error)
}

type quickReplyRepository struct {
	*baseRepository
}

func NewQuickReplyRepository(db *pgxpool.Pool) QuickReplyRepository {
	return &quickReplyRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *quickReplyRepository) Create(ctx context.Context, qr *entity.QuickReply) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if qr.ID == uuid.Nil {
		qr.ID = uuid.New()
	}
	query := `
		INSERT INTO quick_replies (id, tenant_id, shortcut, content_type, body, media)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, created_at, updated_at`
	err := r.q(subCtx).QueryRow(subCtx, query, qr.ID, qr.TenantID, qr.Shortcut, qr.ContentType, qr.Body, qr.Media).
		Scan(&qr.ID, &qr.CreatedAt, &qr.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrShortcutTaken
		}
		return fmt.Errorf("failed to create quick reply: %w", err)
	}
	return nil
}

func (r *quickReplyRepository) List(ctx context.Context, tenantID uuid.UUID, search string) ([]*entity.QuickReply, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM quick_replies
		WHERE tenant_id = $1 AND ($2 = '' OR shortcut ILIKE '%'||$2||'%' OR body ILIKE '%'||$2||'%')
		ORDER BY shortcut ASC
		LIMIT 200`
	var out []*entity.QuickReply
	if err := pgxscan.Select(subCtx, r.q(subCtx), &out, query, tenantID, search); err != nil {
		return nil, fmt.Errorf("failed to list quick replies: %w", err)
	}
	return out, nil
}

func (r *quickReplyRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.QuickReply, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var qr entity.QuickReply
	if err := pgxscan.Get(subCtx, r.q(subCtx), &qr, `SELECT * FROM quick_replies WHERE id = $1`, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("quick reply not found")
		}
		return nil, fmt.Errorf("failed to find quick reply: %w", err)
	}
	return &qr, nil
}

func (r *quickReplyRepository) Update(ctx context.Context, qr *entity.QuickReply) (*entity.QuickReply, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE quick_replies SET shortcut = $1, content_type = $2, body = $3, media = $4
		WHERE id = $5
		RETURNING id, tenant_id, shortcut, content_type, body, media, created_at, updated_at`
	var out entity.QuickReply
	err := r.q(subCtx).QueryRow(subCtx, query, qr.Shortcut, qr.ContentType, qr.Body, qr.Media, qr.ID).
		Scan(&out.ID, &out.TenantID, &out.Shortcut, &out.ContentType, &out.Body, &out.Media, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrShortcutTaken
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("quick reply not found")
		}
		return nil, fmt.Errorf("failed to update quick reply: %w", err)
	}
	return &out, nil
}

func (r *quickReplyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tag, err := r.q(subCtx).Exec(subCtx, `DELETE FROM quick_replies WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete quick reply: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("quick reply not found")
	}
	return nil
}

func (r *quickReplyRepository) Count(ctx context.Context, tenantID uuid.UUID) (int, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var count int
	if err := r.q(subCtx).QueryRow(subCtx, `SELECT COUNT(*) FROM quick_replies WHERE tenant_id = $1`, tenantID).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count quick replies: %w", err)
	}
	return count, nil
}
