package handlers

import (
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type TokenHandler struct {
	ctxinject    *middlewares.ContextMiddleware
	tokenService *services.TokenService
	logger       *zap.Logger
}

func NewTokenHandler(
	ctxinject *middlewares.ContextMiddleware,
	tokenService *services.TokenService,
	logger *zap.Logger,
) *TokenHandler {
	return &TokenHandler{
		ctxinject:    ctxinject,
		tokenService: tokenService,
		logger:       logger,
	}
}
func (h *TokenHandler) GetCSRFToken(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	token, err := h.tokenService.StoreCSRFToken(ctx)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "CSRF token generated", fiber.Map{"csrf_token": token})
}
