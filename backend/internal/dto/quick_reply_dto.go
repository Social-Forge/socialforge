package dto

// CreateQuickReplyRequest creates a canned reply. media carries already-uploaded
// file references (url/mime_type/...) for non-text types; text-only uses body.
type CreateQuickReplyRequest struct {
	Shortcut    string                   `json:"shortcut" validate:"required,min=1,max=255"`
	ContentType string                   `json:"content_type" validate:"required,oneof=text image video document"`
	Body        string                   `json:"body,omitempty" validate:"omitempty,max=8000"`
	Media       []map[string]interface{} `json:"media,omitempty"`
}

type UpdateQuickReplyRequest struct {
	Shortcut    string                   `json:"shortcut" validate:"required,min=1,max=255"`
	ContentType string                   `json:"content_type" validate:"required,oneof=text image video document"`
	Body        string                   `json:"body,omitempty" validate:"omitempty,max=8000"`
	Media       []map[string]interface{} `json:"media,omitempty"`
}
