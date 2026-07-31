package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	ConversationStatusOpen       = "open"
	ConversationStatusUnassigned = "unassigned"
	ConversationStatusCompleted  = "completed"
	ConversationStatusArchived   = "archived"
)

type Convertation struct {
	ID                     uuid.UUID      `json:"id" db:"id"`
	TenantID               uuid.UUID      `json:"tenant_id" db:"tenant_id" validate:"required"`
	ChannelID              uuid.UUID      `json:"channel_id" db:"channel_id" validate:"required"`
	ContactID              uuid.UUID      `json:"contact_id" db:"contact_id" validate:"required"`
	AssignedAgentID        *uuid.UUID     `json:"assigned_agent_id" db:"assigned_agent_id"`
	Status                 string         `json:"status" db:"status" validate:"required,oneof=open unassigned completed archived"`
	IsPinned               bool           `json:"is_pinned"`
	IsArchived             bool           `json:"is_archived" db:"is_archived"`
	UnreadCount            int            `json:"unread_count" db:"unread_count"`
	LastMessageAt          NullTime       `json:"last_message_at" db:"last_message_at"`
	ServiceWindowExpiresAt NullTime       `json:"service_window_expires_at" db:"service_window_expires_at"`
	Metadata               *MetDataConfig `json:"metadata,omitempty" db:"metadata"`
	CreatedAt              time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at" db:"updated_at"`
}

type MetDataConfig map[string]interface{}

func (m MetDataConfig) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *MetDataConfig) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, m)
}
