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

type ChannelHandler struct {
	ctxinject *middlewares.ContextMiddleware
	service   *services.ChannelService
	logger    *zap.Logger
}

func NewChannelHandler(ctxinject *middlewares.ContextMiddleware, service *services.ChannelService, logger *zap.Logger) *ChannelHandler {
	return &ChannelHandler{ctxinject: ctxinject, service: service, logger: logger}
}

func (h *ChannelHandler) tenantID(c fiber.Ctx) (string, bool) {
	tid, ok := c.Locals("tenant_id").(string)
	if !ok || tid == "" {
		return "", false
	}
	return tid, true
}

func (h *ChannelHandler) isOwner(c fiber.Ctx) bool {
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

func (h *ChannelHandler) List(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	channels, err := h.service.List(ctx, tenantID, c.Query("type"), c.Query("division_id"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Channels retrieved successfully", channels)
}

func (h *ChannelHandler) Get(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	channel, err := h.service.Get(ctx, tenantID, c.Params("id"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusNotFound, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Channel retrieved successfully", channel)
}

func (h *ChannelHandler) Create(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.isOwner(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can add channels", nil)
	}

	var req dto.CreateChannelRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}

	channel, err := h.service.Create(ctx, tenantID, &req)
	if err != nil {
		h.logger.Error("Failed to create channel", zap.String("tenant_id", tenantID), zap.Error(err))
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusCreated, "Channel created successfully", channel)
}

func (h *ChannelHandler) Update(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.isOwner(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can update channels", nil)
	}

	var req dto.UpdateChannelRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}

	channel, err := h.service.Update(ctx, tenantID, c.Params("id"), &req)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Channel updated successfully", channel)
}

func (h *ChannelHandler) GetAutoResponse(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	ar, err := h.service.GetAutoResponse(ctx, tenantID, c.Params("id"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Auto-response retrieved", ar)
}

func (h *ChannelHandler) SetAutoResponse(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.isOwner(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can set auto-response", nil)
	}
	var req dto.SetAutoResponseRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}
	ar, err := h.service.SetAutoResponse(ctx, tenantID, c.Params("id"), &req)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Auto-response saved", ar)
}

func (h *ChannelHandler) Connect(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.isOwner(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can connect channels", nil)
	}
	info, err := h.service.Connect(ctx, tenantID, c.Params("id"))
	if err != nil {
		h.logger.Error("Failed to connect channel", zap.String("tenant_id", tenantID), zap.Error(err))
		return helpers.Respond(c, fiber.StatusBadGateway, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Channel connect initiated", info)
}

func (h *ChannelHandler) Delete(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.isOwner(c) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can delete channels", nil)
	}
	if err := h.service.Delete(ctx, tenantID, c.Params("id")); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Channel deleted successfully", nil)
}
