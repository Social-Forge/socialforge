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

type ConversationRepository interface {
	BaseRepository
	// FindOrCreateOpen returns the active (open/unassigned) conversation for a
	// contact on a channel, creating a new unassigned one if none exists.
	FindOrCreateOpen(ctx context.Context, tenantID, channelID, contactID uuid.UUID) (*entity.Convertation, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Convertation, error)
	TouchInbound(ctx context.Context, id uuid.UUID, at time.Time) error
	TouchOutbound(ctx context.Context, id uuid.UUID, at time.Time) error
	List(ctx context.Context, tenantID uuid.UUID, status string, assignedAgentID *uuid.UUID) ([]*entity.Convertation, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	Assign(ctx context.Context, id uuid.UUID, agentID *uuid.UUID) error
	SetPinned(ctx context.Context, id uuid.UUID, pinned bool) error
	SetArchived(ctx context.Context, id uuid.UUID, archived bool) error
	MarkRead(ctx context.Context, id uuid.UUID) error
	// PickAgentForDivision returns the least-loaded active agent/supervisor in a
	// division (fewest open conversations) for auto-assignment, or nil if none.
	PickAgentForDivision(ctx context.Context, tenantID, divisionID uuid.UUID) (*uuid.UUID, error)
}

type conversationRepository struct {
	*baseRepository
}

func NewConversationRepository(db *pgxpool.Pool) ConversationRepository {
	return &conversationRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *conversationRepository) FindOrCreateOpen(ctx context.Context, tenantID, channelID, contactID uuid.UUID) (*entity.Convertation, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	findQuery := `
		SELECT * FROM conversations
		WHERE contact_id = $1 AND channel_id = $2
			AND status IN ('open', 'unassigned') AND is_archived = false
		ORDER BY created_at DESC LIMIT 1`

	var conv entity.Convertation
	err := pgxscan.Get(subCtx, r.q(subCtx), &conv, findQuery, contactID, channelID)
	if err == nil {
		return &conv, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to find conversation: %w", err)
	}

	insertQuery := `
		INSERT INTO conversations (id, tenant_id, channel_id, contact_id, status)
		VALUES ($1, $2, $3, $4, 'unassigned')
		RETURNING id, tenant_id, channel_id, contact_id, assigned_agent_id, status,
			is_pinned, is_archived, unread_count, last_message_at,
			service_window_expires_at, metadata, created_at, updated_at`

	var created entity.Convertation
	if err := r.q(subCtx).QueryRow(subCtx, insertQuery, uuid.New(), tenantID, channelID, contactID).Scan(
		&created.ID, &created.TenantID, &created.ChannelID, &created.ContactID, &created.AssignedAgentID,
		&created.Status, &created.IsPinned, &created.IsArchived, &created.UnreadCount,
		&created.LastMessageAt, &created.ServiceWindowExpiresAt, &created.Metadata,
		&created.CreatedAt, &created.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}
	return &created, nil
}

func (r *conversationRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Convertation, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var conv entity.Convertation
	if err := pgxscan.Get(subCtx, r.q(subCtx), &conv, `SELECT * FROM conversations WHERE id = $1`, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to find conversation: %w", err)
	}
	return &conv, nil
}

// TouchInbound bumps last_message_at and increments unread on inbound messages.
func (r *conversationRepository) TouchInbound(ctx context.Context, id uuid.UUID, at time.Time) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	_, err := r.q(subCtx).Exec(subCtx,
		`UPDATE conversations SET last_message_at = $1, unread_count = unread_count + 1 WHERE id = $2`, at, id)
	if err != nil {
		return fmt.Errorf("failed to touch conversation: %w", err)
	}
	return nil
}

// TouchOutbound bumps last_message_at without changing unread on outbound.
func (r *conversationRepository) TouchOutbound(ctx context.Context, id uuid.UUID, at time.Time) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	_, err := r.q(subCtx).Exec(subCtx,
		`UPDATE conversations SET last_message_at = $1 WHERE id = $2`, at, id)
	if err != nil {
		return fmt.Errorf("failed to touch conversation: %w", err)
	}
	return nil
}

func (r *conversationRepository) List(ctx context.Context, tenantID uuid.UUID, status string, assignedAgentID *uuid.UUID) ([]*entity.Convertation, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM conversations
		WHERE tenant_id = $1
			AND ($2 = '' OR status = $2)
			AND ($3::uuid IS NULL OR assigned_agent_id = $3)
		ORDER BY COALESCE(last_message_at, created_at) DESC
		LIMIT 200`

	var convs []*entity.Convertation
	if err := pgxscan.Select(subCtx, r.q(subCtx), &convs, query, tenantID, status, assignedAgentID); err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	return convs, nil
}

func (r *conversationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tag, err := r.q(subCtx).Exec(subCtx, `UPDATE conversations SET status = $1 WHERE id = $2`, status, id)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("conversation not found")
	}
	return nil
}

func (r *conversationRepository) SetPinned(ctx context.Context, id uuid.UUID, pinned bool) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	tag, err := r.q(subCtx).Exec(subCtx, `UPDATE conversations SET is_pinned = $1 WHERE id = $2`, pinned, id)
	if err != nil {
		return fmt.Errorf("failed to set pinned: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("conversation not found")
	}
	return nil
}

func (r *conversationRepository) SetArchived(ctx context.Context, id uuid.UUID, archived bool) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	status := ""
	if archived {
		status = entity.ConversationStatusArchived
	}
	// Archiving also moves status to 'archived'; unarchiving reopens as unassigned/open based on assignment.
	query := `UPDATE conversations SET is_archived = $1,
		status = CASE WHEN $1 THEN 'archived'
		              WHEN assigned_agent_id IS NULL THEN 'unassigned'
		              ELSE 'open' END
		WHERE id = $2`
	tag, err := r.q(subCtx).Exec(subCtx, query, archived, id)
	if err != nil {
		return fmt.Errorf("failed to set archived: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("conversation not found")
	}
	_ = status
	return nil
}

func (r *conversationRepository) MarkRead(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	_, err := r.q(subCtx).Exec(subCtx, `UPDATE conversations SET unread_count = 0 WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to mark read: %w", err)
	}
	return nil
}

func (r *conversationRepository) PickAgentForDivision(ctx context.Context, tenantID, divisionID uuid.UUID) (*uuid.UUID, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT ut.user_id
		FROM division_members dm
		JOIN user_tenants ut ON ut.id = dm.user_tenant_id AND ut.is_active = true AND ut.deleted_at IS NULL
		JOIN roles r ON r.id = ut.role_id AND r.slug IN ('agent', 'supervisor')
		LEFT JOIN conversations c ON c.assigned_agent_id = ut.user_id
			AND c.tenant_id = $1 AND c.status = 'open' AND c.is_archived = false
		WHERE dm.division_id = $2 AND dm.is_active = true AND dm.deleted_at IS NULL
		GROUP BY ut.user_id
		ORDER BY COUNT(c.id) ASC, random()
		LIMIT 1`

	var agentID uuid.UUID
	err := r.q(subCtx).QueryRow(subCtx, query, tenantID, divisionID).Scan(&agentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // no active agent -> leave unassigned
		}
		return nil, fmt.Errorf("failed to pick agent: %w", err)
	}
	return &agentID, nil
}

func (r *conversationRepository) Assign(ctx context.Context, id uuid.UUID, agentID *uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	status := entity.ConversationStatusUnassigned
	if agentID != nil {
		status = entity.ConversationStatusOpen
	}
	tag, err := r.q(subCtx).Exec(subCtx,
		`UPDATE conversations SET assigned_agent_id = $1, status = $2 WHERE id = $3`, agentID, status, id)
	if err != nil {
		return fmt.Errorf("failed to assign conversation: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("conversation not found")
	}
	return nil
}
