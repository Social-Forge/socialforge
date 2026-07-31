package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	OutboxStatusPending    = "pending"
	OutboxStatusProcessing = "processing"
	OutboxStatusSent       = "sent"
	OutboxStatusFailed     = "failed"
)

// MessageOutbox tracks reliable delivery of outbound messages to providers.
type MessageOutbox struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	TenantID    uuid.UUID  `json:"tenant_id" db:"tenant_id"`
	MessageID   uuid.UUID  `json:"message_id" db:"message_id"`
	Status      string     `json:"status" db:"status"`
	Attempts    int        `json:"attempts" db:"attempts"`
	MaxAttempts int        `json:"max_attempts" db:"max_attempts"`
	NextRetryAt NullTime   `json:"next_retry_at" db:"next_retry_at"`
	LastError   NullString `json:"last_error,omitempty" db:"last_error"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}
