package handlers

import (
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type sendMessageRequest struct {
	Text string `json:"text" validate:"required,min=1,max=8000"`
}

type ConversationHandler struct {
	ctxinject *middlewares.ContextMiddleware
	convSvc   *services.ConversationService
	outbound  *services.OutboundService
	logger    *zap.Logger
}

func NewConversationHandler(
	ctxinject *middlewares.ContextMiddleware,
	convSvc *services.ConversationService,
	outbound *services.OutboundService,
	logger *zap.Logger,
) *ConversationHandler {
	return &ConversationHandler{ctxinject: ctxinject, convSvc: convSvc, outbound: outbound, logger: logger}
}

func (h *ConversationHandler) tenantID(c fiber.Ctx) (string, bool) {
	tid, ok := c.Locals("tenant_id").(string)
	if !ok || tid == "" {
		return "", false
	}
	return tid, true
}

func (h *ConversationHandler) List(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	convs, err := h.convSvc.List(ctx, tenantID, c.Query("status"), c.Query("agent_id"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Conversations retrieved successfully", convs)
}

func (h *ConversationHandler) ListMessages(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	limit, _ := strconv.Atoi(c.Query("limit"))
	messages, err := h.convSvc.ListMessages(ctx, tenantID, c.Params("id"), limit)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Messages retrieved successfully", messages)
}

func (h *ConversationHandler) SendMessage(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	agentID, _ := c.Locals("user_id").(string)

	var req sendMessageRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}

	msg, err := h.outbound.SendText(ctx, tenantID, c.Params("id"), agentID, req.Text)
	if err != nil {
		h.logger.Error("Failed to send message", zap.String("tenant_id", tenantID), zap.Error(err))
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusCreated, "Message sent", msg)
}
