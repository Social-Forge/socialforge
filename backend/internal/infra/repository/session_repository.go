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

type SessionRepository interface {
	Create(ctx context.Context, session *entity.Session) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Session, error)
	FindByToken(ctx context.Context, token string) (*entity.Session, error)
	Update(ctx context.Context, session *entity.Session) (*entity.Session, error)
	Delete(ctx context.Context, id uuid.UUID) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, opts *ListOptions) ([]*entity.Session, int64, error)
	ListByUser(ctx context.Context, userID uuid.UUID, opts *ListOptions) ([]*entity.Session, int64, error)
	FindActiveByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Session, error)
	InvalidateSession(ctx context.Context, token string) error
	InvalidateAllUserSessions(ctx context.Context, userID uuid.UUID) error
	RevokeSession(ctx context.Context, id uuid.UUID) error
	RefreshSession(ctx context.Context, accToken, refToken string, expiresAt time.Time) error
	CleanupExpiredSessions(ctx context.Context) (int64, error)
	CountActiveByUser(ctx context.Context, userID uuid.UUID) (int64, error)
}

type sessionRepository struct {
	*baseRepository
}

func NewSessionRepository(db *pgxpool.Pool) SessionRepository {
	return &sessionRepository{
		baseRepository: NewBaseRepository(db).(*baseRepository),
	}
}
func (r *sessionRepository) Create(ctx context.Context, session *entity.Session) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		INSERT INTO sessions (
			id, user_id, access_token, refresh_token, ip_address, user_agent,
			expires_at, is_revoked, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (access_token, refresh_token)
		DO NOTHING
		RETURNING id, created_at, updated_at
	`

	args := []interface{}{
		session.ID,
		session.UserID,
		session.AccessToken,
		session.RefreshToken,
		session.IPAddress,
		session.UserAgent,
		session.ExpiresAt,
		session.IsRevoked,
		session.CreatedAt,
	}

	err := r.db.QueryRow(subCtx, query, args...).Scan(
		&session.ID,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}
func (r *sessionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Session, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM sessions
		WHERE id = $1 AND deleted_at IS NULL
	`

	var session entity.Session
	err := pgxscan.Get(subCtx, r.db, &session, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to find session: %w", err)
	}

	return &session, nil
}
func (r *sessionRepository) FindByToken(ctx context.Context, token string) (*entity.Session, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM sessions
		WHERE access_token = $1 OR refresh_token = $1
		  AND is_revoked = false
		  AND expires_at > NOW()
		  AND deleted_at IS NULL
		LIMIT 1
	`

	var session entity.Session
	err := pgxscan.Get(subCtx, r.db, &session, query, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("session not found or expired")
		}
		return nil, fmt.Errorf("failed to find session: %w", err)
	}

	return &session, nil
}
func (r *sessionRepository) Update(ctx context.Context, session *entity.Session) (*entity.Session, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE sessions
		SET ip_address = $1, user_agent = $2, expires_at = $3, is_revoked = $4, last_activity_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING id, user_id, access_token, refresh_token, ip_address, user_agent,
				  expires_at, is_revoked,
				  created_at, updated_at
	`

	args := []interface{}{
		session.IPAddress,
		session.UserAgent,
		session.ExpiresAt,
		session.IsRevoked,
		session.ID,
	}

	var updated entity.Session
	err := r.db.QueryRow(subCtx, query, args...).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.AccessToken,
		&updated.RefreshToken,
		&updated.IPAddress,
		&updated.UserAgent,
		&updated.ExpiresAt,
		&updated.IsRevoked,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "sessions_access_token_unique":
				return nil, fmt.Errorf("duplicate access token: %w", err)
			case "sessions_refresh_token_unique":
				return nil, fmt.Errorf("duplicate refresh token: %w", err)
			default:
				return nil, fmt.Errorf("duplicate session token: %w", err)
			}
		}
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return &updated, nil
}
func (r *sessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE sessions
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	cmdTag, err := r.db.Exec(subCtx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}
func (r *sessionRepository) List(ctx context.Context, opts *ListOptions) ([]*entity.Session, int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if opts == nil {
		opts = NewListOptions()
	}

	countQb := r.buildBaseQuery("SELECT COUNT(*) FROM sessions", opts.Filter)
	countQuery, countArgs := countQb.Build()

	var totalRows int64
	err := r.db.QueryRow(subCtx, countQuery, countArgs...).Scan(&totalRows)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sessions: %w", err)
	}

	if totalRows == 0 {
		return []*entity.Session{}, 0, nil
	}

	qb := r.buildBaseQuery("SELECT * FROM sessions", opts.Filter)

	if opts.OrderBy != "" {
		qb.OrderByField(opts.OrderBy, opts.OrderDir)
	} else {
		qb.OrderByField("created_at", "DESC")
	}

	if opts.Pagination != nil {
		qb.WithLimit(opts.Pagination.Limit)
		qb.WithOffset(opts.Pagination.GetOffset())
	}

	query, args := qb.Build()

	var sessions []*entity.Session
	err = pgxscan.Select(subCtx, r.db, &sessions, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.Session{}, 0, nil
		}
		return nil, 0, fmt.Errorf("failed to list sessions: %w", err)
	}

	return sessions, totalRows, nil
}
func (r *sessionRepository) ListByUser(ctx context.Context, userID uuid.UUID, opts *ListOptions) ([]*entity.Session, int64, error) {
	if opts == nil {
		opts = NewListOptions()
	}
	if opts.Filter == nil {
		opts.Filter = &Filter{}
	}
	opts.Filter.UserID = &userID

	return r.List(ctx, opts)
}
func (r *sessionRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `DELETE FROM sessions WHERE id = $1`

	cmdTag, err := r.db.Exec(subCtx, query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete session: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}
func (r *sessionRepository) FindActiveByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Session, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM sessions
		WHERE user_id = $1
		  AND is_revoked = false
		  AND expires_at > NOW()
		  AND deleted_at IS NULL
		ORDER BY last_activity_at DESC
	`

	var sessions []*entity.Session
	err := pgxscan.Select(subCtx, r.db, &sessions, query, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.Session{}, nil
		}
		return nil, fmt.Errorf("failed to find active sessions: %w", err)
	}

	return sessions, nil
}
func (r *sessionRepository) InvalidateSession(ctx context.Context, token string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE sessions
		SET is_revoked = true
		WHERE token = $1 AND deleted_at IS NULL
	`

	cmdTag, err := r.db.Exec(subCtx, query, token)
	if err != nil {
		return fmt.Errorf("failed to invalidate session: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}
func (r *sessionRepository) InvalidateAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE sessions
		SET is_revoked = true
		WHERE user_id = $1 AND is_revoked = false AND deleted_at IS NULL
	`

	_, err := r.db.Exec(subCtx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to invalidate all sessions: %w", err)
	}

	return nil
}
func (r *sessionRepository) RevokeSession(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE sessions
		SET is_revoked = true
		WHERE id = $1 AND deleted_at IS NULL
	`

	cmdTag, err := r.db.Exec(subCtx, query, id)
	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}
func (r *sessionRepository) RefreshSession(ctx context.Context, accToken, refToken string, expiresAt time.Time) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE sessions
		SET expires_at = $1, last_activity_at = NOW()
		WHERE (access_token = $2 OR refresh_token = $3) AND is_revoked = false AND deleted_at IS NULL
	`

	cmdTag, err := r.db.Exec(subCtx, query, expiresAt, accToken, refToken)
	if err != nil {
		return fmt.Errorf("failed to refresh session: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("session not found or already revoked")
	}

	return nil
}
func (r *sessionRepository) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		DELETE FROM sessions
		WHERE expires_at < NOW() OR is_revoked = true
	`

	cmdTag, err := r.db.Exec(subCtx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	return cmdTag.RowsAffected(), nil
}
func (r *sessionRepository) CountActiveByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT COUNT(*)
		FROM sessions
		WHERE user_id = $1
		  AND is_revoked = false
		  AND expires_at > NOW()
		  AND deleted_at IS NULL
	`

	var count int64
	err := r.db.QueryRow(subCtx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count active sessions: %w", err)
	}

	return count, nil
}
func (r *sessionRepository) buildBaseQuery(baseQuery string, filter *Filter) *QueryBuilder {
	qb := NewQueryBuilder(baseQuery)

	if filter == nil {
		qb.Where("deleted_at IS NULL")
		return qb
	}

	if filter.IncludeDeleted != nil && *filter.IncludeDeleted {
		qb.Where("deleted_at IS NOT NULL")
	} else {
		qb.Where("deleted_at IS NULL")
	}

	if filter.UserID != nil {
		qb.Where("user_id = $?", *filter.UserID)
	}

	if filter.Extra != nil {
		if isRevoked, ok := filter.Extra["is_revoked"].(bool); ok {
			qb.Where("is_revoked = $?", isRevoked)
		}

		if isActive, ok := filter.Extra["is_active"].(bool); ok && isActive {
			qb.Where("expires_at > NOW()")
			qb.Where("is_revoked = false")
		}
	}

	return qb
}
