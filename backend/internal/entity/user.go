package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	UserStatusActive    = "active"
	UserStatusSuspended = "suspended"
	UserStatusInactive  = "inactive"
)

type User struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	TenantID        uuid.UUID  `json:"tenant_id" db:"tenant_id"`
	RoleID          uuid.UUID  `json:"role_id" db:"role_id"`
	Email           string     `json:"email" db:"email" validate:"required,email,max=255"`
	Username        string     `json:"username,omitempty" db:"username"`
	PasswordHash    string     `json:"-" db:"password_hash"`
	FullName        string     `json:"full_name" db:"full_name" validate:"required,max=255"`
	Phone           NullString `json:"phone,omitempty" db:"phone" validate:"omitempty,max=20"`
	AvatarURL       NullString `json:"avatar_url,omitempty" db:"avatar_url"`
	TwoFaSecret     NullString `json:"two_fa_secret,omitempty" db:"two_fa_secret"`
	Status          string     `json:"status" db:"status" validate:"required,oneof=active suspended inactive"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	IsVerified      bool       `json:"is_verified" db:"is_verified"`
	EmailVerifiedAt NullTime   `json:"email_verified_at,omitempty" db:"email_verified_at"`
	LastLoginAt     NullTime   `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt       NullTime   `json:"deleted_at,omitempty" db:"deleted_at"`
	Role            *Role      `json:"role,omitempty"`
	Tenant          *Tenant    `json:"tenant,omitempty"`
}

type RolePermissionSummary struct {
	RoleName           string `json:"role_name,omitempty"`
	PermissionName     string `json:"permission_name,omitempty"`
	PermissionResource string `json:"permission_resource,omitempty"`
	PermissionAction   string `json:"permission_action,omitempty"`
}

type UserTenant struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	TenantID  uuid.UUID `json:"tenant_id" db:"tenant_id"`
	RoleID    uuid.UUID `json:"role_id" db:"role_id"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type UserTenantWithDetails struct {
	User           User                    `json:"user"`
	UserTenant     UserTenant              `json:"user_tenant"`
	Tenant         Tenant                  `json:"tenant"`
	Role           Role                    `json:"role"`
	RolePermissions []RolePermissionSummary `json:"role_permissions,omitempty"`
	Metadata       map[string]any          `json:"metadata,omitempty"`
}

type UserTenantWithDetailsNested = UserTenantWithDetails

func (u *UserTenantWithDetails) ToResponse() *UserResponse {
	if u == nil {
		return nil
	}
	response := u.User.ToResponse()
	response.Tenant = &u.Tenant
	response.Role = &u.Role
	response.UserTenant = &u.UserTenant
	response.RolePermissions = u.RolePermissions
	response.Metadata = u.Metadata
	return response
}

type UserResponse struct {
	ID              uuid.UUID               `json:"id"`
	TenantID        uuid.UUID               `json:"tenant_id,omitempty"`
	RoleID          uuid.UUID               `json:"role_id,omitempty"`
	Email           string                  `json:"email"`
	Username        string                  `json:"username,omitempty"`
	FullName        string                  `json:"full_name,omitempty"`
	Phone           string                  `json:"phone,omitempty"`
	AvatarURL       string                  `json:"avatar_url,omitempty"`
	TwoFaSecret     string                  `json:"two_fa_secret,omitempty"`
	IsVerified      bool                    `json:"is_verified"`
	IsActive        bool                    `json:"is_active"`
	EmailVerifiedAt NullTime                `json:"email_verified_at,omitempty"`
	LastLoginAt     NullTime                `json:"last_login_at,omitempty"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
	Tenant          *Tenant                 `json:"tenant,omitempty"`
	Role            *Role                   `json:"role,omitempty"`
	UserTenant      *UserTenant             `json:"user_tenant,omitempty"`
	RolePermissions []RolePermissionSummary `json:"role_permissions,omitempty"`
	Metadata        map[string]any          `json:"metadata,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:              u.ID,
		TenantID:        u.TenantID,
		RoleID:          u.RoleID,
		Email:           u.Email,
		Username:        u.Username,
		FullName:        u.FullName,
		Phone:           u.Phone.String,
		AvatarURL:       u.AvatarURL.String,
		TwoFaSecret:     u.TwoFaSecret.String,
		IsVerified:      u.IsVerified,
		IsActive:        u.IsActive,
		EmailVerifiedAt: u.EmailVerifiedAt,
		LastLoginAt:     u.LastLoginAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
		Tenant:          u.Tenant,
		Role:            u.Role,
	}
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt.Valid
}

func (u *User) CanLogin() bool {
	return u.IsActive && !u.IsDeleted()
}

func (u *User) IsInactive() bool {
	return u.Status == UserStatusInactive && !u.IsDeleted()
}
func (u *User) IsSuspended() bool {
	return u.Status == UserStatusSuspended && !u.IsDeleted()
}

func (u *User) MarkAsVerified() {
	now := time.Now()
	u.Status = UserStatusActive
	u.IsVerified = true
	u.EmailVerifiedAt = NewNullTimePtr(&now)
}

func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = NewNullTimePtr(&now)
}
func (u *User) IsTwoFaActive() bool {
	return u.TwoFaSecret.Valid
}
