package minioclient

import (
	"context"
	"errors"
	"fmt"
	"github/socialforge/config"
	"github/socialforge/internal/infra/contextpool"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

var (
	MinioStorage *MinioClient
	minioOnce    sync.Once
)

type MinioClient struct {
	client     *minio.Client
	config     *config.MinIOConfig
	logger     *zap.Logger
	bucketName string
	isUp       bool
	mu         sync.RWMutex
}

func NewMinioClient(ctx context.Context, cfg *config.MinIOConfig, logger *zap.Logger) (*MinioClient, error) {
	var initErr error
	minioOnce.Do(func() {
		if cfg.Endpoint == "" {
			initErr = errors.New("minio endpoint is required")
			config.Logger.Error("MinIO configuration missing")
			return
		}
		if cfg.AccessKey == "" || cfg.SecretKey == "" {
			initErr = errors.New("minio access key and secret key are required")
			config.Logger.Error("MinIO configuration missing")
			return
		}
		if cfg.BucketName == "" {
			initErr = errors.New("minio bucket name is required")
			config.Logger.Error("MinIO configuration missing")
			return
		}
		client, err := minio.New(cfg.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: cfg.UseSSL,
		})
		if err != nil {
			initErr = fmt.Errorf("failed to create MinIO client: %w", err)
			config.Logger.Error("MinIO client creation failed", zap.Error(err))
			return
		}

		subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 10*time.Second)
		defer cancel()

		exists, err := client.BucketExists(subCtx, cfg.BucketName)
		if err != nil {
			initErr = fmt.Errorf("failed to check bucket existence: %w", err)
			config.Logger.Error("MinIO bucket check failed", zap.Error(err))
			return
		}

		if !exists {
			err = client.MakeBucket(subCtx, cfg.BucketName, minio.MakeBucketOptions{})
			if err != nil {
				initErr = fmt.Errorf("failed to create bucket: %w", err)
				config.Logger.Error("MinIO bucket creation failed", zap.Error(err))
				return
			}
			config.Logger.Info("MinIO bucket created", zap.String("bucket", cfg.BucketName))
		}

		MinioStorage = &MinioClient{
			client:     client,
			config:     cfg,
			logger:     config.Logger,
			bucketName: cfg.BucketName,
			isUp:       true,
		}

		config.Logger.Info("✅ MinIO client initialized successfully",
			zap.String("endpoint", cfg.Endpoint),
			zap.String("bucket", cfg.BucketName),
			zap.Bool("ssl", cfg.UseSSL),
		)
	})
	if initErr != nil {
		return nil, initErr
	}
	return MinioStorage, nil
}
func GetMinio() (*MinioClient, error) {
	if MinioStorage == nil {
		return nil, errors.New("minio not initialized: call NewMinioClient first")
	}
	return MinioStorage, nil
}
