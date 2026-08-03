package handlers

import (
	"github/socialforge/internal/dto"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type WorkingHoursHandler struct {
	ctxinject *middlewares.ContextMiddleware
	service   *services.WorkingHoursService
	logger    *zap.Logger
}

func NewWorkingHoursHandler(ctxinject *middlewares.ContextMiddleware, service *services.WorkingHoursService, logger *zap.Logger) *WorkingHoursHandler {
	return &WorkingHoursHandler{ctxinject: ctxinject, service: service, logger: logger}
}

func (h *WorkingHoursHandler) tenantID(c fiber.Ctx) (string, bool) {
	tid, ok := c.Locals("tenant_id").(string)
	if !ok || tid == "" {
		return "", false
	}
	return tid, true
}

func (h *WorkingHoursHandler) canManage(c fiber.Ctx) bool {
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

func (h *WorkingHoursHandler) List(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	items, err := h.service.List(ctx, tenantID, c.Params("userId"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Working hours retrieved", items)
}

func (h *WorkingHoursHandler) Set(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.canManage(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only owner or supervisor can set working hours", nil)
	}
	var req dto.SetWorkingHoursRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}
	if err := h.service.Replace(ctx, tenantID, c.Params("userId"), &req); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Working hours saved", nil)
}
