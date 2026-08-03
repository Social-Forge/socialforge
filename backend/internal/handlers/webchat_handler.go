package handlers

import (
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type webchatMessageRequest struct {
	VisitorID string `json:"visitor_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Text      string `json:"text" validate:"required,min=1,max=8000"`
}

type WebchatHandler struct {
	ctxinject   *middlewares.ContextMiddleware
	ingestion   *services.IngestionService
	wsURL       string
	logger      *zap.Logger
}

func NewWebchatHandler(ctxinject *middlewares.ContextMiddleware, ingestion *services.IngestionService, wsURL string, logger *zap.Logger) *WebchatHandler {
	return &WebchatHandler{ctxinject: ctxinject, ingestion: ingestion, wsURL: wsURL, logger: logger}
}

// Session hands the widget a stable visitor id + the realtime endpoint.
func (h *WebchatHandler) Session(c fiber.Ctx) error {
	return helpers.Respond(c, fiber.StatusOK, "Webchat session", fiber.Map{
		"visitor_id": "visitor_" + uuid.NewString(),
		"ws_url":     h.wsURL,
	})
}

// Message ingests a visitor's inbound message through the standard pipeline.
func (h *WebchatHandler) Message(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req webchatMessageRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}
	if req.VisitorID == "" {
		req.VisitorID = "visitor_" + uuid.NewString()
	}

	conversationID, err := h.ingestion.IngestWebchat(ctx, c.Params("channelId"), req.VisitorID, req.Name, req.Text)
	if err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	return helpers.Respond(c, fiber.StatusCreated, "Message received", fiber.Map{
		"visitor_id":       req.VisitorID,
		"conversation_id":  conversationID,
		"realtime_channel": "conversation:" + conversationID,
	})
}
