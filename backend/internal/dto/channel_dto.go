package dto

// CreateChannelRequest is used by a tenant owner to add a messaging channel
// assigned to a division.
type CreateChannelRequest struct {
	DivisionID  string                 `json:"division_id" validate:"required,uuid4"`
	Type        string                 `json:"type" validate:"required,oneof=whatsapp_waha whatsapp_meta messenger instagram telegram"`
	Name        string                 `json:"name" validate:"required,min=2,max=255"`
	WahaEngine  string                 `json:"waha_engine,omitempty" validate:"omitempty,oneof=WEBJS NOWEB GOWS"`
	Credentials map[string]interface{} `json:"credentials,omitempty"`
}

// UpdateChannelRequest updates mutable channel fields.
type UpdateChannelRequest struct {
	Name       string  `json:"name" validate:"required,min=2,max=255"`
	DivisionID string  `json:"division_id,omitempty" validate:"omitempty,uuid4"`
	AIAgentID  *string `json:"ai_agent_id,omitempty" validate:"omitempty,uuid4"`
}
