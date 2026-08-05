package handlers

import (
	"github/socialforge/internal/dto"
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type PlanHandler struct {
	ctxinject *middlewares.ContextMiddleware
	service   *services.PlanService
	logger    *zap.Logger
}

func NewPlanHandler(ctxinject *middlewares.ContextMiddleware, service *services.PlanService, logger *zap.Logger) *PlanHandler {
	return &PlanHandler{ctxinject: ctxinject, service: service, logger: logger}
}

// List returns the public plan catalog (active plans only).
func (h *PlanHandler) List(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	plans, err := h.service.List(ctx, true)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Plans retrieved successfully", plans)
}

func (h *PlanHandler) Get(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	plan, err := h.service.GetByCode(ctx, c.Params("code"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusNotFound, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Plan retrieved successfully", plan)
}

// ListAll returns all plans including inactive (superadmin).
func (h *PlanHandler) ListAll(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	plans, err := h.service.List(ctx, false)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Plans retrieved successfully", plans)
}

func (h *PlanHandler) Create(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	var req dto.CreatePlanRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}
	plan, err := h.service.Create(ctx, &req)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusCreated, "Plan created successfully", plan)
}

func (h *PlanHandler) Update(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	var req dto.UpdatePlanRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}
	plan, err := h.service.Update(ctx, c.Params("id"), &req)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Plan updated successfully", plan)
}

func (h *PlanHandler) Delete(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	if err := h.service.Delete(ctx, c.Params("id")); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Plan deleted successfully", nil)
}
