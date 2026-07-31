package handlers

import (
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/infra/repository"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type SuperadminHandler struct {
	ctxinject *middlewares.ContextMiddleware
	service   *services.SuperadminService
	logger    *zap.Logger
}

func NewSuperadminHandler(ctxinject *middlewares.ContextMiddleware, service *services.SuperadminService, logger *zap.Logger) *SuperadminHandler {
	return &SuperadminHandler{ctxinject: ctxinject, service: service, logger: logger}
}

func (h *SuperadminHandler) ListTenants(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	opts := repository.NewListOptions()
	if s := c.Query("search"); s != "" {
		opts.Filter.Search = s
	}

	tenants, total, err := h.service.ListTenants(ctx, opts)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Tenants retrieved successfully", fiber.Map{
		"tenants": tenants,
		"total":   total,
	})
}
