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

type LabelHandler struct {
	ctxinject *middlewares.ContextMiddleware
	service   *services.LabelService
	logger    *zap.Logger
}

func NewLabelHandler(ctxinject *middlewares.ContextMiddleware, service *services.LabelService, logger *zap.Logger) *LabelHandler {
	return &LabelHandler{ctxinject: ctxinject, service: service, logger: logger}
}

func (h *LabelHandler) tenantID(c fiber.Ctx) (string, bool) {
	tid, ok := c.Locals("tenant_id").(string)
	if !ok || tid == "" {
		return "", false
	}
	return tid, true
}

func (h *LabelHandler) hasAnyRole(c fiber.Ctx, roles ...string) bool {
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

func labelStatus(err error) int {
	if errors.Is(err, repository.ErrLabelNameTaken) {
		return fiber.StatusConflict
	}
	return fiber.StatusInternalServerError
}

func (h *LabelHandler) List(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	labels, err := h.service.List(ctx, tenantID)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Labels retrieved successfully", labels)
}

func (h *LabelHandler) Create(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	var req dto.CreateLabelRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}
	label, err := h.service.Create(ctx, tenantID, &req)
	if err != nil {
		return helpers.Respond(c, labelStatus(err), err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusCreated, "Label created successfully", label)
}

func (h *LabelHandler) Update(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	var req dto.UpdateLabelRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}
	label, err := h.service.Update(ctx, tenantID, c.Params("id"), &req)
	if err != nil {
		return helpers.Respond(c, labelStatus(err), err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Label updated successfully", label)
}

func (h *LabelHandler) Delete(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.hasAnyRole(c, entity.RoleTenantOwner, entity.RoleSupervisor) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only owner or supervisor can delete labels", nil)
	}
	if err := h.service.Delete(ctx, tenantID, c.Params("id")); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Label deleted successfully", nil)
}

func (h *LabelHandler) ListForConversation(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	labels, err := h.service.ListForConversation(ctx, tenantID, c.Params("conversationId"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Conversation labels retrieved successfully", labels)
}

func (h *LabelHandler) Attach(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	var req dto.AttachLabelRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}
	if err := h.service.Attach(ctx, tenantID, c.Params("conversationId"), req.LabelID); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Label attached", nil)
}

func (h *LabelHandler) Detach(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if err := h.service.Detach(ctx, tenantID, c.Params("conversationId"), c.Params("labelId")); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Label detached", nil)
}
