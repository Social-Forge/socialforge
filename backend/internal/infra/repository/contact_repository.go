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

type ContactRepository interface {
	BaseRepository
	// FindOrCreate upserts a contact by (channel_id, external_id) and returns it,
	// plus created=true when the contact was newly inserted (first-time customer).
	FindOrCreate(ctx context.Context, contact *entity.Contact) (result *entity.Contact, created bool, err error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Contact, error)
	List(ctx context.Context, tenantID uuid.UUID, channelID *uuid.UUID, search string) ([]*entity.Contact, error)
	SetBlocked(ctx context.Context, id uuid.UUID, blocked bool) error
	Update(ctx context.Context, contact *entity.Contact) (*entity.Contact, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type contactRepository struct {
	*baseRepository
}

func NewContactRepository(db *pgxpool.Pool) ContactRepository {
	return &contactRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *contactRepository) FindOrCreate(ctx context.Context, c *entity.Contact) (*entity.Contact, bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.Attributes == nil {
		c.Attributes = &entity.AttributesConfig{}
	}

	// Upsert: existing contact keeps its id; display_name/avatar refreshed.
	// (xmax = 0) is true only for freshly inserted rows -> a first-time contact.
	query := `
		INSERT INTO contacts (id, tenant_id, channel_id, external_id, display_name, avatar_url, attributes)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (channel_id, external_id) DO UPDATE SET
			display_name = COALESCE(NULLIF(EXCLUDED.display_name, ''), contacts.display_name),
			avatar_url = COALESCE(EXCLUDED.avatar_url, contacts.avatar_url),
			updated_at = now()
		RETURNING id, tenant_id, channel_id, external_id, display_name, avatar_url,
			is_blocked, attributes, created_at, updated_at, (xmax = 0) AS created`

	var out entity.Contact
	var created bool
	err := r.q(subCtx).QueryRow(subCtx, query,
		c.ID, c.TenantID, c.ChannelID, c.ExternalID, c.DisplayName, c.AvatarURL, c.Attributes,
	).Scan(
		&out.ID, &out.TenantID, &out.ChannelID, &out.ExternalID, &out.DisplayName,
		&out.AvatarURL, &out.IsBlocked, &out.Attributes, &out.CreatedAt, &out.UpdatedAt, &created,
	)
	if err != nil {
		return nil, false, fmt.Errorf("failed to upsert contact: %w", err)
	}
	return &out, created, nil
}

func (r *contactRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Contact, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var contact entity.Contact
	if err := pgxscan.Get(subCtx, r.q(subCtx), &contact, `SELECT * FROM contacts WHERE id = $1`, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("contact not found")
		}
		return nil, fmt.Errorf("failed to find contact: %w", err)
	}
	return &contact, nil
}

func (r *contactRepository) List(ctx context.Context, tenantID uuid.UUID, channelID *uuid.UUID, search string) ([]*entity.Contact, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM contacts
		WHERE tenant_id = $1
			AND ($2::uuid IS NULL OR channel_id = $2)
			AND ($3 = '' OR display_name ILIKE '%'||$3||'%' OR external_id ILIKE '%'||$3||'%')
		ORDER BY updated_at DESC
		LIMIT 200`

	var contacts []*entity.Contact
	if err := pgxscan.Select(subCtx, r.q(subCtx), &contacts, query, tenantID, channelID, search); err != nil {
		return nil, fmt.Errorf("failed to list contacts: %w", err)
	}
	return contacts, nil
}

func (r *contactRepository) SetBlocked(ctx context.Context, id uuid.UUID, blocked bool) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tag, err := r.q(subCtx).Exec(subCtx, `UPDATE contacts SET is_blocked = $1 WHERE id = $2`, blocked, id)
	if err != nil {
		return fmt.Errorf("failed to set blocked: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("contact not found")
	}
	return nil
}

func (r *contactRepository) Update(ctx context.Context, c *entity.Contact) (*entity.Contact, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE contacts SET display_name = $1, avatar_url = $2, attributes = $3
		WHERE id = $4
		RETURNING id, tenant_id, channel_id, external_id, display_name, avatar_url,
			is_blocked, attributes, created_at, updated_at`
	var out entity.Contact
	err := r.q(subCtx).QueryRow(subCtx, query, c.DisplayName, c.AvatarURL, c.Attributes, c.ID).Scan(
		&out.ID, &out.TenantID, &out.ChannelID, &out.ExternalID, &out.DisplayName,
		&out.AvatarURL, &out.IsBlocked, &out.Attributes, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("contact not found")
		}
		return nil, fmt.Errorf("failed to update contact: %w", err)
	}
	return &out, nil
}

func (r *contactRepository) Delete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tag, err := r.q(subCtx).Exec(subCtx, `DELETE FROM contacts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("contact not found")
	}
	return nil
}
