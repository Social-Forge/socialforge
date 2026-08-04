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

type AIAgentHandler struct {
	ctxinject *middlewares.ContextMiddleware
	service   *services.AIAgentService
	logger    *zap.Logger
}

func NewAIAgentHandler(ctxinject *middlewares.ContextMiddleware, service *services.AIAgentService, logger *zap.Logger) *AIAgentHandler {
	return &AIAgentHandler{ctxinject: ctxinject, service: service, logger: logger}
}

func (h *AIAgentHandler) tenantID(c fiber.Ctx) (string, bool) {
	tid, ok := c.Locals("tenant_id").(string)
	if !ok || tid == "" {
		return "", false
	}
	return tid, true
}

func (h *AIAgentHandler) isOwner(c fiber.Ctx) bool {
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

func (h *AIAgentHandler) List(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	items, err := h.service.List(ctx, tenantID)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "AI agents retrieved successfully", items)
}

func (h *AIAgentHandler) Get(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	agent, err := h.service.Get(ctx, tenantID, c.Params("id"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusNotFound, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "AI agent retrieved successfully", agent)
}

func (h *AIAgentHandler) Create(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.isOwner(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can manage AI agents", nil)
	}
	var req dto.CreateAIAgentRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}
	agent, err := h.service.Create(ctx, tenantID, &req)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusCreated, "AI agent created successfully", agent)
}

func (h *AIAgentHandler) Update(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.isOwner(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can manage AI agents", nil)
	}
	var req dto.UpdateAIAgentRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}
	agent, err := h.service.Update(ctx, tenantID, c.Params("id"), &req)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "AI agent updated successfully", agent)
}

func (h *AIAgentHandler) Delete(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.isOwner(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can manage AI agents", nil)
	}
	if err := h.service.Delete(ctx, tenantID, c.Params("id")); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "AI agent deleted successfully", nil)
}
