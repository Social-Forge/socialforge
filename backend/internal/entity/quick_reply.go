package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	QuickReplyTypeText     = "text"
	QuickReplyTypeImage    = "image"
	QuickReplyTypeVideo    = "video"
	QuickReplyTypeDocument = "document"
)

// QuickReplyMedia is the list of attached media (empty for text-only). Stored as
// JSONB; the column default '{}' is tolerated and read as empty.
type QuickReplyMedia []map[string]interface{}

func (m QuickReplyMedia) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "[]", nil
	}
	return json.Marshal(m)
}

func (m *QuickReplyMedia) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	s := string(b)
	if s == "" || s == "{}" {
		*m = nil
		return nil
	}
	return json.Unmarshal(b, m)
}

// QuickReply is a reusable canned response with a "/shortcut" trigger.
type QuickReply struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	TenantID    uuid.UUID       `json:"tenant_id" db:"tenant_id"`
	Shortcut    string          `json:"shortcut" db:"shortcut" validate:"required,min=1,max=255"`
	ContentType string          `json:"content_type" db:"content_type" validate:"required,oneof=text image video document"`
	Body        NullString      `json:"body,omitempty" db:"body"`
	Media       QuickReplyMedia `json:"media,omitempty" db:"media"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}
