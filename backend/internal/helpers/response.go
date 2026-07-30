package helpers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

var (
	ErrExpiredSession       string = "Session has expired, Please login"
	ErrUnauthorizedSession  string = "Unauthorized, session not found, Please login"
	ErrUnauthorizedDomain   string = "Unauthorized, domain not registered"
	ErrBadInputRequest      string = "Invalid form input"
	ErrInvalidOriginRequest string = "Access not permitted, Invalid request origin"
	ErrInvalidUuid          string = "Invalid ID format"
	ErrInvalidQueryParams   string = "Invalid query parameters"
	ErrInvalidFormData      string = "Invalid form data format"
	ErrInvalidRequestParams string = "Invalid request parameters"
)

type ApiResponse struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}
type Meta struct {
	Total int `json:"total,omitempty"`
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}
type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Stack   string `json:"stack,omitempty"` // Only include in development
}
type PaginateMeta struct {
	Total      int  `json:"total"`
	Limit      int  `json:"limit"`
	Offset     int  `json:"offset,omitempty"`
	TotalPages int  `json:"total_pages,omitempty"`
	HasMore    bool `json:"has_more,omitempty"`
}

func Respond(c fiber.Ctx, status int, message string, payload interface{}) error {
	success := status >= 200 && status < 300

	response := ApiResponse{
		Status:  status,
		Success: success,
		Message: message,
	}

	if success {
		response.Data = payload
	} else {
		response.Error = payload
	}

	return c.Status(status).JSON(response)
}
func ErrorDataNotFound(tableName string, obj interface{}) error {
	return fmt.Errorf("%s not found : %v", tableName, obj)
}
func ErrorDataConflict(tableName string, obj interface{}) error {
	return fmt.Errorf("%s already exists : %v", tableName, obj)
}
