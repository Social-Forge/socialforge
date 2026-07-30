package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	ID          uuid.UUID         `json:"id" db:"id"`
	TenantID    uuid.UUID         `json:"tenant_id" db:"tenant_id" validate:"required"`
	ChannelID   uuid.UUID         `json:"channel_id" db:"channel_id" validate:"required"`
	ExternalID  string            `json:"external_id" db:"external_id" validate:"required,max=1000"`
	DisplayName string            `json:"display_name" db:"display_name" validate:"required,max=255"`
	AvatarURL   NullString        `json:"avatar_url" db:"avatar_url"`
	IsBlocked   bool              `json:"is_blocked" db:"is_blocked"`
	Attributes  *AttributesConfig `json:"attributes,omitempty" db:"attributes"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
}
type AttributesConfig map[string]interface{}

func (a AttributesConfig) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

func (a *AttributesConfig) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}
