package handlers

import (
	"errors"
	"github/socialforge/internal/dto"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/infra/repository"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type QuickReplyHandler struct {
	ctxinject *middlewares.ContextMiddleware
	service   *services.QuickReplyService
	logger    *zap.Logger
}

func NewQuickReplyHandler(ctxinject *middlewares.ContextMiddleware, service *services.QuickReplyService, logger *zap.Logger) *QuickReplyHandler {
	return &QuickReplyHandler{ctxinject: ctxinject, service: service, logger: logger}
}

func (h *QuickReplyHandler) tenantID(c fiber.Ctx) (string, bool) {
	tid, ok := c.Locals("tenant_id").(string)
	if !ok || tid == "" {
		return "", false
	}
	return tid, true
}

func (h *QuickReplyHandler) canManage(c fiber.Ctx) bool {
	names, ok := c.Locals("role_name").([]string)
	if !ok {
		return false
	}
	for _, n := range names {
		if n == entity.RoleTenantOwner || n == entity.RoleSupervisor {
			return true
		}
	}
	return false
}

func qrStatus(err error) int {
	if errors.Is(err, repository.ErrShortcutTaken) {
		return fiber.StatusConflict
	}
	return fiber.StatusInternalServerError
}

// List is available to any tenant member (agents use them via "/" typeahead).
func (h *QuickReplyHandler) List(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	items, err := h.service.List(ctx, tenantID, c.Query("search"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Quick replies retrieved successfully", items)
}

func (h *QuickReplyHandler) Create(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.canManage(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only owner or supervisor can manage quick replies", nil)
	}
	var req dto.CreateQuickReplyRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}
	qr, err := h.service.Create(ctx, tenantID, &req)
	if err != nil {
		return helpers.Respond(c, qrStatus(err), err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusCreated, "Quick reply created successfully", qr)
}

func (h *QuickReplyHandler) Update(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.canManage(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only owner or supervisor can manage quick replies", nil)
	}
	var req dto.UpdateQuickReplyRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}
	qr, err := h.service.Update(ctx, tenantID, c.Params("id"), &req)
	if err != nil {
		return helpers.Respond(c, qrStatus(err), err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Quick reply updated successfully", qr)
}

func (h *QuickReplyHandler) Delete(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.canManage(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only owner or supervisor can manage quick replies", nil)
	}
	if err := h.service.Delete(ctx, tenantID, c.Params("id")); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Quick reply deleted successfully", nil)
}
