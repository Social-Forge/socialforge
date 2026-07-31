package entity

import (
	"time"

	"github.com/google/uuid"
)

// WebhookEvent records an inbound provider event for idempotent processing.
// A unique (channel_id, provider_event_id) prevents double-processing when a
// provider redelivers the same event.
type WebhookEvent struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	ChannelID       uuid.UUID  `json:"channel_id" db:"channel_id"`
	Provider        string     `json:"provider" db:"provider"`
	ProviderEventID string     `json:"provider_event_id" db:"provider_event_id"`
	EventType       NullString `json:"event_type,omitempty" db:"event_type"`
	Payload         []byte     `json:"-" db:"payload"`
	ProcessedAt     NullTime   `json:"processed_at,omitempty" db:"processed_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}
