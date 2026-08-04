package entity

import (
	"time"

	"github.com/google/uuid"
)

// AI asset types (mirror the DB CHECK constraint).
const (
	AIAssetTypeImage    = "image"
	AIAssetTypeVideo    = "video"
	AIAssetTypeDocument = "document"
)

// AIAsset is a reusable media object (stored in MinIO under StorageKey) that an
// AI agent can send to customers, referenced by playbooks via asset_ids. A fresh
// presigned URL is minted at send time from StorageKey.
type AIAsset struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	TenantID    uuid.UUID  `json:"tenant_id" db:"tenant_id"`
	AIAgentID   uuid.UUID  `json:"ai_agent_id" db:"ai_agent_id"`
	Name        string     `json:"name" db:"name"`
	Type        string     `json:"type" db:"type"`
	StorageKey  string     `json:"storage_key" db:"storage_key"`
	MimeType    NullString `json:"mime_type" db:"mime_type"`
	Size        NullInt32  `json:"size" db:"size"`
	Description NullString `json:"description" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}
