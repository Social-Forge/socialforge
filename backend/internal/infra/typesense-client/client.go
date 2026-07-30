package typesenseclient

import (
	"context"
	"errors"
	"github/socialforge/config"
	"github/socialforge/internal/infra/contextpool"
	"sync"
	"time"

	"github.com/typesense/typesense-go/v4/typesense"
	"go.uber.org/zap"
)

var (
	TypesenseStorage *TypesenseClient
	typesenseOnce    sync.Once
)

type TypesenseClient struct {
	client *typesense.Client
	config *config.TypeSenseConfig
	logger *zap.Logger
	isUp   bool
	mu     sync.RWMutex
}

func NewTypesenseClient(ctx context.Context, cfg *config.TypeSenseConfig, logger *zap.Logger) (*TypesenseClient, error) {
	var initErr error
	typesenseOnce.Do(func() {
		if cfg == nil {
			initErr = errors.New("typesense config is required")
			return
		}
		if cfg.ApiKey == "" {
			initErr = errors.New("typesense api key is required")
			logger.Error("Typesense configuration missing: API key is empty")
			return
		}
		if cfg.Host == "" || cfg.Port == 0 {
			initErr = errors.New("typesense host and port are required")
			logger.Error("Typesense configuration missing: host or port is empty")
			return
		}

		client := typesense.NewClient(
			typesense.WithServer(cfg.GetURL()),
			typesense.WithAPIKey(cfg.ApiKey),
		)

		checkCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 5*time.Second)
		defer cancel()

		TypesenseStorage = &TypesenseClient{
			client: client,
			config: cfg,
			logger: logger,
			isUp:   true,
		}

		if _, err := client.Debug(checkCtx); err != nil {
			TypesenseStorage.mu.Lock()
			TypesenseStorage.isUp = false
			TypesenseStorage.mu.Unlock()
			logger.Warn("Typesense connection check failed; client initialized but marked unavailable", zap.Error(err))
		}

		logger.Info("Typesense client initialized successfully",
			zap.String("url", cfg.GetURL()),
		)
	})

	if initErr != nil {
		return nil, initErr
	}
	return TypesenseStorage, nil
}

func GetTypesense() (*TypesenseClient, error) {
	if TypesenseStorage == nil {
		return nil, errors.New("typesense not initialized: call NewTypesenseClient first")
	}
	return TypesenseStorage, nil
}

func (ts *TypesenseClient) Client() *typesense.Client {
	return ts.client
}

func (ts *TypesenseClient) IsUp() bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.isUp
}
