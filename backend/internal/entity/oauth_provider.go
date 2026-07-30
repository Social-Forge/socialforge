package entity

import (
	"time"

	"github.com/google/uuid"
)

type OAuthProvider struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id" validate:"required"`
	ProviderName string    `json:"provider_name" db:"provider_name" validate:"required"`
	ProviderID   string    `json:"provider_id" db:"provider_id" validate:"required, oneof=google facebook twitter"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
