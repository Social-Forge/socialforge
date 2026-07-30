package services

import (
	"context"
	"errors"
	"fmt"
	"github/socialforge/internal/dto"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/infra/contextpool"
	minioclient "github/socialforge/internal/infra/minio-client"
	"github/socialforge/internal/infra/repository"
	"github/socialforge/internal/utils"
	"mime/multipart"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	tenantRepo     repository.TenantRepository
	divisiRepo     repository.DivisionRepository
	userTenantRepo repository.UserTenantRepository
	userHelper     *helpers.UserHelper
	tokenHelper    *helpers.TokenHelper
	logger         *zap.Logger
	minio          *minioclient.MinioClient
}

func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	tenantRepo repository.TenantRepository,
	divisiRepo repository.DivisionRepository,
	userTenantRepo repository.UserTenantRepository,
	userHelper *helpers.UserHelper,
	tokenHelper *helpers.TokenHelper,
	logger *zap.Logger,
	minio *minioclient.MinioClient,
) *UserService {
	return &UserService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		tenantRepo:     tenantRepo,
		divisiRepo:     divisiRepo,
		userTenantRepo: userTenantRepo,
		userHelper:     userHelper,
		tokenHelper:    tokenHelper,
		logger:         logger,
		minio:          minio,
	}
}
func (s *UserService) GetUserByID(ctx context.Context, userID string) (*entity.UserResponse, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	userTenant, err := s.userRepo.GetUserTenantWithDetailsByUserID(subCtx, userIDUUID)
	if err != nil {
		s.logger.Error("Failed to get user tenant",
			zap.String("user_id", userTenant.UserTenant.UserID.String()), // ✅ Log user ID yang benar
			zap.String("expected_user_id", userIDUUID.String()),
			zap.Error(err),
		)
		return nil, dto.ErrUserNotFound
	}
	s.logger.Info("User found",
		zap.String("user_id", userID),
	)

	return userTenant.ToResponse(), nil
}
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*entity.UserResponse, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userEmail, err := s.userRepo.FindByEmail(subCtx, email)
	if err != nil {
		return nil, dto.ErrUserNotFound
	}
	userTenant, err := s.userRepo.GetUserTenantWithDetailsByUserID(subCtx, userEmail.ID)
	if err != nil {
		return nil, dto.ErrUserNotFound
	}
	s.logger.Info("User found",
		zap.String("email", email),
	)

	return userTenant.ToResponse(), nil
}
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*entity.UserResponse, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userUsername, err := s.userRepo.FindByUsername(subCtx, username)
	if err != nil {
		return nil, dto.ErrUserNotFound
	}
	userTenant, err := s.userRepo.GetUserTenantWithDetailsByUserID(subCtx, userUsername.ID)
	if err != nil {
		return nil, dto.ErrUserNotFound
	}
	s.logger.Info("User found",
		zap.String("username", username),
	)

	return userTenant.ToResponse(), nil
}
func (s *UserService) Logout(ctx context.Context, userID string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.FindByID(subCtx, userUUID)
	if err != nil || user == nil {
		return dto.ErrUserNotFound
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.tokenHelper.DeleteAllUserSessions(subCtx, user.ID); err != nil {
			errChan <- fmt.Errorf("failed to clear auth cache: %w", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.userHelper.ClearAllSessionCache(subCtx, userID); err != nil {
			errChan <- fmt.Errorf("failed to clear user cache: %w", err)
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return fmt.Errorf("failed to logout user: %w", err)
		}
	}

	return nil
}
func (s *UserService) ChangeAvatar(ctx context.Context, userID string, avatarFile *multipart.FileHeader) (string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.FindByID(subCtx, userUUID)
	if err != nil || user == nil {
		return "", dto.ErrUserNotFound
	}

	objectName := fmt.Sprintf("user/%s/avatar", userID)

	var wg sync.WaitGroup
	if user.AvatarURL.Valid && user.AvatarURL.String != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Gunakan background context untuk cleanup ops
			cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cleanupCancel()
			if errDelete := s.minio.DeleteFile(cleanupCtx, objectName); errDelete != nil {
				s.logger.Error("Failed to delete old avatar",
					zap.String("user_id", userID),
					zap.Error(errDelete),
				)
			}
		}()
	}

	file, err := avatarFile.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open avatar file: %w", err)
	}
	defer file.Close()

	_, err = s.minio.UploadImage(subCtx, objectName, file, avatarFile.Size)
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar: %w", err)
	}

	presignedURL, err := s.minio.GetPresignedURL(ctx, objectName, time.Hour*24*7)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	wg.Wait()

	newAvatarURL, err := s.userRepo.UpdateAvatar(subCtx, userUUID, presignedURL)
	if err != nil {
		rollbackCtx, rollbackCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer rollbackCancel()
		if errRollback := s.minio.DeleteFile(rollbackCtx, objectName); errRollback != nil {
			s.logger.Error("Failed to rollback delete avatar",
				zap.String("user_id", userID),
				zap.Error(errRollback),
			)
		}

		return "", fmt.Errorf("failed to update avatar: %w", err)
	}

	return newAvatarURL, nil
}
func (s *UserService) GetAvatarURL(ctx context.Context, userID string) (string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.FindByID(subCtx, userUUID)
	if err != nil || user == nil {
		return "", dto.ErrUserNotFound
	}
	if !user.AvatarURL.Valid || user.AvatarURL.String == "" {
		return "", errors.New("user has no avatar")
	}

	presignedUrl, err := s.minio.GetPresignedURL(ctx, user.AvatarURL.String, time.Hour*24*7)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedUrl, nil
}
func (s *UserService) UpdateProfile(ctx context.Context, userID string, req *dto.UpdateProfileRequest) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.FindByID(subCtx, userUUID)
	if err != nil || user == nil {
		return nil, dto.ErrUserNotFound
	}

	payload := &entity.User{
		ID:              user.ID,
		FullName:        req.FullName,
		Email:           req.Email,
		Username:        req.Username,
		Phone:           entity.NewNullString(req.Phone),
		AvatarURL:       entity.NewNullString(user.AvatarURL.String),
		IsActive:        user.IsActive,
		IsVerified:      user.IsVerified,
		EmailVerifiedAt: user.EmailVerifiedAt,
	}
	userUpdate, err := s.userRepo.Update(subCtx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	userUpdate.PasswordHash = "[Hidden]"
	userUpdate.TwoFaSecret = entity.NewNullString(user.TwoFaSecret.String)

	return userUpdate, nil
}
func (s *UserService) UpdatePassword(ctx context.Context, userID string, req *dto.ChangePasswordRequest) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.FindByID(subCtx, userUUID)
	if err != nil || user == nil {
		return dto.ErrUserNotFound
	}

	if utils.VerifyPassword(user.PasswordHash, req.NewPassword) {
		return fmt.Errorf("new password cannot be the same as the old password")
	}

	if !utils.IsStrongPassword(req.NewPassword) {
		return fmt.Errorf("new password is not strong enough")
	}

	hashedPassword, err := utils.GeneratePasswordHash(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.userRepo.UpdatePassword(subCtx, userUUID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
func (s *UserService) EnableTwoFactor(ctx context.Context, userID string, req *dto.EnableTwoFactorRequest) (string, string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", "", fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.FindByID(subCtx, userUUID)
	if err != nil || user == nil {
		return "", "", dto.ErrUserNotFound
	}

	switch req.Status {
	case true:
		if user.TwoFaSecret.Valid && user.TwoFaSecret.String != "" {
			return "", "", fmt.Errorf("two-factor authentication is already enabled")
		}
		s.userHelper.ClearTemp2FASecret(subCtx, userUUID)

		qrcode, secret, err := s.userHelper.Generate2FAQRCode(user.ID, user.FullName)
		if err != nil {
			return "", "", fmt.Errorf("failed to generate two-factor authentication secret: %w", err)
		}
		payload := map[string]interface{}{
			"qr_url":  qrcode,
			"secret":  secret,
			"user_id": user.ID,
		}
		if err := s.userHelper.SetTemp2FASecret(subCtx, userUUID, payload); err != nil {
			return "", "", fmt.Errorf("failed to set two-factor authentication status: %w", err)
		}

		return qrcode, secret, nil
	case false:
		if !user.TwoFaSecret.Valid || user.TwoFaSecret.String == "" {
			return "", "", fmt.Errorf("two-factor authentication is not enabled")
		}
		if err := s.userRepo.UpdateTwoFaSecret(subCtx, userUUID, nil); err != nil {
			return "", "", fmt.Errorf("failed to disable two-factor authentication: %w", err)
		}
		if err := s.userHelper.ClearTemp2FASecret(subCtx, userUUID); err != nil {
			return "", "", fmt.Errorf("failed to clear two-factor authentication status: %w", err)
		}
		return "", "", nil
	default:
		return "", "", fmt.Errorf("invalid status value")
	}
}
func (s *UserService) ActivateTwoFactor(ctx context.Context, userID string, req *dto.ActivateTwoFactorRequest) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.FindByID(subCtx, userUUID)
	if err != nil || user == nil {
		return dto.ErrUserNotFound
	}

	if user.TwoFaSecret.Valid || user.TwoFaSecret.String != "" {
		return fmt.Errorf("two-factor authentication already enabled")
	}

	payload, err := s.userHelper.GetTemp2FASecret(subCtx, userUUID)
	if err != nil {
		return fmt.Errorf("failed to get two-factor authentication status: %w", err)
	}
	secret, ok := payload["secret"].(string)
	if !ok || secret == "" {
		return fmt.Errorf("two-factor authentication secret not found")
	}

	valid, err := s.userHelper.Verify2FA(subCtx, userUUID, req.Code, secret)
	if err != nil {
		return fmt.Errorf("failed to verify two-factor authentication token: %w", err)
	}
	if !valid {
		return fmt.Errorf("invalid two-factor authentication token")
	}

	if err := s.userRepo.UpdateTwoFaSecret(subCtx, userUUID, &secret); err != nil {
		return fmt.Errorf("failed to update two-factor authentication secret: %w", err)
	}
	if err := s.userHelper.ClearTemp2FASecret(subCtx, userUUID); err != nil {
		return fmt.Errorf("failed to clear two-factor authentication status: %w", err)
	}
	return nil
}
