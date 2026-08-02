package entity

import (
	"time"

	"github.com/google/uuid"
)

// ConversationLabel attaches a Label to a Conversation.
type ConversationLabel struct {
	ID             uuid.UUID `json:"id" db:"id"`
	TenantID       uuid.UUID `json:"tenant_id" db:"tenant_id"`
	ConversationID uuid.UUID `json:"conversation_id" db:"conversation_id"`
	LabelID        uuid.UUID `json:"label_id" db:"label_id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
