package entity

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse is the API-facing shape of a user, flattened with the active
// tenant/role context. is_active/is_verified are computed (see User.IsActive /
// User.IsVerified) since the users table only stores status + email_verified_at.
type UserResponse struct {
	ID              uuid.UUID   `json:"id"`
	TenantID        uuid.UUID   `json:"tenant_id,omitempty"`
	RoleID          uuid.UUID   `json:"role_id,omitempty"`
	Email           string      `json:"email"`
	FullName        string      `json:"full_name"`
	Phone           string      `json:"phone,omitempty"`
	AvatarURL       string      `json:"avatar_url,omitempty"`
	TwoFaSecret     string      `json:"-"`
	Status          string      `json:"status,omitempty"`
	IsActive        bool        `json:"is_active"`
	IsVerified      bool        `json:"is_verified"`
	EmailVerifiedAt NullTime    `json:"email_verified_at,omitempty"`
	LastLoginAt     NullTime    `json:"last_login_at,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	Tenant          *Tenant     `json:"tenant,omitempty"`
	UserTenant      *UserTenant `json:"user_tenant,omitempty"`
	Role            *Role       `json:"role,omitempty"`
}

// NewUserResponse builds a UserResponse from a user entity and its tenant/role
// context. Pass nil for tenant/userTenant/role when they are not needed.
func NewUserResponse(u *User, tenant *Tenant, userTenant *UserTenant, role *Role) *UserResponse {
	if u == nil {
		return nil
	}
	resp := &UserResponse{
		ID:              u.ID,
		Email:           u.Email,
		FullName:        u.FullName,
		Phone:           u.Phone.String,
		AvatarURL:       u.AvatarURL.String,
		TwoFaSecret:     u.TwoFaSecret.String,
		Status:          u.Status,
		IsActive:        u.IsActive(),
		IsVerified:      u.IsVerified(),
		EmailVerifiedAt: u.EmailVerifiedAt,
		LastLoginAt:     u.LastLoginAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
		Tenant:          tenant,
		UserTenant:      userTenant,
		Role:            role,
	}
	if tenant != nil {
		resp.TenantID = tenant.ID
	}
	if role != nil {
		resp.RoleID = role.ID
	}
	return resp
}

// UserTenantWithDetails is the aggregate loaded for authentication/session
// building: the user together with its active membership, tenant and role.
// Permissions were removed in favour of level/name-based roles, so this holds
// no permission list — role name + level is the authority.
type UserTenantWithDetails struct {
	User       User           `json:"user"`
	UserTenant UserTenant     `json:"user_tenant"`
	Tenant     Tenant         `json:"tenant"`
	Role       Role           `json:"role"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

// ToResponse flattens the aggregate into the API-facing UserResponse.
func (d *UserTenantWithDetails) ToResponse() *UserResponse {
	return NewUserResponse(&d.User, &d.Tenant, &d.UserTenant, &d.Role)
}

// RoleNames returns the role identifiers used for authorization checks
// (name, plus slug when it differs).
func (d *UserTenantWithDetails) RoleNames() []string {
	names := []string{d.Role.Name}
	if d.Role.Slug != "" && d.Role.Slug != d.Role.Name {
		names = append(names, d.Role.Slug)
	}
	return names
}
