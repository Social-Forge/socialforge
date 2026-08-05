package entity

import (
	"time"

	"github.com/google/uuid"
)

// Subscription statuses (mirror the DB CHECK). Note the historical typo
// 'trailing' (= trialing) is preserved to match the constraint.
const (
	SubscriptionStatusTrialing = "trailing"
	SubscriptionStatusActive   = "active"
	SubscriptionStatusPastDue  = "past_due"
	SubscriptionStatusCanceled = "canceled"
	SubscriptionStatusExpired  = "expired"
)

// Subscription is the billing record binding a tenant to a plan for a period.
type Subscription struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	TenantID           uuid.UUID `json:"tenant_id" db:"tenant_id"`
	PlanID             uuid.UUID `json:"plan_id" db:"plan_id"`
	Status             string    `json:"status" db:"status"`
	CurrentPeriodStart NullTime  `json:"current_period_start" db:"current_period_start"`
	CurrentPeriodEnd   NullTime  `json:"current_period_end" db:"current_period_end"`
	CancelAtPeriodEnd  bool      `json:"cancel_at_period_end" db:"cancel_at_period_end"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}
