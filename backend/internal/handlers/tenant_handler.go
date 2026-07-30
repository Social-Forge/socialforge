package handlers

import (
	"github/socialforge/internal/dto"
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type TenantHandler struct {
	ctxinject *middlewares.ContextMiddleware
	service   *services.TenantService
	logger    *zap.Logger
}

func NewTenantHandler(ctxinject *middlewares.ContextMiddleware, service *services.TenantService, logger *zap.Logger) *TenantHandler {
	return &TenantHandler{
		ctxinject: ctxinject,
		service:   service,
		logger:    logger,
	}
}

func (h *TenantHandler) UpdateLogo(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid tenant ID", nil)
	}

	c.Request().Header.Set("Content-Type", "multipart/form-data")

	form, err := c.MultipartForm()
	if err != nil {
		h.logger.Error("Failed to parse multipart form",
			zap.String("tenant_id", tenantID),
			zap.Error(err))
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid form data", nil)
	}
	defer form.RemoveAll()

	files, exists := form.File["logo"]
	if !exists || len(files) == 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, "Logo file is required", nil)
	}

	logoFile := files[0]

	const maxFileSize = 5 << 20 // 5MB
	if logoFile.Size > maxFileSize {
		return helpers.Respond(c, fiber.StatusBadRequest,
			"File too large. Maximum size is 5MB", nil)
	}

	fileHeader := logoFile.Header.Get("Content-Type")
	if !h.isValidImageType(fileHeader) {
		return helpers.Respond(c, fiber.StatusBadRequest,
			"Invalid image type. Only JPEG, JPG, PNG, GIF, and WebP are allowed", nil)
	}

	logoURL, err := h.service.ChangeLogo(ctx, tenantID, logoFile)
	if err != nil {
		h.logger.Error("Failed to change logo",
			zap.String("tenant_id", tenantID),
			zap.Error(err))

		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.Respond(c, fiber.StatusOK, "Logo updated successfully", fiber.Map{
		"logo_url": logoURL,
	})
}
func (h *TenantHandler) UpdateInfo(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	tenantID := c.Params("tenantID")
	if tenantID == "" {
		return helpers.Respond(c, fiber.StatusBadRequest, "Tenant ID is required", nil)
	}

	var req dto.UpdateTenantRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, "Invalid request payload", nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}

	updateTenant, err := h.service.UpdateInfo(ctx, tenantID, &req)
	if err != nil {
		h.logger.Error("Failed to update tenant info",
			zap.String("tenant_id", tenantID),
			zap.Error(err))
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.Respond(c, fiber.StatusOK, "Tenant info updated successfully", updateTenant)
}

func (h *TenantHandler) isValidImageType(contentType string) bool {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	return allowedTypes[contentType]
}
