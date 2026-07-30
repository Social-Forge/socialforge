package services

import (
	"context"
	"fmt"
	"github/socialforge/internal/dto"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	minioclient "github/socialforge/internal/infra/minio-client"
	"github/socialforge/internal/infra/repository"
	"mime/multipart"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TenantService struct {
	tenantRepo repository.TenantRepository
	logger     *zap.Logger
	minio      *minioclient.MinioClient
}

func NewTenantService(tenantRepo repository.TenantRepository, logger *zap.Logger, minio *minioclient.MinioClient) *TenantService {
	return &TenantService{
		tenantRepo: tenantRepo,
		logger:     logger,
		minio:      minio,
	}
}
func (s *TenantService) UpdateInfo(ctx context.Context, tenantID string, req *dto.UpdateTenantRequest) (*entity.Tenant, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant ID: %w", err)
	}

	existTenant, err := s.tenantRepo.FindByID(subCtx, tenantUUID)
	if err != nil {
		return nil, fmt.Errorf("tenant not found: %w", err)
	}

	payload := &entity.Tenant{
		ID:                 existTenant.ID,
		Name:               req.Name,
		Slug:               req.Slug,
		LogoURL:            existTenant.LogoURL,
		Description:        entity.NewNullString(req.Description),
		SubscriptionPlan:   existTenant.SubscriptionPlan,
		SubscriptionStatus: existTenant.SubscriptionStatus,
		TrialEndsAt:        existTenant.TrialEndsAt,
		IsActive:           existTenant.IsActive,
	}

	updateTenant, err := s.tenantRepo.Update(subCtx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	return updateTenant, nil
}
func (s *TenantService) ChangeLogo(ctx context.Context, tenantID string, logoFile *multipart.FileHeader) (string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return "", fmt.Errorf("invalid tenant ID: %w", err)
	}

	existTenant, err := s.tenantRepo.FindByID(subCtx, tenantUUID)
	if err != nil {
		return "", fmt.Errorf("tenant not found: %w", err)
	}
	objectName := fmt.Sprintf("tenant/%s/logo", tenantID)

	var wg sync.WaitGroup
	if existTenant.LogoURL.Valid && existTenant.LogoURL.String != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cleanupCancel()
			if errDelete := s.minio.DeleteFile(cleanupCtx, objectName); errDelete != nil {
				s.logger.Error("Failed to delete old logo",
					zap.String("tenant_id", tenantID),
					zap.Error(errDelete),
				)
			}
		}()
	}

	file, err := logoFile.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open logo file: %w", err)
	}
	defer file.Close()

	_, err = s.minio.UploadImage(subCtx, objectName, file, logoFile.Size)
	if err != nil {
		return "", fmt.Errorf("failed to upload logo: %w", err)
	}

	presignedURL, err := s.minio.GetPresignedURL(ctx, objectName, time.Hour*24*7)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	wg.Wait()

	newLogoUrl, err := s.tenantRepo.UpdateLogo(subCtx, tenantUUID, presignedURL)
	if err != nil {
		rollbackCtx, rollbackCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer rollbackCancel()
		if errRollback := s.minio.DeleteFile(rollbackCtx, objectName); errRollback != nil {
			s.logger.Error("Failed to rollback delete logo",
				zap.String("tenant_id", tenantID),
				zap.Error(errRollback),
			)
		}

		return "", fmt.Errorf("failed to update logo: %w", err)
	}

	return newLogoUrl, nil
}
