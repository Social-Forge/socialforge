package entity

import (
	"time"

	"github.com/google/uuid"
)

// Ledger reason codes (mirror the DB CHECK constraint).
const (
	AICreditReasonGrant      = "grant"
	AICreditReasonTopup      = "topup"
	AICreditReasonDebit      = "debit"
	AICreditReasonAdjustment = "adjustment"
)

// AICreditLedger is one entry in a tenant's AI credit ledger. Every AI
// generation records a `debit` row capturing the token usage and the resulting
// balance for auditability.
type AICreditLedger struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	TenantID       uuid.UUID  `json:"tenant_id" db:"tenant_id"`
	ConversationID *uuid.UUID `json:"conversation_id,omitempty" db:"conversation_id"`
	MessageID      *uuid.UUID `json:"message_id,omitempty" db:"message_id"`
	Delta          int        `json:"delta" db:"delta"`
	BalanceAfter   int        `json:"balance_after" db:"balance_after"`
	Reason         string     `json:"reason" db:"reason"`
	Model          string     `json:"model,omitempty" db:"model"`
	InputTokens    int        `json:"input_tokens" db:"input_tokens"`
	OutputTokens   int        `json:"output_tokens" db:"output_tokens"`
	CostUSD        float64    `json:"cost_usd" db:"cost_usd"`
	Credit         float64    `json:"credit" db:"credit"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}
