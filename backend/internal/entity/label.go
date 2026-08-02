package entity

import (
	"time"

	"github.com/google/uuid"
)

// Label is a tenant-scoped tag used to mark conversations (e.g. "new customer",
// "pending payment", "delivered"). Names are unique within a tenant.
type Label struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TenantID  uuid.UUID `json:"tenant_id" db:"tenant_id"`
	Name      string    `json:"name" db:"name" validate:"required,min=1,max=255"`
	Color     string    `json:"color" db:"color" validate:"omitempty,max=32"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
