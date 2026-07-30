package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleLevelTenantOwner = 0
	RoleLevelSupervisor  = 1
	RoleLevelAgent       = 2
	RoleLevelSuperAdmin  = 3
)
const (
	RoleTenantOwner = "tenant_owner"
	RoleSupervisor  = "supervisor"
	RoleAgent       = "agent"
	RoleSuperAdmin  = "superadmin"
)

type Role struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name" validate:"required,max=50"`
	Slug        string     `json:"slug" db:"slug" validate:"required,max=50"`
	Description NullString `json:"description,omitempty" db:"description"`
	Level       int        `json:"level" db:"level" validate:"required,min=0,max=3"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// TableName returns the table name for Role
func (Role) TableName() string {
	return "roles"
}

// IsSuperAdmin checks if role is superadmin
func (r *Role) IsSuperAdmin() bool {
	return r.Level == RoleLevelSuperAdmin
}

// IsTenantOwner checks if role is tenant owner
func (r *Role) IsTenantOwner() bool {
	return r.Level == RoleLevelTenantOwner
}

// CanManageTenant checks if role can manage tenant
func (r *Role) CanManageTenant() bool {
	return r.Level <= RoleLevelTenantOwner
}

// HasHigherLevelThan checks if this role has higher level than another role
func (r *Role) HasHigherLevelThan(other *Role) bool {
	return r.Level < other.Level // Lower number = higher level
}
