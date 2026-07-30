package entity

import (
	"time"

	"github.com/google/uuid"
)

type DivisionMember struct {
	ID           uuid.UUID   `json:"id" db:"id"`
	UserTenantID uuid.UUID   `json:"user_tenant_id" db:"user_tenant_id"`
	DivisionID   uuid.UUID   `json:"division_id" db:"division_id"`
	IsActive     bool        `json:"is_active" db:"is_active"`
	JoinedAt     time.Time   `json:"joined_at" db:"joined_at"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
	DeletedAt    NullTime    `json:"deleted_at,omitempty" db:"deleted_at"`
	UserTenant   *UserTenant `json:"user_tenant,omitempty" db:"-"`
	Division     *Division   `json:"division,omitempty" db:"-"`
}
