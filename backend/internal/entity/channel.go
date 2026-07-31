package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	ChannelTypeWhatsAppWaha = "whatsapp_waha"
	ChannelTypeWhatsAppMeta = "whatsapp_meta"
	ChannelTypeMessenger    = "messenger"
	ChannelTypeInstagram    = "instagram"
	ChannelTypeTelegram     = "telegram"
)

const (
	ChannelStatusDisconnected = "disconnected"
	ChannelStatusConnected    = "connected"
	ChannelStatusConnecting   = "connecting"
	ChannelStatusFailed       = "failed"
)

type Channel struct {
	ID              uuid.UUID          `json:"id" db:"id"`
	TenantID        uuid.UUID          `json:"tenant_id" db:"tenant_id" validate:"required"`
	DivisionID      uuid.UUID          `json:"division_id" db:"division_id" validate:"required"`
	AIAgentID       *uuid.UUID         `json:"ai_agent_id,omitempty" db:"ai_agent_id"`
	Type            string             `json:"type" db:"type" validate:"required,oneof=whatsapp_waha whatsapp_meta messenger instagram telegram"`
	Name            string             `json:"name" db:"name" validate:"required,max=255"`
	Status          string             `json:"status" db:"status" validate:"required,oneof=disconnected connected connecting failed"`
	ExternalID      NullString         `json:"external_id,omitempty" db:"external_id"`
	WahaEngine      NullString         `json:"waha_engine,omitempty" db:"waha_engine"`
	WahaSessionName NullString         `json:"waha_session_name,omitempty" db:"waha_session_name"`
	WebhookSecret   NullString         `json:"webhook_secret,omitempty" db:"webhook_secret"`
	Credentials     *CredentialsConfig `json:"credentials,omitempty" db:"credentials"`
	Settings        *ChannelSetting    `json:"settings,omitempty" db:"settings"`
	CreatedAt       time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" db:"updated_at"`
}

type CredentialsConfig map[string]interface{}
type ChannelSetting map[string]interface{}

func (c CredentialsConfig) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

func (c *CredentialsConfig) Scan(value interface{}) error {
	if value == nil {
		*c = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, c)
}
func (c ChannelSetting) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

func (c *ChannelSetting) Scan(value interface{}) error {
	if value == nil {
		*c = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, c)
}
