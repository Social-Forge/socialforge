package entity

import (
	"time"

	"github.com/google/uuid"
)

// PaymentEvent is an audit record of a provider webhook callback (idempotency +
// traceability). external_id is the provider event/reference id used for dedup.
type PaymentEvent struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	TenantID   uuid.UUID  `json:"tenant_id" db:"tenant_id"`
	InvoiceID  *uuid.UUID `json:"invoice_id,omitempty" db:"invoice_id"`
	Provider   string     `json:"provider" db:"provider"`
	EventType  string     `json:"event_type" db:"event_type"`
	ExternalID NullString `json:"external_id" db:"external_id"`
	Payload    JSONMap    `json:"payload" db:"payload"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}
