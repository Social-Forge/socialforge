package dto

import (
	"errors"
	"github/socialforge/internal/entity"
)

var (
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrPhoneAlreadyExists    = errors.New("phone already exists")
	ErrWeakPassword          = errors.New("password too weak")
)

type UpdateProfileRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"email"`
	Username string `json:"username" validate:"required"`
	Phone    string `json:"phone" validate:"omitempty,e164"`
}
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
type EnableTwoFactorRequest struct {
	Status bool `json:"status" validate:"boolean"`
}
type ActivateTwoFactorRequest struct {
	Code string `json:"code" validate:"required"`
}
type UserProfileResponse struct {
	User      *entity.User       `json:"user"`
	Tenant    *entity.Tenant     `json:"tenant"`
	Role      *entity.Role       `json:"role"`
	Divisions []*entity.Division `json:"divisions"`
}

type UserTenantResponse struct {
	User      *entity.User   `json:"user"`
	Tenant    *entity.Tenant `json:"tenant"`
	Role      *entity.Role   `json:"role"`
	Divisions []string       `json:"divisions"` // List of division names/slugs
}
