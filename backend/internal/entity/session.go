package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id" validate:"required"`
	AccessToken    string     `json:"access_token" db:"access_token" validate:"required"`
	RefreshToken   string     `json:"refresh_token" db:"refresh_token" validate:"required"`
	ExpiresAt      time.Time  `json:"expires_at" db:"expires_at" validate:"required"`
	IPAddress      *string    `json:"ip_address" db:"ip_address"`
	UserAgent      *string    `json:"user_agent" db:"user_agent"`
	LastActivityAt *time.Time `json:"last_activity_at" db:"last_activity_at"`
	IsRevoked      bool       `json:"is_revoked" db:"is_revoked"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

func (Session) TableName() string {
	return "sessions"
}
func (t *Session) IsExpired() bool {
	return t.ExpiresAt.Before(time.Now())
}
func (t *Session) IsSessionRevoked() bool {
	return t.IsRevoked
}
