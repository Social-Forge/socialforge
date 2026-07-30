package dto

import (
	"errors"
)

var (
	ErrTenantNotFound = errors.New("tenant not found")
)

type UpdateTenantRequest struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
}
