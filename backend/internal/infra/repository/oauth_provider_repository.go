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

type OAuthProviderRepository interface {
	Create(ctx context.Context, provider *entity.OAuthProvider) error
	FindByProviderAndID(ctx context.Context, providerName, providerID string) (*entity.OAuthProvider, error)
	FindByUserIDAndProvider(ctx context.Context, userID uuid.UUID, providerName string) (*entity.OAuthProvider, error)
	DeleteByUserIDAndProvider(ctx context.Context, userID uuid.UUID, providerName string) error
}

type oauthProviderRepository struct {
	*baseRepository
}

func NewOAuthProviderRepository(db *pgxpool.Pool) OAuthProviderRepository {
	return &oauthProviderRepository{
		baseRepository: NewBaseRepository(db).(*baseRepository),
	}
}

func (r *oauthProviderRepository) Create(ctx context.Context, provider *entity.OAuthProvider) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		INSERT INTO oauth_providers (
			id, user_id, provider_name, provider_id, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(subCtx, query,
		provider.ID,
		provider.UserID,
		provider.ProviderName,
		provider.ProviderID,
		provider.CreatedAt,
		provider.UpdatedAt,
	).Scan(&provider.ID, &provider.CreatedAt, &provider.UpdatedAt)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			return fmt.Errorf("oauth provider link already exists: %w", err)
		}
		return fmt.Errorf("failed to create oauth provider: %w", err)
	}

	return nil
}

func (r *oauthProviderRepository) FindByProviderAndID(ctx context.Context, providerName, providerID string) (*entity.OAuthProvider, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, provider_name, provider_id, created_at, updated_at
		FROM oauth_providers
		WHERE provider_name = $1 AND provider_id = $2
		LIMIT 1
	`

	var provider entity.OAuthProvider
	err := pgxscan.Get(subCtx, r.db, &provider, query, providerName, providerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find oauth provider: %w", err)
	}

	return &provider, nil
}

func (r *oauthProviderRepository) FindByUserIDAndProvider(ctx context.Context, userID uuid.UUID, providerName string) (*entity.OAuthProvider, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, provider_name, provider_id, created_at, updated_at
		FROM oauth_providers
		WHERE user_id = $1 AND provider_name = $2
		LIMIT 1
	`

	var provider entity.OAuthProvider
	err := pgxscan.Get(subCtx, r.db, &provider, query, userID, providerName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find oauth provider by user: %w", err)
	}

	return &provider, nil
}

func (r *oauthProviderRepository) DeleteByUserIDAndProvider(ctx context.Context, userID uuid.UUID, providerName string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	query := `
		DELETE FROM oauth_providers
		WHERE user_id = $1 AND provider_name = $2
	`

	cmdTag, err := r.db.Exec(subCtx, query, userID, providerName)
	if err != nil {
		return fmt.Errorf("failed to delete oauth provider: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("oauth provider not found")
	}

	return nil
}
