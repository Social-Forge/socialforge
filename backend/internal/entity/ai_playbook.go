package entity

import (
	"time"

	"github.com/google/uuid"
)

// AIPlaybook is a keyword-triggered instruction attached to an AI agent. When an
// inbound message matches any keyword, the highest-priority active playbook
// steers the reply (its instruction) and optionally attaches assets.
type AIPlaybook struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	TenantID    uuid.UUID       `json:"tenant_id" db:"tenant_id"`
	AIAgentID   uuid.UUID       `json:"ai_agent_id" db:"ai_agent_id"`
	Name        string          `json:"name" db:"name"`
	Keywords    JSONStringSlice `json:"keywords" db:"keywords"`
	Instruction string          `json:"instruction" db:"instruction"`
	AssetIDs    JSONStringSlice `json:"asset_ids" db:"asset_ids"`
	Priority    int             `json:"priority" db:"priority"`
	IsActive    bool            `json:"is_active" db:"is_active"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}
