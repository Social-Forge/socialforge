package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserTenant struct {
	ID              uuid.UUID        `json:"id" db:"id"`
	UserID          uuid.UUID        `json:"user_id" db:"user_id"`
	TenantID        uuid.UUID        `json:"tenant_id" db:"tenant_id"`
	RoleID          uuid.UUID        `json:"role_id" db:"role_id"`
	IsActive        bool             `json:"is_active" db:"is_active"`
	CreatedAt       time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt       NullTime         `json:"deleted_at,omitempty" db:"deleted_at"`
	User            *User            `json:"user,omitempty" db:"-"`
	Tenant          *Tenant          `json:"tenant,omitempty" db:"-"`
	Role            *Role            `json:"role,omitempty" db:"-"`
	DivisionMembers []DivisionMember `json:"division_members,omitempty" db:"-"`
}
