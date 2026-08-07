package handlers

import (
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type ContactHandler struct {
	ctxinject *middlewares.ContextMiddleware
	service   *services.ContactService
	logger    *zap.Logger
}

func NewContactHandler(ctxinject *middlewares.ContextMiddleware, service *services.ContactService, logger *zap.Logger) *ContactHandler {
	return &ContactHandler{ctxinject: ctxinject, service: service, logger: logger}
}

func (h *ContactHandler) tenantID(c fiber.Ctx) (string, bool) {
	tid, ok := c.Locals("tenant_id").(string)
	if !ok || tid == "" {
		return "", false
	}
	return tid, true
}

// List returns a paginated page of contacts with pagination meta.
func (h *ContactHandler) List(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	p := helpers.ParsePageParams(c)
	contacts, total, err := h.service.List(ctx, tenantID, c.Query("channel_id"), c.Query("search"), p.Limit, p.Offset)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.RespondWithMeta(c, fiber.StatusOK, "Contacts retrieved successfully", contacts, helpers.NewPageMeta(p, total))
}

func (h *ContactHandler) Block(c fiber.Ctx) error {
	return h.setBlocked(c, true)
}
func (h *ContactHandler) Unblock(c fiber.Ctx) error {
	return h.setBlocked(c, false)
}
func (h *ContactHandler) setBlocked(c fiber.Ctx, blocked bool) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if err := h.service.SetBlocked(ctx, tenantID, c.Params("id"), blocked); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	msg := "Contact unblocked"
	if blocked {
		msg = "Contact blocked"
	}
	return helpers.Respond(c, fiber.StatusOK, msg, nil)
}

func (h *ContactHandler) Delete(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if err := h.service.Delete(ctx, tenantID, c.Params("id")); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Contact deleted successfully", nil)
}
