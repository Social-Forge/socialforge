package entity

import (
	"time"

	"github.com/google/uuid"
)

// Subscription addon types (mirror the DB CHECK).
const (
	AddonTypeChannelSlot = "channel_slot"
	AddonTypeAgentSlot   = "agent_slot"
	AddonTypeAICredits   = "ai_credits"
)

// SubscriptionAddon is a purchased top-up on top of the base plan: extra channel
// slots, agent seats, or AI credits. `meta` carries type-specific detail
// (e.g. {"channelType":"telegram"} or {"credits":10000}).
type SubscriptionAddon struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TenantID  uuid.UUID `json:"tenant_id" db:"tenant_id"`
	Type      string    `json:"type" db:"type"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Meta      JSONMap   `json:"meta" db:"meta"`
	ExpiresAt NullTime  `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
