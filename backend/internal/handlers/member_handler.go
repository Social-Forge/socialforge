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

type MemberHandler struct {
	ctxinject *middlewares.ContextMiddleware
	service   *services.MemberService
	logger    *zap.Logger
}

func NewMemberHandler(ctxinject *middlewares.ContextMiddleware, service *services.MemberService, logger *zap.Logger) *MemberHandler {
	return &MemberHandler{ctxinject: ctxinject, service: service, logger: logger}
}

func (h *MemberHandler) tenantID(c fiber.Ctx) (string, bool) {
	tid, ok := c.Locals("tenant_id").(string)
	if !ok || tid == "" {
		return "", false
	}
	return tid, true
}

func (h *MemberHandler) isOwner(c fiber.Ctx) bool {
	names, ok := c.Locals("role_name").([]string)
	if !ok {
		return false
	}
	for _, n := range names {
		if n == entity.RoleTenantOwner {
			return true
		}
	}
	return false
}

func (h *MemberHandler) List(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.isOwner(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can manage members", nil)
	}

	members, err := h.service.List(ctx, tenantID, c.Query("role"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Members retrieved successfully", members)
}

func (h *MemberHandler) Create(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.isOwner(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can add members", nil)
	}

	var req dto.CreateMemberRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}

	member, err := h.service.Create(ctx, tenantID, &req)
	if err != nil {
		h.logger.Error("Failed to create member", zap.String("tenant_id", tenantID), zap.Error(err))
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusCreated, "Member created successfully", member)
}

func (h *MemberHandler) Update(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.isOwner(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can update members", nil)
	}

	var req dto.UpdateMemberRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}

	if err := h.service.Update(ctx, tenantID, c.Params("userTenantID"), &req); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Member updated successfully", nil)
}

func (h *MemberHandler) Delete(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.isOwner(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can remove members", nil)
	}

	if err := h.service.Delete(ctx, tenantID, c.Params("userTenantID")); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Member removed successfully", nil)
}
