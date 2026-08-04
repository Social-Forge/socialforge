package dto

// ---- Knowledge ----

type CreateAIKnowledgeRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=500"`
	Content string `json:"content" validate:"required,min=1"`
}

type UpdateAIKnowledgeRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=500"`
	Content string `json:"content" validate:"required,min=1"`
}

// ---- Playbook ----

type CreateAIPlaybookRequest struct {
	Name        string   `json:"name" validate:"required,min=1,max=255"`
	Keywords    []string `json:"keywords" validate:"required,min=1,dive,min=1"`
	Instruction string   `json:"instruction" validate:"required,min=1"`
	AssetIDs    []string `json:"asset_ids,omitempty" validate:"omitempty,dive,uuid"`
	Priority    int      `json:"priority,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

type UpdateAIPlaybookRequest struct {
	Name        string   `json:"name" validate:"required,min=1,max=255"`
	Keywords    []string `json:"keywords" validate:"required,min=1,dive,min=1"`
	Instruction string   `json:"instruction" validate:"required,min=1"`
	AssetIDs    []string `json:"asset_ids,omitempty" validate:"omitempty,dive,uuid"`
	Priority    int      `json:"priority,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

// ---- Asset ----
// Asset metadata CRUD. The binary is uploaded to MinIO separately (refinement);
// storage_key references the object. type ∈ image|video|document.

type CreateAIAssetRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Type        string `json:"type" validate:"required,oneof=image video document"`
	StorageKey  string `json:"storage_key" validate:"required,min=1"`
	MimeType    string `json:"mime_type,omitempty"`
	Size        int    `json:"size,omitempty" validate:"omitempty,min=0"`
	Description string `json:"description,omitempty"`
}

type UpdateAIAssetRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Type        string `json:"type" validate:"required,oneof=image video document"`
	StorageKey  string `json:"storage_key" validate:"required,min=1"`
	MimeType    string `json:"mime_type,omitempty"`
	Size        int    `json:"size,omitempty" validate:"omitempty,min=0"`
	Description string `json:"description,omitempty"`
}
