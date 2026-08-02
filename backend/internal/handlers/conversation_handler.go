package handlers

import (
	"context"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/infra/repository"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/services"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type sendMessageRequest struct {
	Text string `json:"text" validate:"required,min=1,max=8000"`
}

type assignRequest struct {
	AgentID string `json:"agent_id" validate:"required,uuid4"`
}

type editMessageRequest struct {
	Text string `json:"text" validate:"required,min=1,max=8000"`
}

type forwardMessageRequest struct {
	TargetConversationID string `json:"target_conversation_id" validate:"required,uuid4"`
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

func (h *ConversationHandler) hasAnyRole(c fiber.Ctx, roles ...string) bool {
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

// action is a helper for simple conversation mutation endpoints.
func (h *ConversationHandler) action(c fiber.Ctx, okMsg string, fn func(ctx context.Context, tenantID, id string) error) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if err := fn(ctx, tenantID, c.Params("id")); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, okMsg, nil)
}

func (h *ConversationHandler) Assign(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if !h.hasAnyRole(c, entity.RoleTenantOwner, entity.RoleSupervisor) {
		return helpers.Respond(c, fiber.StatusForbidden, "Only owner or supervisor can assign conversations", nil)
	}
	var req assignRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}
	if err := h.convSvc.Assign(ctx, tenantID, c.Params("id"), req.AgentID); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Conversation assigned", nil)
}

func (h *ConversationHandler) Unassign(c fiber.Ctx) error {
	return h.action(c, "Conversation unassigned", h.convSvc.Unassign)
}

// --- Message actions ---

func (h *ConversationHandler) PinMessage(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if err := h.convSvc.SetMessagePinned(ctx, tenantID, c.Params("id"), c.Params("messageId"), true); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Message pinned", nil)
}

func (h *ConversationHandler) UnpinMessage(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if err := h.convSvc.SetMessagePinned(ctx, tenantID, c.Params("id"), c.Params("messageId"), false); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Message unpinned", nil)
}

func (h *ConversationHandler) EditMessage(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	var req editMessageRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}
	msg, err := h.convSvc.EditMessage(ctx, tenantID, c.Params("id"), c.Params("messageId"), req.Text)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Message edited", msg)
}

func (h *ConversationHandler) DeleteMessage(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	if err := h.convSvc.DeleteMessage(ctx, tenantID, c.Params("id"), c.Params("messageId")); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Message deleted", nil)
}

func (h *ConversationHandler) ForwardMessage(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	agentID, _ := c.Locals("user_id").(string)

	var req forwardMessageRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}

	src, err := h.convSvc.GetMessage(ctx, tenantID, c.Params("messageId"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusNotFound, err.Error(), nil)
	}
	msg, err := h.outbound.SendText(ctx, tenantID, req.TargetConversationID, agentID, src.Body.String)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusCreated, "Message forwarded", msg)
}
func (h *ConversationHandler) Complete(c fiber.Ctx) error {
	return h.action(c, "Conversation marked completed", h.convSvc.Complete)
}
func (h *ConversationHandler) Reopen(c fiber.Ctx) error {
	return h.action(c, "Conversation reopened", h.convSvc.Reopen)
}
func (h *ConversationHandler) MarkRead(c fiber.Ctx) error {
	return h.action(c, "Conversation marked read", h.convSvc.MarkRead)
}
func (h *ConversationHandler) Pin(c fiber.Ctx) error {
	return h.action(c, "Conversation pinned", func(ctx context.Context, t, id string) error { return h.convSvc.SetPinned(ctx, t, id, true) })
}
func (h *ConversationHandler) Unpin(c fiber.Ctx) error {
	return h.action(c, "Conversation unpinned", func(ctx context.Context, t, id string) error { return h.convSvc.SetPinned(ctx, t, id, false) })
}
func (h *ConversationHandler) Archive(c fiber.Ctx) error {
	return h.action(c, "Conversation archived", func(ctx context.Context, t, id string) error { return h.convSvc.SetArchived(ctx, t, id, true) })
}
func (h *ConversationHandler) Unarchive(c fiber.Ctx) error {
	return h.action(c, "Conversation unarchived", func(ctx context.Context, t, id string) error { return h.convSvc.SetArchived(ctx, t, id, false) })
}

func (h *ConversationHandler) List(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}

	f := repository.ConversationListFilter{
		Status: c.Query("status"),
		Search: c.Query("search"),
	}
	if v := c.Query("channel_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			f.ChannelID = &id
		}
	}
	if v := c.Query("agent_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			f.AgentID = &id
		}
	}
	if v := c.Query("label_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			f.LabelID = &id
		}
	}
	if v := c.Query("archived"); v != "" {
		b := v == "true" || v == "1"
		f.Archived = &b
	}
	if v := c.Query("date_from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.DateFrom = t
		}
	}
	if v := c.Query("date_to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.DateTo = t
		}
	}

	convs, err := h.convSvc.List(ctx, tenantID, f)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Conversations retrieved successfully", convs)
}

func (h *ConversationHandler) Unread(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)
	tenantID, ok := h.tenantID(c)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant context is required", nil)
	}
	total, err := h.convSvc.TotalUnread(ctx, tenantID, c.Query("agent_id"))
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusOK, "Unread total retrieved", fiber.Map{"total_unread": total})
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
