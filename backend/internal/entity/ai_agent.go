package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AIAgent struct {
	ID               uuid.UUID                 `json:"id" db:"id"`
	TenantID         uuid.UUID                 `json:"tenant_id" db:"tenant_id" validate:"required"`
	Name             string                    `json:"name" db:"name" validate:"required"`
	Provider         string                    `json:"provider" db:"provider" validate:"required, oneof=claude openai google"`
	Model            string                    `json:"model" db:"model" validate:"required"`
	SystemPrompt     string                    `json:"system_prompt" db:"system_prompt"`
	Persona          *AiPersonaConfig          `json:"persona" db:"persona"`
	Safety           *AiSafetyConfig           `json:"safety" db:"safety"`
	Guardrails       *AiSafetyGuardrailsConfig `json:"guardrails" db:"guardrails"`
	Temperature      float64                   `json:"temperature" db:"temperature"`
	MaxTokens        int                       `json:"max_tokens" db:"max_tokens"`
	AutoReplyEnabled bool                      `json:"auto_reply_enabled" db:"auto_reply_enabled"`
	WorkingHours     *WorkingHours             `json:"working_hours" db:"working_hours"`
	IsActive         bool                      `json:"is_active" db:"is_active"`
	CreatedAt        time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time                 `json:"updated_at" db:"updated_at"`
}

type AiPersonaConfig map[string]interface{}
type AiSafetyConfig map[string]interface{}
type AiSafetyGuardrailsConfig map[string]interface{}
type WorkingHours map[string]interface{}

func (a AiPersonaConfig) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}
func (a *AiPersonaConfig) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}
func (a AiSafetyConfig) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}
func (a *AiSafetyConfig) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}
func (a AiSafetyGuardrailsConfig) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}
func (a *AiSafetyGuardrailsConfig) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}
func (w WorkingHours) Value() (driver.Value, error) {
	if w == nil {
		return nil, nil
	}
	return json.Marshal(w)
}
func (w *WorkingHours) Scan(value interface{}) error {
	if value == nil {
		*w = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, w)
}
