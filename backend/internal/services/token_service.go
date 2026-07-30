package services

import (
	"context"
	"errors"
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/repository"
	"github/socialforge/internal/utils"

	"time"

	"go.uber.org/zap"
)

type TokenService struct {
	tokenRepo   repository.TokenRepository
	tokenHelper *helpers.TokenHelper
	logger      *zap.Logger
}

func NewTokenService(
	tokenRepo repository.TokenRepository,
	tokenHelper *helpers.TokenHelper,
	logger *zap.Logger,
) *TokenService {
	return &TokenService{
		tokenRepo:   tokenRepo,
		tokenHelper: tokenHelper,
		logger:      logger,
	}
}
func (s *TokenService) StoreCSRFToken(ctx context.Context) (string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 10*time.Second)
	defer cancel()

	token := utils.GenerateRandomToken()
	if token == "" {
		return "", errors.New("token is empty")
	}
	err := s.tokenHelper.StoreSessionCSRF(subCtx, token)
	if err != nil {
		return "", err
	}
	return token, nil
}
