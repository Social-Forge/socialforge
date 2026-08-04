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

type AIResourceHandler struct {
	ctxinject *middlewares.ContextMiddleware
	service   *services.AIResourceService
	logger    *zap.Logger
}

func NewAIResourceHandler(ctxinject *middlewares.ContextMiddleware, service *services.AIResourceService, logger *zap.Logger) *AIResourceHandler {
	return &AIResourceHandler{ctxinject: ctxinject, service: service, logger: logger}
}

func (h *AIResourceHandler) tenantID(c fiber.Ctx) (string, bool) {
	tid, ok := c.Locals("tenant_id").(string)
	if !ok || tid == "" {
		return "", false
	}
	return tid, true
}

func (h *AIResourceHandler) isOwner(c fiber.Ctx) bool {
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

func (h *AIResourceHandler) guard(c fiber.Ctx, mutate bool) (string, bool) {
	tid, ok := h.tenantID(c)
	if !ok {
		_ = helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
		return "", false
	}
	if mutate && !h.isOwner(c) {
		_ = helpers.Respond(c, fiber.StatusForbidden, "Only the tenant owner can manage AI resources", nil)
		return "", false
	}
	return tid, true
}

func bindValidate[T any](c fiber.Ctx, req *T) bool {
	if err := c.Bind().Body(req); err != nil {
		_ = helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
		return false
	}
	if errs := helpers.ValidateStruct(*req); len(errs) > 0 {
		_ = helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
		return false
	}
	return true
}

// ---------------- Knowledge ----------------

func (h *AIResourceHandler) ListKnowledge(c fiber.Ctx) error {
	tid, ok := h.guard(c, false)
	if !ok {
		return nil
	}
	items, err := h.service.ListKnowledge(h.ctxinject.HandlerContext(c), tid, c.Params("id"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Knowledge retrieved successfully", items)
}

func (h *AIResourceHandler) CreateKnowledge(c fiber.Ctx) error {
	tid, ok := h.guard(c, true)
	if !ok {
		return nil
	}
	var req dto.CreateAIKnowledgeRequest
	if !bindValidate(c, &req) {
		return nil
	}
	item, err := h.service.CreateKnowledge(h.ctxinject.HandlerContext(c), tid, c.Params("id"), &req)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusCreated, "Knowledge created successfully", item)
}

func (h *AIResourceHandler) UpdateKnowledge(c fiber.Ctx) error {
	tid, ok := h.guard(c, true)
	if !ok {
		return nil
	}
	var req dto.UpdateAIKnowledgeRequest
	if !bindValidate(c, &req) {
		return nil
	}
	item, err := h.service.UpdateKnowledge(h.ctxinject.HandlerContext(c), tid, c.Params("id"), c.Params("childId"), &req)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Knowledge updated successfully", item)
}

func (h *AIResourceHandler) DeleteKnowledge(c fiber.Ctx) error {
	tid, ok := h.guard(c, true)
	if !ok {
		return nil
	}
	if err := h.service.DeleteKnowledge(h.ctxinject.HandlerContext(c), tid, c.Params("id"), c.Params("childId")); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Knowledge deleted successfully", nil)
}

// ---------------- Playbook ----------------

func (h *AIResourceHandler) ListPlaybooks(c fiber.Ctx) error {
	tid, ok := h.guard(c, false)
	if !ok {
		return nil
	}
	items, err := h.service.ListPlaybooks(h.ctxinject.HandlerContext(c), tid, c.Params("id"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Playbooks retrieved successfully", items)
}

func (h *AIResourceHandler) CreatePlaybook(c fiber.Ctx) error {
	tid, ok := h.guard(c, true)
	if !ok {
		return nil
	}
	var req dto.CreateAIPlaybookRequest
	if !bindValidate(c, &req) {
		return nil
	}
	item, err := h.service.CreatePlaybook(h.ctxinject.HandlerContext(c), tid, c.Params("id"), &req)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusCreated, "Playbook created successfully", item)
}

func (h *AIResourceHandler) UpdatePlaybook(c fiber.Ctx) error {
	tid, ok := h.guard(c, true)
	if !ok {
		return nil
	}
	var req dto.UpdateAIPlaybookRequest
	if !bindValidate(c, &req) {
		return nil
	}
	item, err := h.service.UpdatePlaybook(h.ctxinject.HandlerContext(c), tid, c.Params("id"), c.Params("childId"), &req)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Playbook updated successfully", item)
}

func (h *AIResourceHandler) DeletePlaybook(c fiber.Ctx) error {
	tid, ok := h.guard(c, true)
	if !ok {
		return nil
	}
	if err := h.service.DeletePlaybook(h.ctxinject.HandlerContext(c), tid, c.Params("id"), c.Params("childId")); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Playbook deleted successfully", nil)
}

// ---------------- Asset ----------------

func (h *AIResourceHandler) ListAssets(c fiber.Ctx) error {
	tid, ok := h.guard(c, false)
	if !ok {
		return nil
	}
	items, err := h.service.ListAssets(h.ctxinject.HandlerContext(c), tid, c.Params("id"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Assets retrieved successfully", items)
}

func (h *AIResourceHandler) CreateAsset(c fiber.Ctx) error {
	tid, ok := h.guard(c, true)
	if !ok {
		return nil
	}
	var req dto.CreateAIAssetRequest
	if !bindValidate(c, &req) {
		return nil
	}
	item, err := h.service.CreateAsset(h.ctxinject.HandlerContext(c), tid, c.Params("id"), &req)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusCreated, "Asset created successfully", item)
}

func (h *AIResourceHandler) UpdateAsset(c fiber.Ctx) error {
	tid, ok := h.guard(c, true)
	if !ok {
		return nil
	}
	var req dto.UpdateAIAssetRequest
	if !bindValidate(c, &req) {
		return nil
	}
	item, err := h.service.UpdateAsset(h.ctxinject.HandlerContext(c), tid, c.Params("id"), c.Params("childId"), &req)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Asset updated successfully", item)
}

func (h *AIResourceHandler) DeleteAsset(c fiber.Ctx) error {
	tid, ok := h.guard(c, true)
	if !ok {
		return nil
	}
	if err := h.service.DeleteAsset(h.ctxinject.HandlerContext(c), tid, c.Params("id"), c.Params("childId")); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Asset deleted successfully", nil)
}
