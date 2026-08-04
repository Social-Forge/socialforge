package dto

// CreateAIAgentRequest configures an AI customer-service agent's identity,
// persona and guardrails.
type CreateAIAgentRequest struct {
	Name         string                 `json:"name" validate:"required,min=1,max=255"`
	Provider     string                 `json:"provider" validate:"required,oneof=claude openai google"`
	Model        string                 `json:"model,omitempty"`
	SystemPrompt string                 `json:"system_prompt" validate:"required,min=1"`
	Persona      map[string]interface{} `json:"persona,omitempty"`    // name, soul, tone, gender, style, greeting
	Safety       map[string]interface{} `json:"safety,omitempty"`     // sensitive topics -> handoff
	Guardrails   map[string]interface{} `json:"guardrails,omitempty"` // list of instructions/limits
	Temperature  float64                `json:"temperature,omitempty" validate:"omitempty,min=0,max=2"`
	MaxTokens    int                    `json:"max_tokens,omitempty" validate:"omitempty,min=64,max=8192"`
	AutoReply    *bool                  `json:"auto_reply_enabled,omitempty"`
}

type UpdateAIAgentRequest struct {
	Name         string                 `json:"name" validate:"required,min=1,max=255"`
	Provider     string                 `json:"provider" validate:"required,oneof=claude openai google"`
	Model        string                 `json:"model,omitempty"`
	SystemPrompt string                 `json:"system_prompt" validate:"required,min=1"`
	Persona      map[string]interface{} `json:"persona,omitempty"`
	Safety       map[string]interface{} `json:"safety,omitempty"`
	Guardrails   map[string]interface{} `json:"guardrails,omitempty"`
	Temperature  float64                `json:"temperature,omitempty" validate:"omitempty,min=0,max=2"`
	MaxTokens    int                    `json:"max_tokens,omitempty" validate:"omitempty,min=64,max=8192"`
	AutoReply    *bool                  `json:"auto_reply_enabled,omitempty"`
	IsActive     *bool                  `json:"is_active,omitempty"`
}
