package entity

import (
	"time"

	"github.com/google/uuid"
)

type DivisionMember struct {
	ID         uuid.UUID `json:"id" db:"id"`
	TenantID   uuid.UUID `json:"tenant_id" db:"tenant_id" validate:"required"`
	DivisionID uuid.UUID `json:"division_id" db:"division_id" validate:"required"`
	UserID     uuid.UUID `json:"user_id" db:"user_id" validate:"required"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
