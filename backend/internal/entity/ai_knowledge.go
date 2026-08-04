package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// JSONStringSlice is a []string persisted as a JSONB array (keywords, asset_ids).
type JSONStringSlice []string

func (s JSONStringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

func (s *JSONStringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}
	return json.Unmarshal(bytes, s)
}

// AIKnowledge is a single knowledge-base entry (title + content) attached to an
// AI agent. `embedding` is reserved for pgvector RAG (Fase 5D); for now it stays
// an empty JSON array and retrieval is not vectorized.
type AIKnowledge struct {
	ID         uuid.UUID `json:"id" db:"id"`
	TenantID   uuid.UUID `json:"tenant_id" db:"tenant_id"`
	AIAgentID  uuid.UUID `json:"ai_agent_id" db:"ai_agent_id"`
	Title      string    `json:"title" db:"title"`
	Content    string    `json:"content" db:"content"`
	TokenCount int       `json:"token_count" db:"token_count"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
