package services

import (
	"context"
	"errors"
	"fmt"
	"github/socialforge/internal/dto"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/repository"
	typesenseclient "github/socialforge/internal/infra/typesense-client"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/utils"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/markbates/goth"
	"go.uber.org/zap"
)

type AuthService struct {
	userRepo      repository.UserRepository
	oauthRepo     repository.OAuthProviderRepository
	roleRepo      repository.RoleRepository
	sessionRepo   repository.SessionRepository
	tenantRepo    repository.TenantRepository
	tokenRepo     repository.TokenRepository
	searchIndex   *typesenseclient.SearchIndexService
	rateLimiter   *middlewares.RateLimiterMiddleware
	tokenHelper   *helpers.TokenHelper
	authHelper    *helpers.AuthHelper
	userHelper    *helpers.UserHelper
	logger        *zap.Logger
	jwtSecret     string
	jwtExpiry     time.Duration
	refreshExpiry time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	oauthRepo repository.OAuthProviderRepository,
	roleRepo repository.RoleRepository,
	sessionRepo repository.SessionRepository,
	tenantRepo repository.TenantRepository,
	tokenRepo repository.TokenRepository,
	searchIndex *typesenseclient.SearchIndexService,
	rateLimiter *middlewares.RateLimiterMiddleware,
	tokenHelper *helpers.TokenHelper,
	authHelper *helpers.AuthHelper,
	userHelper *helpers.UserHelper,
	logger *zap.Logger,
	jwtSecret string,
	jwtExpiryHours int,
	refreshExpiryHours int,
) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		oauthRepo:     oauthRepo,
		roleRepo:      roleRepo,
		sessionRepo:   sessionRepo,
		tenantRepo:    tenantRepo,
		tokenRepo:     tokenRepo,
		searchIndex:   searchIndex,
		rateLimiter:   rateLimiter,
		tokenHelper:   tokenHelper,
		authHelper:    authHelper,
		userHelper:    userHelper,
		logger:        logger,
		jwtSecret:     jwtSecret,
		jwtExpiry:     time.Duration(jwtExpiryHours) * time.Hour,
		refreshExpiry: time.Duration(refreshExpiryHours) * time.Hour,
	}
}
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest, ip, platform string) (*dto.LoginResponse, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	blockKey := fmt.Sprintf("block:login:%s", ip)
	if blocked := s.IsBlockedAttempt(subCtx, blockKey); blocked {
		return nil, errors.New("too many attempts. please try again later")
	}

	user, err := s.userRepo.FindByEmail(subCtx, req.Email)

	if err != nil {
		attemptsKey := fmt.Sprintf("delay:%s:%s", "login", ip)
		attempts := s.ShouldBlockCredential(subCtx, attemptsKey)
		_ = s.SetExpireAttemptCredential(subCtx, attemptsKey, time.Hour)

		if attempts >= 3 {
			_ = s.SetBlockedAttemptCredential(subCtx, blockKey, "1", 30*time.Minute)
		}
		remaining := 3 - attempts

		return nil, fmt.Errorf("invalid credentials. %d attempts remaining", remaining)
	}

	validatePassword := utils.VerifyPassword(user.PasswordHash, req.Password)
	if !validatePassword {
		attemptsKey := fmt.Sprintf("delay:%s:%s", "login", ip)
		attempts := s.ShouldBlockCredential(subCtx, attemptsKey)
		_ = s.SetExpireAttemptCredential(subCtx, attemptsKey, time.Hour)

		if attempts >= 3 {
			_ = s.SetBlockedAttemptCredential(subCtx, blockKey, "1", 30*time.Minute)
		}
		remaining := 3 - attempts

		return nil, fmt.Errorf("invalid credentials. %d attempts remaining", remaining)
	}

	if !user.CanLogin() {
		s.logger.Warn("Login failed: user inactive",
			zap.String("email", req.Email),
			zap.Any("user_id", user.ID),
		)
		return nil, dto.ErrUserInactive
	}
	if !user.EmailVerifiedAt.Valid || user.EmailVerifiedAt.Time.IsZero() {
		tokenVerify, expiredAt := utils.GenerateEmailToken()

		tokenPayload := &entity.Token{
			ID:        uuid.New(),
			UserID:    user.ID,
			Token:     tokenVerify,
			Type:      string(dto.EmailVerification),
			IsUsed:    false,
			ExpiresAt: expiredAt,
			CreatedAt: time.Now(),
		}

		_, err = s.tokenRepo.Create(ctx, tokenPayload)
		if err != nil {
			s.logger.Error("Failed to create token",
				zap.Any("token_id", tokenPayload.ID),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to create token: %w", err)
		}
		metadata := &dto.SendMailMetaData{
			Token:     tokenVerify,
			Type:      dto.EmailVerification,
			To:        user.Email,
			User:      &entity.User{ID: user.ID, Email: user.Email, FullName: user.FullName, Phone: user.Phone},
			Password:  req.Password,
			ExpiredAt: expiredAt,
		}

		if err = s.authHelper.SendEmail(metadata); err != nil {
			s.logger.Error("Failed to send verification email",
				zap.String("email", user.Email),
				zap.Error(err),
			)
		}

		return &dto.LoginResponse{
			AccessToken:  "",
			RefreshToken: "",
			TwoFaToken:   "",
			TokenType:    "",
			ExpiresIn:    0,
			Status:       "require_email_verification",
			User:         nil,
		}, nil
	}

	if user.TwoFaSecret.Valid && user.TwoFaSecret.String != "" {
		twoFaToken := uuid.New().String()
		if err = s.userHelper.Set2FaStatus(subCtx, twoFaToken, "pending_2fa", user.ID.String()); err != nil {
			s.logger.Error("Failed to save 2FA token",
				zap.Any("user_id", user.ID),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to save 2FA token: %w", err)
		}

		return &dto.LoginResponse{
			AccessToken:  "",
			RefreshToken: "",
			TwoFaToken:   twoFaToken,
			TokenType:    "",
			ExpiresIn:    0,
			Status:       "two_fa_required",
			User:         nil,
		}, nil
	}

	userTenant, err := s.userRepo.GetUserTenantWithDetailsByUserID(subCtx, user.ID)
	if err != nil {
		s.logger.Error("Failed to get user tenant",
			zap.Any("user_id", user.ID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user tenant: %w", err)
	}
	if userTenant.UserTenant.ID == uuid.Nil {
		return nil, fmt.Errorf("user tenant relationship not found for user %s", user.ID)
	}
	if userTenant.Tenant.ID == uuid.Nil {
		return nil, fmt.Errorf("tenant not found for user %s", user.ID)
	}
	if userTenant.Role.ID == uuid.Nil {
		return nil, fmt.Errorf("role not found for user %s", user.ID)
	}

	roleNames := s.GetUserRoleNames(userTenant)

	var accTokenExp time.Duration
	var refreshTokenExp time.Duration
	if platform == "mobile" {
		accTokenExp = time.Duration(168) * time.Hour
		refreshTokenExp = time.Duration(336) * time.Hour
	} else {
		accTokenExp = s.jwtExpiry
		refreshTokenExp = s.refreshExpiry
	}

	tokenPayload := &entity.RedisSessionData{
		UserID:       user.ID,
		Email:        user.Email,
		TenantID:     userTenant.Tenant.ID,
		UserTenantID: userTenant.UserTenant.ID,
		RoleID:       userTenant.Role.ID,
		RoleName:     roleNames,
		SessionID:    uuid.New().String(),
		Metadata:     userTenant.Metadata,
		IssuedAt:     time.Now().Unix(),
		LastAccessed: time.Now().Unix(),
	}

	accessToken, err := utils.GenerateJWT(s.jwtSecret, tokenPayload, accTokenExp)
	if err != nil {
		s.logger.Error("Failed to generate access token",
			zap.Any("user_id", user.ID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateJWT(s.jwtSecret, tokenPayload, refreshTokenExp)
	if err != nil {
		s.logger.Error("Failed to generate refresh token",
			zap.Any("user_id", user.ID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		s.logger.Warn("Failed to update last login",
			zap.Any("user_id", user.ID),
			zap.Error(err),
		)
	}

	if err := s.tokenHelper.SetSessionToken(ctx, tokenPayload, s.jwtExpiry); err != nil {
		s.logger.Error("Failed to save session",
			zap.Any("user_id", user.ID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to save session: %w", err)
	}
	s.tokenHelper.DeleteAllExceptCurrent(ctx, tokenPayload.SessionID)

	s.logger.Info("User logged in successfully",
		zap.Any("user_id", user.ID),
		zap.String("email", user.Email),
		zap.String("session_id", tokenPayload.SessionID),
	)

	return &dto.LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TwoFaToken:       "",
		TokenType:        "Bearer",
		ExpiresIn:        int64(accTokenExp.Seconds()),
		ExpiresRefreshIn: int64(refreshTokenExp.Seconds()),
		Status:           "accepted",
		User: &entity.UserResponse{
			ID:              user.ID,
			Email:           user.Email,
			FullName:        user.FullName,
			Phone:           user.Phone.String,
			AvatarURL:       user.AvatarURL.String,
			TwoFaSecret:     user.TwoFaSecret.String,
			EmailVerifiedAt: user.EmailVerifiedAt,
			LastLoginAt:     user.LastLoginAt,
			CreatedAt:       user.CreatedAt,
			UpdatedAt:       user.UpdatedAt,
		},
	}, nil
}
func (s *AuthService) Register(ctx context.Context, req *dto.RegisterUserRequest) (*entity.User, error) {
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		s.logger.Warn("Registration failed: email already registered",
			zap.String("email", req.Email),
		)
		return nil, dto.ErrEmailAlreadyExists
	}
	exists, err := s.userRepo.ExistsByPhone(ctx, req.Phone)
	if err == nil && exists {
		s.logger.Warn("Registration failed: phone already registered",
			zap.String("phone", req.Phone),
		)
		return nil, dto.ErrPhoneAlreadyExists
	}
	if !utils.IsStrongPassword(req.Password) {
		s.logger.Warn("Registration failed: weak password",
			zap.String("email", req.Email),
		)
		return nil, dto.ErrWeakPassword
	}

	hashPassword, err := utils.GeneratePasswordHash(req.Password)
	if err != nil {
		s.logger.Error("Failed to hash password",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userUUID := uuid.New()
	tenantUUID := uuid.New()

	role, err := s.roleRepo.GetByName(ctx, "tenant_owner")
	if err != nil {
		s.logger.Error("Failed to get role",
			zap.String("role_name", entity.RoleTenantOwner),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	tenant := &entity.Tenant{
		ID:                 tenantUUID,
		Name:               req.FullName,
		Slug:               utils.GenerateSlugUnicodeV2(req.FullName),
		MaxDivisions:       1,
		MaxAgents:          1,
		MaxQuickReplies:    5,
		MaxWahaWhatsApp:    0,
		MaxMetaWhatsApp:    0,
		MaxMetaMessenger:   1,
		MaxInstagram:       1,
		MaxTelegram:        1,
		MaxWebChat:         1,
		MaxLinkChat:        1,
		AiCredits:          0,
		SubscriptionStatus: entity.StatusActive,
		SubscriptionPlan:   entity.PlanFree,
		IsActive:           true,
		TrialEndsAt:        entity.NewNullTimeFromNow(0, 0, 7),
		Settings:           &entity.TenantSettingConfig{},
		CreatedAt:          time.Now(),
	}
	user := &entity.User{
		ID:              userUUID,
		Email:           req.Email,
		FullName:        req.FullName,
		Phone:           entity.NewNullString(req.Phone),
		PasswordHash:    hashPassword,
		Status:          entity.UserStatusActive,
		EmailVerifiedAt: entity.NullTime{},
		LastLoginAt:     entity.NewNullTime(time.Now()),
		CreatedAt:       time.Now(),
	}

	// Create tenant, user and the owner membership atomically so a user always
	// has a valid tenant membership (the identity chain login depends on).
	if err = s.userRepo.WithTransaction(ctx, func(tx pgx.Tx) error {
		if err := s.tenantRepo.CreateTx(ctx, tx, tenant); err != nil {
			return fmt.Errorf("failed to create tenant: %w", err)
		}
		if err := s.userRepo.CreateTx(ctx, tx, user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		membership := &entity.UserTenant{
			UserID:   user.ID,
			TenantID: tenant.ID,
			RoleID:   role.ID,
			IsActive: true,
		}
		if err := s.userRepo.CreateUserTenantTx(ctx, tx, membership); err != nil {
			return fmt.Errorf("failed to create user membership: %w", err)
		}
		return nil
	}); err != nil {
		s.logger.Error("Failed to register user",
			zap.Any("user_id", user.ID),
			zap.Any("tenant_id", tenant.ID),
			zap.Error(err),
		)
		return nil, err
	}

	tokenVerify, expiredAt := utils.GenerateEmailToken()

	tokenPayload := &entity.Token{
		ID:        uuid.New(),
		UserID:    userUUID,
		Token:     tokenVerify,
		Type:      string(dto.EmailVerification),
		IsUsed:    false,
		ExpiresAt: expiredAt,
		CreatedAt: time.Now(),
	}

	_, err = s.tokenRepo.Create(ctx, tokenPayload)
	if err != nil {
		s.logger.Error("Failed to create token",
			zap.Any("token_id", tokenPayload.ID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create token: %w", err)
	}
	metadata := &dto.SendMailMetaData{
		Token:     tokenVerify,
		Type:      dto.EmailVerification,
		To:        req.Email,
		User:      &entity.User{ID: userUUID, Email: req.Email, FullName: req.FullName, Phone: entity.NewNullString(req.Phone)},
		Password:  req.Password,
		ExpiredAt: expiredAt,
	}

	if err := s.authHelper.SendEmail(metadata); err != nil {
		s.logger.Error("Failed to send verification email",
			zap.String("email", req.Email),
			zap.Error(err),
		)
	}
	if s.searchIndex != nil {
		if err := s.searchIndex.UpsertTenant(ctx, tenant); err != nil {
			s.logger.Warn("Failed to sync tenant to typesense", zap.Error(err))
		}
		if err := s.searchIndex.UpsertUser(ctx, user); err != nil {
			s.logger.Warn("Failed to sync user to typesense", zap.Error(err))
		}
	}
	return user, nil
}

func (s *AuthService) OAuthLoginOrRegister(ctx context.Context, gothUser *goth.User, platform, ip string) (*dto.LoginResponse, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	providerName := strings.ToLower(strings.TrimSpace(gothUser.Provider))
	providerID := strings.TrimSpace(gothUser.UserID)
	if providerName == "" || providerID == "" {
		return nil, fmt.Errorf("invalid oauth profile")
	}

	if oauthLink, err := s.oauthRepo.FindByProviderAndID(subCtx, providerName, providerID); err == nil && oauthLink != nil {
		user, err := s.userRepo.FindByID(subCtx, oauthLink.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to load oauth user: %w", err)
		}
		return s.buildOAuthLoginResponse(subCtx, user, platform, ip)
	}

	var user *entity.User
	if gothUser.Email != "" {
		if existingUser, err := s.userRepo.FindByEmail(subCtx, gothUser.Email); err == nil && existingUser != nil {
			user = existingUser
			if !existingUser.EmailVerifiedAt.Valid {
				if err := s.userRepo.SetEmailVerified(subCtx, existingUser.ID, true); err != nil {
					s.logger.Warn("Failed to verify oauth user email", zap.Error(err))
				} else {
					now := time.Now()
					user.Status = entity.UserStatusActive
					user.EmailVerifiedAt = entity.NewNullTimePtr(&now)
				}
			}
		}
	}

	if user == nil {
		if strings.TrimSpace(gothUser.Email) == "" {
			return nil, fmt.Errorf("oauth email is required")
		}

		displayName := strings.TrimSpace(gothUser.Name)
		if displayName == "" {
			displayName = strings.TrimSpace(gothUser.NickName)
		}
		if displayName == "" {
			displayName = strings.TrimSpace(gothUser.Email)
		}
		if displayName == "" {
			displayName = providerName
		}

		role, err := s.roleRepo.GetByName(subCtx, "tenant_owner")
		if err != nil || role == nil {
			return nil, fmt.Errorf("failed to resolve oauth role: %w", err)
		}

		tenant := &entity.Tenant{
			ID:                 uuid.New(),
			Name:               displayName,
			Slug:               utils.GenerateSlugUnicodeV2(displayName),
			MaxDivisions:       1,
			MaxAgents:          1,
			MaxQuickReplies:    1,
			MaxWahaWhatsApp:    0,
			MaxMetaWhatsApp:    0,
			MaxMetaMessenger:   1,
			MaxInstagram:       1,
			MaxTelegram:        1,
			MaxWebChat:         1,
			MaxLinkChat:        1,
			AiCredits:          0,
			SubscriptionStatus: entity.StatusActive,
			SubscriptionPlan:   entity.PlanFree,
			IsActive:           true,
			TrialEndsAt:        entity.NewNullTimeFromNow(0, 0, 7),
			Settings:           &entity.TenantSettingConfig{},
			CreatedAt:          time.Now(),
		}
		if err := s.tenantRepo.Create(subCtx, tenant); err != nil {
			return nil, fmt.Errorf("failed to create oauth tenant: %w", err)
		}

		password, err := utils.GenerateRandomPassword()
		if err != nil {
			return nil, fmt.Errorf("failed to generate oauth password: %w", err)
		}
		hashPassword, err := utils.GeneratePasswordHash(password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash oauth password: %w", err)
		}
		now := time.Now()
		user = &entity.User{
			ID:              uuid.New(),
			Email:           gothUser.Email,
			FullName:        displayName,
			Phone:           entity.NullString{},
			AvatarURL:       entity.NewNullString(gothUser.AvatarURL),
			PasswordHash:    hashPassword,
			Status:          entity.UserStatusActive,
			EmailVerifiedAt: entity.NewNullTimePtr(&now),
			LastLoginAt:     entity.NewNullTimePtr(&now),
			CreatedAt:       now,
		}
		if err := s.userRepo.Create(subCtx, user); err != nil {
			return nil, fmt.Errorf("failed to create oauth user: %w", err)
		}
	}

	link := &entity.OAuthProvider{
		ID:           uuid.New(),
		UserID:       user.ID,
		ProviderName: providerName,
		ProviderID:   providerID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if existingLink, err := s.oauthRepo.FindByProviderAndID(subCtx, providerName, providerID); err == nil && existingLink == nil {
		if err := s.oauthRepo.Create(subCtx, link); err != nil {
			return nil, fmt.Errorf("failed to create oauth link: %w", err)
		}
	}

	if s.searchIndex != nil {
		if details, err := s.userRepo.GetUserTenantWithDetailsByUserID(subCtx, user.ID); err == nil && details != nil {
			if err := s.searchIndex.UpsertTenant(subCtx, &details.Tenant); err != nil {
				s.logger.Warn("Failed to sync oauth tenant to typesense", zap.Error(err))
			}
		}
		if err := s.searchIndex.UpsertUser(subCtx, user); err != nil {
			s.logger.Warn("Failed to sync oauth user to typesense", zap.Error(err))
		}
	}

	return s.buildOAuthLoginResponse(subCtx, user, platform, ip)
}

func (s *AuthService) buildOAuthLoginResponse(ctx context.Context, user *entity.User, platform, ip string) (*dto.LoginResponse, error) {
	userTenant, err := s.userRepo.GetUserTenantWithDetailsByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load oauth user tenant: %w", err)
	}

	roleNames := s.GetUserRoleNames(userTenant)

	var accTokenExp time.Duration
	var refreshTokenExp time.Duration
	if platform == "mobile" {
		accTokenExp = time.Duration(168) * time.Hour
		refreshTokenExp = time.Duration(336) * time.Hour
	} else {
		accTokenExp = s.jwtExpiry
		refreshTokenExp = s.refreshExpiry
	}

	tokenPayload := &entity.RedisSessionData{
		UserID:       user.ID,
		Email:        user.Email,
		TenantID:     userTenant.Tenant.ID,
		UserTenantID: userTenant.UserTenant.ID,
		RoleID:       userTenant.Role.ID,
		RoleName:     roleNames,
		SessionID:    uuid.New().String(),
		Metadata:     userTenant.Metadata,
		IssuedAt:     time.Now().Unix(),
		LastAccessed: time.Now().Unix(),
	}

	accessToken, err := utils.GenerateJWT(s.jwtSecret, tokenPayload, accTokenExp)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err := utils.GenerateJWT(s.jwtSecret, tokenPayload, refreshTokenExp)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.tokenHelper.SetSessionToken(ctx, tokenPayload, s.jwtExpiry); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}
	s.tokenHelper.DeleteAllExceptCurrent(ctx, tokenPayload.SessionID)

	return &dto.LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        int64(accTokenExp.Seconds()),
		ExpiresRefreshIn: int64(refreshTokenExp.Seconds()),
		Status:           "accepted",
		User:             entity.NewUserResponse(user, &userTenant.Tenant, &userTenant.UserTenant, &userTenant.Role),
	}, nil
}
func (s *AuthService) FindOrCreateUser(ctx context.Context, user *goth.User) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	existingUser, err := s.userRepo.FindByEmail(subCtx, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if existingUser != nil {

		return existingUser, nil
	}

	return nil, nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, tokenString string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	token, err := s.tokenRepo.FindByToken(subCtx, tokenString)
	if err != nil {
		return fmt.Errorf("failed to find token: %w", err)
	}
	if token.IsUsed {
		return errors.New("token already used")
	}
	if token.IsExpired() {
		return errors.New("token expired")
	}
	if token.Type != string(dto.EmailVerification) {
		return errors.New("invalid token type")
	}

	user, err := s.userRepo.FindByID(subCtx, token.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user.Status == entity.UserStatusActive {
		return errors.New("email already verified")
	}

	if err := s.userRepo.SetEmailVerified(subCtx, user.ID, true); err != nil {
		return fmt.Errorf("failed to set email verified: %w", err)
	}

	if err := s.tokenRepo.HardDeleteByToken(subCtx, token.Token); err != nil {
		return fmt.Errorf("failed to hard delete token: %w", err)
	}

	return nil
}
func (s *AuthService) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest, ip string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	blockKey := fmt.Sprintf("block:forgot:%s", ip)
	if blocked := s.IsBlockedAttempt(subCtx, blockKey); blocked {
		return errors.New("too many attempts. Please try again later")
	}

	user, err := s.userRepo.FindByEmail(subCtx, req.Email)
	if err != nil {
		attemptsKey := fmt.Sprintf("delay:%s:%s", "forgot", ip)
		attempts, errIncrement := s.userHelper.IncrementAndGet(subCtx, attemptsKey, time.Hour)
		if errIncrement != nil {
			return fmt.Errorf("failed to increment attempts: %w", errIncrement)
		}
		if attempts > 3 {
			errBlock := s.userHelper.SetBlockedAttemptCredential(subCtx, blockKey, 1, 30*time.Minute)
			if errBlock != nil {
				return fmt.Errorf("failed to set blocked attempt: %w", errBlock)
			}
			s.userHelper.ResetCounter(subCtx, attemptsKey)
			return errors.New("too many attempts. Please try again later")
		}

		remaining := 3 - attempts
		return fmt.Errorf("failed to get user: %w. %d attempts remaining", err, remaining)
	}
	if user.Status != entity.UserStatusActive {
		return errors.New("user is inactive")
	}

	tokenVerify, expiredAt := utils.GenerateEmailToken()

	tokenPayload := &entity.Token{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     tokenVerify,
		Type:      string(dto.ResetPassword),
		IsUsed:    false,
		ExpiresAt: expiredAt,
		CreatedAt: time.Now(),
	}
	newToken, err := s.tokenRepo.CreateOrGetExist(subCtx, tokenPayload)
	if err != nil {
		s.logger.Error("Failed to create or update token",
			zap.Any("token_id", tokenPayload.ID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create or update token: %w", err)
	}
	metadata := &dto.SendMailMetaData{
		Token:     newToken.Token,
		Type:      dto.ResetPassword,
		To:        req.Email,
		User:      &entity.User{ID: user.ID, Email: req.Email, FullName: user.FullName, Phone: user.Phone},
		ExpiredAt: newToken.ExpiresAt,
	}

	if err := s.authHelper.SendEmail(metadata); err != nil {
		s.logger.Error("Failed to send verification email",
			zap.String("email", req.Email),
			zap.Error(err),
		)
	}
	return nil
}
func (s *AuthService) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	token, err := s.tokenRepo.FindByToken(subCtx, req.Token)
	if err != nil {
		return fmt.Errorf("failed to find token: %w", err)
	}
	if token.IsUsed {
		return errors.New("token already used")
	}
	if token.IsExpired() {
		return errors.New("token expired")
	}
	if token.Type != string(dto.ResetPassword) {
		return errors.New("invalid token type")
	}

	user, err := s.userRepo.FindByID(subCtx, token.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	isMatchPass := utils.VerifyPassword(user.PasswordHash, req.NewPassword)
	if isMatchPass {
		return errors.New("new password cannot be the same as old password")
	}

	isWeakPass := utils.IsStrongPassword(req.NewPassword)
	if !isWeakPass {
		return errors.New("password is weak")
	}

	hashPass, err := utils.GeneratePasswordHash(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %w", err)
	}

	if err := s.userRepo.UpdatePassword(subCtx, token.UserID, hashPass); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	if err := s.tokenRepo.HardDeleteByToken(subCtx, token.Token); err != nil {
		return fmt.Errorf("failed to hard delete token: %w", err)
	}

	return nil
}
func (s *AuthService) VerifyTwoFactor(ctx context.Context, req *dto.VerifyTwoFactorRequest, ip, platform string) (*dto.LoginResponse, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	blockKey := fmt.Sprintf("block:verify2fa:%s", ip)
	if blocked := s.IsBlockedAttempt(subCtx, blockKey); blocked {
		return nil, errors.New("too many attempts. Please try again later")
	}

	userID, err := s.userHelper.Get2FaStatus(subCtx, req.Token, "pending_2fa")
	if err != nil {
		return nil, fmt.Errorf("failed to get 2fa status: %w", err)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	userTenant, err := s.userRepo.GetUserTenantWithDetailsByUserID(subCtx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tenant: %w", err)
	}
	if userTenant.UserTenant.ID == uuid.Nil {
		return nil, fmt.Errorf("user tenant relationship not found for user %s", userTenant.User.ID)
	}
	if userTenant.Tenant.ID == uuid.Nil {
		return nil, fmt.Errorf("tenant not found for user %s", userTenant.User.ID)
	}
	if userTenant.Role.ID == uuid.Nil {
		return nil, fmt.Errorf("role not found for user %s", userTenant.User.ID)
	}

	if !userTenant.User.IsActive() {
		return nil, errors.New("user is inactive")
	}

	if !userTenant.User.TwoFaSecret.Valid || userTenant.User.TwoFaSecret.String == "" {
		return nil, errors.New("two factor authentication not enabled")
	}

	valid, err := s.userHelper.Verify2FA(subCtx, userTenant.User.ID, req.OTP, userTenant.User.TwoFaSecret.String)
	if err != nil || !valid {
		attemptsKey := fmt.Sprintf("delay:%s:%s", "verify2fa", ip)
		attempts, errIncrement := s.userHelper.IncrementAndGet(subCtx, attemptsKey, time.Hour)
		if errIncrement != nil {
			return nil, fmt.Errorf("failed to increment attempts: %w", errIncrement)
		}
		if attempts > 3 {
			errBlock := s.userHelper.SetBlockedAttemptCredential(subCtx, blockKey, 1, 30*time.Minute)
			if errBlock != nil {
				return nil, fmt.Errorf("failed to set blocked attempt: %w", errBlock)
			}
			s.userHelper.ResetCounter(subCtx, attemptsKey)
			return nil, errors.New("too many attempts. Please try again later")
		}
		remaining := 3 - int(attempts)
		return nil, fmt.Errorf("failed to validate 2fa: %w. %d attempts remaining", err, remaining)
	}
	if err = s.userHelper.Clear2FaStatus(subCtx, req.Token, "pending_2fa"); err != nil {
		s.logger.Error("Failed to clear 2fa status",
			zap.String("token", req.Token),
			zap.Error(err),
		)
	}

	roleNames := s.GetUserRoleNames(userTenant)

	var accTokenExp time.Duration
	var refreshTokenExp time.Duration
	if platform == "mobile" {
		accTokenExp = time.Duration(168) * time.Hour
		refreshTokenExp = time.Duration(336) * time.Hour
	} else {
		accTokenExp = s.jwtExpiry
		refreshTokenExp = s.refreshExpiry
	}

	tokenPayload := &entity.RedisSessionData{
		UserID:       userTenant.User.ID,
		Email:        userTenant.User.Email,
		TenantID:     userTenant.Tenant.ID,
		UserTenantID: userTenant.UserTenant.ID,
		RoleID:       userTenant.Role.ID,
		RoleName:     roleNames,
		SessionID:    uuid.New().String(),
		IssuedAt:     time.Now().Unix(),
		LastAccessed: time.Now().Unix(),
	}

	accessToken, err := utils.GenerateJWT(s.jwtSecret, tokenPayload, accTokenExp)
	if err != nil {
		s.logger.Error("Failed to generate access token",
			zap.Any("user_id", tokenPayload.UserID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateJWT(s.jwtSecret, tokenPayload, refreshTokenExp)
	if err != nil {
		s.logger.Error("Failed to generate refresh token",
			zap.Any("user_id", tokenPayload.UserID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.userRepo.UpdateLastLogin(ctx, tokenPayload.UserID); err != nil {
		s.logger.Warn("Failed to update last login",
			zap.Any("user_id", tokenPayload.UserID),
			zap.Error(err),
		)
	}

	if err := s.tokenHelper.SetSessionToken(ctx, tokenPayload, s.jwtExpiry); err != nil {
		s.logger.Error("Failed to save session",
			zap.Any("user_id", tokenPayload.UserID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to save session: %w", err)
	}
	s.tokenHelper.DeleteAllExceptCurrent(ctx, tokenPayload.SessionID)

	s.logger.Info("User logged in successfully",
		zap.Any("user_id", tokenPayload.UserID),
		zap.String("email", tokenPayload.Email),
		zap.String("session_id", tokenPayload.SessionID),
	)

	return &dto.LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TwoFaToken:       "",
		TokenType:        "Bearer",
		ExpiresIn:        int64(accTokenExp.Seconds()),
		ExpiresRefreshIn: int64(refreshTokenExp.Seconds()),
		Status:           "accepted",
		User:             entity.NewUserResponse(&userTenant.User, &userTenant.Tenant, &userTenant.UserTenant, &userTenant.Role),
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*dto.JWTClaims, error) {
	token, err := utils.VerifyJWT(tokenString, s.jwtSecret)
	if err != nil {
		s.logger.Debug("Token validation failed", zap.Error(err))
		return nil, dto.ErrInvalidToken
	}
	if !token.Valid {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("unauthorized - invalid or expired token: %w", err)
		}
		return nil, dto.ErrInvalidToken
	}

	return nil, dto.ErrInvalidToken
}
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken, platform string) (*dto.LoginResponse, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	userUUID, err := uuid.Parse(claims.UserID)
	if err != nil {
		s.logger.Error("Invalid user ID format",
			zap.Any("user_id", claims.UserID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	userTenant, err := s.userRepo.GetUserTenantWithDetailsByUserID(subCtx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tenant: %w", err)
	}
	if userTenant.UserTenant.ID == uuid.Nil {
		return nil, fmt.Errorf("user tenant relationship not found for user %s", userTenant.User.ID)
	}
	if userTenant.Tenant.ID == uuid.Nil {
		return nil, fmt.Errorf("tenant not found for user %s", userTenant.User.ID)
	}
	if userTenant.Role.ID == uuid.Nil {
		return nil, fmt.Errorf("role not found for user %s", userTenant.User.ID)
	}

	if !userTenant.User.IsActive() {
		return nil, errors.New("user is inactive")
	}

	roleNames := s.GetUserRoleNames(userTenant)

	var accTokenExp time.Duration
	var refreshTokenExp time.Duration
	if platform == "mobile" {
		accTokenExp = time.Duration(168) * time.Hour
		refreshTokenExp = time.Duration(336) * time.Hour
	} else {
		accTokenExp = s.jwtExpiry
		refreshTokenExp = s.refreshExpiry
	}

	tokenPayload := &entity.RedisSessionData{
		UserID:       userTenant.User.ID,
		Email:        userTenant.User.Email,
		TenantID:     userTenant.Tenant.ID,
		UserTenantID: userTenant.UserTenant.ID,
		RoleID:       userTenant.Role.ID,
		RoleName:     roleNames,
		SessionID:    claims.SessionID,
		Metadata:     userTenant.Metadata,
		IssuedAt:     time.Now().Unix(),
		LastAccessed: time.Now().Unix(),
	}

	accessToken, err := utils.GenerateJWT(s.jwtSecret, tokenPayload, accTokenExp)
	if err != nil {
		s.logger.Error("Failed to generate access token",
			zap.Any("user_id", tokenPayload.UserID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err = utils.GenerateJWT(s.jwtSecret, tokenPayload, refreshTokenExp)
	if err != nil {
		s.logger.Error("Failed to generate refresh token",
			zap.Any("user_id", tokenPayload.UserID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.tokenHelper.SetSessionToken(subCtx, tokenPayload, accTokenExp); err != nil {
		return nil, fmt.Errorf("failed to set session token: %w", err)
	}

	s.tokenHelper.DeleteAllExceptCurrent(ctx, tokenPayload.SessionID)

	return &dto.LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        int64(accTokenExp.Seconds()),
		ExpiresRefreshIn: int64(refreshTokenExp.Seconds()),
		Status:           "accepted",
		User: &entity.UserResponse{
			ID:              userTenant.User.ID,
			Email:           userTenant.User.Email,
			FullName:        userTenant.User.FullName,
			Phone:           userTenant.User.Phone.String,
			AvatarURL:       userTenant.User.AvatarURL.String,
			TwoFaSecret:     userTenant.User.TwoFaSecret.String,
			EmailVerifiedAt: userTenant.User.EmailVerifiedAt,
			LastLoginAt:     userTenant.User.LastLoginAt,
			CreatedAt:       userTenant.User.CreatedAt,
			UpdatedAt:       userTenant.User.UpdatedAt,
		},
	}, nil
}
func (s *AuthService) CheckPermission(ctx context.Context, userID uuid.UUID, permission string) error {
	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return dto.ErrUserNotFound
	}

	if !user.CanLogin() {
		return dto.ErrUserInactive
	}

	role, err := s.roleRepo.FindByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	switch role.Slug {
	case entity.RoleSuperAdmin, entity.RoleTenantOwner:
		return nil
	case entity.RoleSupervisor:
		if permission == "" {
			return nil
		}
	case entity.RoleAgent:
		if permission == "conversation.read" || permission == "conversation.write" {
			return nil
		}
	}

	s.logger.Warn("Permission denied",
		zap.Any("user_id", userID),
		zap.String("permission", permission),
	)

	return dto.ErrPermissionDenied
}
func (s *AuthService) CheckPermissions(ctx context.Context, userID uuid.UUID, requiredPermissions []string) error {
	for _, permission := range requiredPermissions {
		if err := s.CheckPermission(ctx, userID, permission); err != nil {
			return err
		}
	}
	return nil
}

func (s *AuthService) SaveSession(ctx context.Context, accToken, refToken string, userID uuid.UUID) error {
	session := &entity.Session{
		ID:             uuid.Must(uuid.NewRandom()),
		UserID:         userID,
		AccessToken:    accToken,
		RefreshToken:   refToken,
		ExpiresAt:      time.Now().Add(s.jwtExpiry),
		LastActivityAt: utils.TimePtr(time.Now()),
		CreatedAt:      time.Now(),
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}
	return nil
}

// GetUserRoleNames returns the role identifiers (name, plus slug when it
// differs) used for level/name-based authorization. Permissions were removed
// in favour of role level/name, so no permission lists are produced.
func (s *AuthService) GetUserRoleNames(metaData *entity.UserTenantWithDetails) []string {
	return metaData.RoleNames()
}
func (s *AuthService) IsBlockedAttempt(ctx context.Context, key string) bool {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	blocked, err := s.userHelper.IsBlockedAttempt(subCtx, key)
	if err != nil {
		return false
	}
	return blocked
}
func (s *AuthService) ShouldBlockCredential(ctx context.Context, key string) int {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	return s.userHelper.ShouldBlockCredential(subCtx, key)
}
func (s *AuthService) SetExpireAttemptCredential(ctx context.Context, key string, expiration time.Duration) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	return s.userHelper.SetExpireAttemptCredential(subCtx, key, expiration)
}
func (s *AuthService) SetBlockedAttemptCredential(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 3*time.Second)
	defer cancel()

	return s.userHelper.SetBlockedAttemptCredential(subCtx, key, val, expiration)
}
