package handlers

import (
	"github/socialforge/internal/dto"
	"github/socialforge/internal/helpers"
	"github/socialforge/internal/middlewares"
	"github/socialforge/internal/services"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type AuthHandler struct {
	ctxinject   *middlewares.ContextMiddleware
	authService *services.AuthService
	rateLimiter *middlewares.RateLimiterMiddleware
	logger      *zap.Logger
}

func NewAuthHandler(
	ctxinject *middlewares.ContextMiddleware,
	authService *services.AuthService,
	rateLimiter *middlewares.RateLimiterMiddleware,
	logger *zap.Logger,
) *AuthHandler {
	return &AuthHandler{
		ctxinject:   ctxinject,
		authService: authService,
		rateLimiter: rateLimiter,
		logger:      logger,
	}
}
func (h *AuthHandler) Register(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req dto.RegisterUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}

	_, err := h.authService.Register(ctx, &req)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.Respond(c, fiber.StatusCreated, "User registered successfully", nil)
}
func (h *AuthHandler) Login(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req dto.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}
	ip, ok := c.Locals("real_ip").(string)
	if !ok {
		ip = c.IP()
	}
	platform, ok := c.Locals("platform").(string)
	if !ok {
		platform = "browser"
	}

	response, err := h.authService.Login(ctx, &req, ip, platform)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	h.rateLimiter.ResetLimitCounters(c)

	switch response.Status {
	case "require_email_verification":
		return helpers.Respond(c, fiber.StatusForbidden, "Email verification required", response)
	case "two_fa_required":
		return helpers.Respond(c, fiber.StatusAccepted, "Two-factor authentication required", response)
	default:
		return helpers.Respond(c, fiber.StatusOK, "Login successful", response)
	}
}
func (h *AuthHandler) VerifyEmail(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req dto.VerifyEmailRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}

	if err := h.authService.VerifyEmail(ctx, req.Token); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.Respond(c, fiber.StatusOK, "Email verified successfully", nil)
}
func (h *AuthHandler) ForgotPassword(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req dto.ForgotPasswordRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}

	ip, ok := c.Locals("real_ip").(string)
	if !ok {
		ip = c.IP()
	}

	if err := h.authService.ForgotPassword(ctx, &req, ip); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	h.rateLimiter.ResetLimitCounters(c)

	return helpers.Respond(c, fiber.StatusOK, "Password reset email sent, please check your email", nil)
}
func (h *AuthHandler) ResetPassword(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req dto.ResetPasswordRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}

	if err := h.authService.ResetPassword(ctx, &req); err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.Respond(c, fiber.StatusOK, "Password reset successfully", nil)
}
func (h *AuthHandler) VerifyTwoFactor(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req dto.VerifyTwoFactorRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}

	ip, ok := c.Locals("real_ip").(string)
	if !ok {
		ip = c.IP()
	}
	platform, ok := c.Locals("platform").(string)
	if !ok {
		platform = "browser"
	}

	payload, err := h.authService.VerifyTwoFactor(ctx, &req, ip, platform)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	h.rateLimiter.ResetLimitCounters(c)

	return helpers.Respond(c, fiber.StatusOK, "Two-factor authentication verified successfully", payload)
}
func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	ctx := h.ctxinject.HandlerContext(c)

	var req dto.RefreshTokenRequest
	if err := c.Bind().Body(&req); err != nil {
		return helpers.Respond(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if errs := helpers.ValidateStruct(req); len(errs) > 0 {
		return helpers.Respond(c, fiber.StatusBadRequest, helpers.ValidationErrors{Errors: errs}.Error(), nil)
	}

	platform, ok := c.Locals("platform").(string)
	if !ok {
		platform = "browser"
	}

	payload, err := h.authService.RefreshToken(ctx, req.RefreshToken, platform)
	if err != nil {
		return helpers.Respond(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	h.rateLimiter.ResetLimitCounters(c)

	return helpers.Respond(c, fiber.StatusOK, "Token refreshed successfully", payload)
}
