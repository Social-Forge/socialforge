package oauth

import (
	"context"
	"fmt"
	"github/socialforge/config"
	redisclient "github/socialforge/internal/infra/redis-client"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"go.uber.org/zap"
)

type OAuthClient struct {
	client       *redisclient.RedisClient
	sessionStore *fiber.Handler
}

func NewOAuthClient(client *redisclient.RedisClient, cfg *config.Config, logger *zap.Logger) *OAuthClient {
	sessionStore := session.New(session.Config{
		Storage: &redisStorageWrapper{client: client},
	})

	goth.UseProviders(
		google.New(cfg.OAuth.GoogleClientId, cfg.OAuth.GoogleClientSecret, fmt.Sprintf("%s%s", cfg.App.URL, "/auth/google/callback"), "email", "profile"),
		facebook.New(cfg.OAuth.FacebookClientId, cfg.OAuth.FacebookClientSecret, fmt.Sprintf("%s%s", cfg.App.URL, "/auth/facebook/callback"), "email"),
		github.New(cfg.OAuth.GithubClientId, cfg.OAuth.GithubClientSecret, fmt.Sprintf("%s%s", cfg.App.URL, "/auth/github/callback"), "user:email"),
	)

	return &OAuthClient{
		client:       client,
		sessionStore: &sessionStore,
	}
}

type redisStorageWrapper struct {
	client *redisclient.RedisClient
}

func (w *redisStorageWrapper) Get(key string) ([]byte, error) {
	return w.client.Get(key)
}

func (w *redisStorageWrapper) Set(key string, val []byte, exp time.Duration) error {
	return w.client.Set(key, val, exp)
}

func (w *redisStorageWrapper) Delete(key string) error {
	return w.client.Delete(key)
}

func (w *redisStorageWrapper) Reset() error {
	return w.client.Reset()
}

func (w *redisStorageWrapper) Close() error {
	return w.client.Close()
}

func (w *redisStorageWrapper) GetWithContext(ctx context.Context, key string) ([]byte, error) {
	return w.client.GetByte(ctx, key)
}

func (w *redisStorageWrapper) SetWithContext(ctx context.Context, key string, val []byte, exp time.Duration) error {
	return w.client.Setbyte(ctx, key, val, exp)
}

func (w *redisStorageWrapper) DeleteWithContext(ctx context.Context, key string) error {
	return w.client.DeleteCache(ctx, key)
}

func (w *redisStorageWrapper) ResetWithContext(ctx context.Context) error {
	return w.client.Reset()
}

func (w *redisStorageWrapper) CloseWithContext(ctx context.Context) error {
	return w.client.CloseClient(ctx)
}
