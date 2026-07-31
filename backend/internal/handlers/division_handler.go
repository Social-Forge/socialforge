package handlers

import (
	"github/socialforge/internal/dto"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/infra/repository"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type DivisionHandler struct {
	ctxinject *middlewares.ContextMiddleware
	service   *services.DivisionService
	logger    *zap.Logger
}

func NewDivisionHandler(ctxinject *middlewares.ContextMiddleware, service *services.DivisionService, logger *zap.Logger) *DivisionHandler {
	return &DivisionHandler{
		ctxinject: ctxinject,
		service:   service,
		logger:    logger,
	}
}

// tenantID reads the authenticated tenant id injected by the auth middleware.
func (h *DivisionHandler) tenantID(c fiber.Ctx) (string, bool) {
	tid, ok := c.Locals("tenant_id").(string)
	if !ok || tid == "" {
		return "", false
	}
	return tid, true
}

// hasAnyRole reports whether the caller holds any of the given role names.
func (h *DivisionHandler) hasAnyRole(c fiber.Ctx, roles ...string) bool {
	names, ok := c.Locals("role_name").([]string)
	if !ok {
		return false
	}
	for _, have := range names {
		for _, want := range roles {
			if have == want {
				return true
			}
		}
	}
	return false
}

func (h *DivisionHandler) Create(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.hasAnyRole(c, entity.RoleTenantOwner) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can create divisions", nil)
	}

	var req dto.CreateDivisionRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}

	division, err := h.service.Create(ctx, tenantID, &req)
	if err != nil {
		h.logger.Error("Failed to create division", zap.String("tenant_id", tenantID), zap.Error(err))
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusCreated, "Division created successfully", division)
}

func (h *DivisionHandler) List(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}

	opts := repository.NewListOptions()
	if s := c.Query("search"); s != "" {
		opts.Filter.Search = s
	}

	divisions, total, err := h.service.List(ctx, tenantID, opts)
	if err != nil {
		h.logger.Error("Failed to list divisions", zap.String("tenant_id", tenantID), zap.Error(err))
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Divisions retrieved successfully", fiber.Map{
		"divisions": divisions,
		"total":     total,
	})
}

func (h *DivisionHandler) Get(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}

	division, err := h.service.Get(ctx, tenantID, c.Params("id"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusNotFound, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Division retrieved successfully", division)
}

func (h *DivisionHandler) Update(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.hasAnyRole(c, entity.RoleTenantOwner) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can update divisions", nil)
	}

	var req dto.UpdateDivisionRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}

	division, err := h.service.Update(ctx, tenantID, c.Params("id"), &req)
	if err != nil {
		h.logger.Error("Failed to update division", zap.String("tenant_id", tenantID), zap.Error(err))
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Division updated successfully", division)
}

func (h *DivisionHandler) Delete(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.hasAnyRole(c, entity.RoleTenantOwner) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can delete divisions", nil)
	}

	if err := h.service.Delete(ctx, tenantID, c.Params("id")); err != nil {
		h.logger.Error("Failed to delete division", zap.String("tenant_id", tenantID), zap.Error(err))
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Division deleted successfully", nil)
}

func (h *DivisionHandler) ListMembers(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}

	members, err := h.service.ListMembers(ctx, tenantID, c.Params("id"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Division members retrieved successfully", members)
}

func (h *DivisionHandler) AddMember(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.hasAnyRole(c, entity.RoleTenantOwner, entity.RoleSupervisor) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only owner or supervisor can manage members", nil)
	}

	var req dto.AddDivisionMemberRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}

	member, err := h.service.AddMember(ctx, tenantID, c.Params("id"), &req)
	if err != nil {
		h.logger.Error("Failed to add division member", zap.String("tenant_id", tenantID), zap.Error(err))
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusCreated, "Member added to division successfully", member)
}

func (h *DivisionHandler) RemoveMember(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.hasAnyRole(c, entity.RoleTenantOwner, entity.RoleSupervisor) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only owner or supervisor can manage members", nil)
	}

	if err := h.service.RemoveMember(ctx, tenantID, c.Params("id"), c.Params("userTenantID")); err != nil {
		h.logger.Error("Failed to remove division member", zap.String("tenant_id", tenantID), zap.Error(err))
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Member removed from division successfully", nil)
}
