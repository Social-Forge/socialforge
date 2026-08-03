package entity

import (
	"time"

	"github.com/google/uuid"
)

// AutoResponse is the per-channel automatic first reply sent to a brand-new
// customer (one not yet in the contact database), unless an AI agent is active.
type AutoResponse struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	TenantID    uuid.UUID       `json:"tenant_id" db:"tenant_id"`
	ChannelID   uuid.UUID       `json:"channel_id" db:"channel_id"`
	IsEnabled   bool            `json:"is_enabled" db:"is_enabled"`
	ContentType string          `json:"content_type" db:"content_type" validate:"required,oneof=text image video document"`
	Body        NullString      `json:"body,omitempty" db:"body"`
	Media       QuickReplyMedia `json:"media,omitempty" db:"media"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}
