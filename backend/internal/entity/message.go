package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	MessageDirectionIn  = "in"
	MessageDirectionOut = "out"

	SenderTypeContact = "contact"
	SenderTypeAgent   = "agent"
	SenderTypeAI      = "ai"
	SenderTypeSystem  = "system"

	MessageStatusPending   = "pending"
	MessageStatusSent      = "sent"
	MessageStatusDelivered = "delivered"
	MessageStatusRead      = "read"
	MessageStatusFailed    = "failed"

	ContentTypeText     = "text"
	ContentTypeImage    = "image"
	ContentTypeVideo    = "video"
	ContentTypeAudio    = "audio"
	ContentTypeDocument = "document"
	ContentTypeLocation = "location"
	ContentTypeSticker  = "sticker"
	ContentTypeTemplate = "template"
)

type Message struct {
	ID                uuid.UUID              `json:"id" db:"id"`
	TenantID          uuid.UUID              `json:"tenant_id" db:"tenant_id" validate:"required"`
	ConversationID    uuid.UUID              `json:"conversation_id" db:"conversation_id" validate:"required"`
	SenderID          *uuid.UUID             `json:"sender_id,omitempty" db:"sender_id"`
	Direction         string                 `json:"direction" db:"direction" validate:"required,oneof=in out"`
	SenderType        string                 `json:"sender_type" db:"sender_type" validate:"required,oneof=contact agent ai system"`
	ContentType       string                 `json:"content_type" db:"content_type" validate:"required,oneof=text image video audio document location contact sticker button list template"`
	Body              NullString             `json:"body,omitempty" db:"body"`
	Media             map[string]interface{} `json:"media,omitempty" db:"media"`
	ProviderMessageID NullString             `json:"provider_message_id,omitempty" db:"provider_message_id"`
	Status            string                 `json:"status" db:"status" validate:"required,oneof=pending sent delivered read failed"`
	ReplyToID         *uuid.UUID             `json:"reply_to_id,omitempty" db:"reply_to_id"`
	IsPinned          bool                   `json:"is_pinned" db:"is_pinned"`
	Error             NullString             `json:"error,omitempty" db:"error"`
	EditedAt          NullTime               `json:"edited_at" db:"edited_at"`
	CreatedAt         time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at" db:"updated_at"`
	DeletedAt         NullTime               `json:"deleted_at,omitempty" db:"deleted_at"`
}
