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

type TokenRepository interface {
	Create(ctx context.Context, token *entity.Token) (*entity.Token, error)
	CreateOrGetExist(ctx context.Context, token *entity.Token) (*entity.Token, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Token, error)
	FindByToken(ctx context.Context, token string) (*entity.Token, error)
	Update(ctx context.Context, token *entity.Token) (*entity.Token, error)
	Delete(ctx context.Context, id uuid.UUID) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	HardDeleteByToken(ctx context.Context, token string) error
	List(ctx context.Context, opts *ListOptions) ([]*entity.Token, int64, error)
	ListByUser(ctx context.Context, userID uuid.UUID, opts *ListOptions) ([]*entity.Token, int64, error)
	FindActiveByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Token, error)
	FindByTokenAndType(ctx context.Context, token, tokenType string) (*entity.Token, error)
	FindByType(ctx context.Context, userID uuid.UUID, tokenType string) (*entity.Token, error)
	RevokeToken(ctx context.Context, id uuid.UUID) error
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID, tokenType string) error
	ValidateToken(ctx context.Context, token, tokenType string) (*entity.Token, error)
	CleanupExpiredTokens(ctx context.Context) (int64, error)
	CountActiveByUser(ctx context.Context, userID uuid.UUID, tokenType string) (int64, error)
}

type tokenRepository struct {
	*baseRepository
}

func NewTokenRepository(db *pgxpool.Pool) TokenRepository {
	return &tokenRepository{
		baseRepository: NewBaseRepository(db).(*baseRepository),
	}
}
func (r *tokenRepository) Create(ctx context.Context, token *entity.Token) (*entity.Token, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		INSERT INTO tokens (
			id, user_id, token, type, expires_at,
			is_used, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (token) DO UPDATE SET
			is_used = EXCLUDED.is_used
		RETURNING id, user_id, token, type, expires_at, is_used, created_at, updated_at
	`

	args := []interface{}{
		token.ID,
		token.UserID,
		token.Token,
		token.Type,
		token.ExpiresAt,
		token.IsUsed,
		token.CreatedAt,
	}

	var newToken entity.Token
	err := r.db.QueryRow(subCtx, query, args...).Scan(
		&newToken.ID,
		&newToken.UserID,
		&newToken.Token,
		&newToken.Type,
		&newToken.ExpiresAt,
		&newToken.IsUsed,
		&newToken.CreatedAt,
		&newToken.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("duplicate token: %w", err)
		}
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return &newToken, nil
}
func (r *tokenRepository) CreateOrGetExist(ctx context.Context, token *entity.Token) (*entity.Token, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	existToken, err := r.FindByType(subCtx, token.UserID, token.Type)
	if err == nil && existToken != nil {
		return existToken, nil
	}
	newToken, errCreate := r.Create(subCtx, token)
	if errCreate != nil {
		return nil, fmt.Errorf("failed to create token: %w", errCreate)
	}
	return newToken, nil
}

func (r *tokenRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Token, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM tokens
		WHERE id = $1 AND deleted_at IS NULL
	`

	var token entity.Token
	err := pgxscan.Get(subCtx, r.db, &token, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to find token: %w", err)
	}

	return &token, nil
}
func (r *tokenRepository) FindByToken(ctx context.Context, token string) (*entity.Token, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM tokens
		WHERE token = $1 AND deleted_at IS NULL
		LIMIT 1
	`

	var tokenEntity entity.Token
	err := pgxscan.Get(subCtx, r.db, &tokenEntity, query, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to find token: %w", err)
	}

	return &tokenEntity, nil
}
func (r *tokenRepository) Update(ctx context.Context, token *entity.Token) (*entity.Token, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE tokens
		SET expires_at = $1, is_used = $2
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING id, user_id, token, type, expires_at, is_used, created_at, updated_at
	`

	args := []interface{}{
		token.ExpiresAt,
		token.IsUsed,
		token.ID,
	}

	var updated entity.Token
	err := r.db.QueryRow(subCtx, query, args...).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.Token,
		&updated.Type,
		&updated.ExpiresAt,
		&updated.IsUsed,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to update token: %w", err)
	}

	return &updated, nil
}
func (r *tokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE tokens
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	cmdTag, err := r.db.Exec(subCtx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("token not found")
	}

	return nil
}
func (r *tokenRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `DELETE FROM tokens WHERE id = $1`

	cmdTag, err := r.db.Exec(subCtx, query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete token: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("token not found")
	}

	return nil
}
func (r *tokenRepository) HardDeleteByToken(ctx context.Context, token string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `DELETE FROM tokens WHERE token = $1`

	cmdTag, err := r.db.Exec(subCtx, query, token)
	if err != nil {
		return fmt.Errorf("failed to hard delete token: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("token not found")
	}

	return nil
}
func (r *tokenRepository) List(ctx context.Context, opts *ListOptions) ([]*entity.Token, int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	if opts == nil {
		opts = NewListOptions()
	}

	// Count total
	countQb := r.buildBaseQuery("SELECT COUNT(*) FROM tokens", opts.Filter)
	countQuery, countArgs := countQb.Build()

	var totalRows int64
	err := r.db.QueryRow(subCtx, countQuery, countArgs...).Scan(&totalRows)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tokens: %w", err)
	}

	if totalRows == 0 {
		return []*entity.Token{}, 0, nil
	}

	// Get data
	qb := r.buildBaseQuery("SELECT * FROM tokens", opts.Filter)

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

	var tokens []*entity.Token
	err = pgxscan.Select(subCtx, r.db, &tokens, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.Token{}, 0, nil
		}
		return nil, 0, fmt.Errorf("failed to list tokens: %w", err)
	}

	return tokens, totalRows, nil
}
func (r *tokenRepository) ListByUser(ctx context.Context, userID uuid.UUID, opts *ListOptions) ([]*entity.Token, int64, error) {
	if opts == nil {
		opts = NewListOptions()
	}
	if opts.Filter == nil {
		opts.Filter = &Filter{}
	}
	opts.Filter.UserID = &userID

	return r.List(ctx, opts)
}
func (r *tokenRepository) FindActiveByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Token, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM tokens
		WHERE user_id = $1
		  AND is_used = false
		  AND expires_at > NOW()
		  AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	var tokens []*entity.Token
	err := pgxscan.Select(subCtx, r.db, &tokens, query, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.Token{}, nil
		}
		return nil, fmt.Errorf("failed to find active tokens: %w", err)
	}

	return tokens, nil
}
func (r *tokenRepository) FindByTokenAndType(ctx context.Context, token, tokenType string) (*entity.Token, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM tokens
		WHERE token = $1 
		  AND type = $2
		  AND deleted_at IS NULL
		LIMIT 1
	`

	var tokenEntity entity.Token
	err := pgxscan.Get(subCtx, r.db, &tokenEntity, query, token, tokenType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to find token: %w", err)
	}

	return &tokenEntity, nil
}
func (r *tokenRepository) FindByType(ctx context.Context, userID uuid.UUID, tokenType string) (*entity.Token, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM tokens
		WHERE user_id = $1
		  AND type = $2
		  AND is_used = false
		  AND expires_at > NOW()
		  AND deleted_at IS NULL
		LIMIT 1
	`

	var tokenEntity entity.Token
	err := pgxscan.Get(subCtx, r.db, &tokenEntity, query, userID, tokenType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to find token: %w", err)
	}

	return &tokenEntity, nil
}
func (r *tokenRepository) RevokeToken(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		UPDATE tokens
		SET is_used = true
		WHERE id = $1 AND deleted_at IS NULL
	`

	cmdTag, err := r.db.Exec(subCtx, query, id)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("token not found or already used")
	}

	return nil
}
func (r *tokenRepository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID, tokenType string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var query string
	var args []interface{}

	if tokenType != "" {
		query = `
			UPDATE tokens
			SET is_used = true, deleted_at = NOW()
			WHERE user_id = $1 
			  AND type = $2
			  AND is_used = false 
			  AND deleted_at IS NULL
		`
		args = []interface{}{userID, tokenType}
	} else {
		query = `
			UPDATE tokens
			SET is_used = true, deleted_at = NOW()
			WHERE user_id = $1 
			  AND is_used = false 
			  AND deleted_at IS NULL
		`
		args = []interface{}{userID}
	}

	cmdTag, err := r.db.Exec(subCtx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to revoke all tokens: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("token not found or already used")
	}

	return nil
}
func (r *tokenRepository) ValidateToken(ctx context.Context, token, tokenType string) (*entity.Token, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT * FROM tokens
		WHERE token = $1 
		  AND type = $2
		  AND is_used = false
		  AND expires_at > NOW()
		  AND deleted_at IS NULL
		LIMIT 1
	`

	var tokenEntity entity.Token
	err := pgxscan.Get(subCtx, r.db, &tokenEntity, query, token, tokenType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("token not found or expired")
		}
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	updateQuery := `
		UPDATE tokens
		SET is_used = true
		WHERE id = $1
	`
	_, _ = r.db.Exec(subCtx, updateQuery, tokenEntity.ID)

	return &tokenEntity, nil
}
func (r *tokenRepository) CleanupExpiredTokens(ctx context.Context) (int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		DELETE FROM tokens
		WHERE expires_at < NOW() OR is_used = true
	`

	cmdTag, err := r.db.Exec(subCtx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	return cmdTag.RowsAffected(), nil
}
func (r *tokenRepository) CountActiveByUser(ctx context.Context, userID uuid.UUID, tokenType string) (int64, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	var query string
	var args []interface{}

	if tokenType != "" {
		query = `
			SELECT COUNT(*)
			FROM tokens
			WHERE user_id = $1
			  AND type = $2
			  AND is_used = false
			  AND expires_at > NOW()
			  AND deleted_at IS NULL
		`
		args = []interface{}{userID, tokenType}
	} else {
		query = `
			SELECT COUNT(*)
			FROM tokens
			WHERE user_id = $1
			  AND is_used = false
			  AND expires_at > NOW()
			  AND deleted_at IS NULL
		`
		args = []interface{}{userID}
	}

	var count int64
	err := r.db.QueryRow(subCtx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count active tokens: %w", err)
	}

	return count, nil
}

func (r *tokenRepository) buildBaseQuery(baseQuery string, filter *Filter) *QueryBuilder {
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
		if tokenType, ok := filter.Extra["type"].(string); ok && tokenType != "" {
			qb.Where("type = $?", tokenType)
		}

		if isUsed, ok := filter.Extra["is_used"].(bool); ok {
			qb.Where("is_used = $?", isUsed)
		}

		if isActive, ok := filter.Extra["is_active"].(bool); ok && isActive {
			qb.Where("expires_at > NOW()")
			qb.Where("is_used = false")
		}
	}

	return qb
}
