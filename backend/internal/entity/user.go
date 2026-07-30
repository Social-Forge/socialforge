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
	ID              uuid.UUID    `json:"id" db:"id"`
	Email           string       `json:"email" db:"email" validate:"required,email,max=255"`
	PasswordHash    string       `json:"-" db:"password_hash"`
	FullName        string       `json:"full_name" db:"full_name" validate:"required,max=255"`
	Phone           NullString   `json:"phone,omitempty" db:"phone" validate:"omitempty,max=20"`
	AvatarURL       NullString   `json:"avatar_url,omitempty" db:"avatar_url"`
	TwoFaSecret     NullString   `json:"two_fa_secret,omitempty" db:"two_fa_secret"`
	Status          string       `json:"status" db:"status" validate:"required,oneof=active suspended inactive"`
	EmailVerifiedAt NullTime     `json:"email_verified_at,omitempty" db:"email_verified_at"`
	LastLoginAt     NullTime     `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt       time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt       NullTime     `json:"deleted_at,omitempty" db:"deleted_at"`
	UserTenants     []UserTenant `json:"user_tenants,omitempty" db:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt.Valid
}

func (u *User) CanLogin() bool {
	return u.Status == UserStatusActive && !u.IsDeleted()
}

func (u *User) IsInactive() bool {
	return u.Status == UserStatusInactive && !u.IsDeleted()
}
func (u *User) IsSuspended() bool {
	return u.Status == UserStatusSuspended && !u.IsDeleted()
}

// IsActive is derived from status (there is no is_active column). It is the
// single source of truth for "active" replacing the removed is_active field.
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive && !u.IsDeleted()
}

// IsVerified is derived from email_verified_at (there is no is_verified column).
func (u *User) IsVerified() bool {
	return u.EmailVerifiedAt.Valid
}

func (u *User) MarkAsVerified() {
	now := time.Now()
	u.Status = UserStatusActive
	u.EmailVerifiedAt = NewNullTimePtr(&now)
}

func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = NewNullTimePtr(&now)
}
func (u *User) IsTwoFaActive() bool {
	return u.TwoFaSecret.Valid
}
