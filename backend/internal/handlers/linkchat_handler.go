package handlers

import (
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type LinkchatHandler struct {
	ctxinject *middlewares.ContextMiddleware
	service   *services.LinkchatService
	logger    *zap.Logger
}

func NewLinkchatHandler(ctxinject *middlewares.ContextMiddleware, service *services.LinkchatService, logger *zap.Logger) *LinkchatHandler {
	return &LinkchatHandler{ctxinject: ctxinject, service: service, logger: logger}
}

// Resolve is the public landing endpoint for a shared division link.
func (h *LinkchatHandler) Resolve(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	info, err := h.service.Resolve(ctx, c.Params("token"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusNotFound, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Link resolved", info)
}
