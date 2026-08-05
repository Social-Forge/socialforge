package dto

// CheckoutRequest starts a payment for either a plan subscription or an addon.
type CheckoutRequest struct {
	Kind     string `json:"kind" validate:"required,oneof=subscription addon"`
	Provider string `json:"provider" validate:"required,oneof=xendit midtrans paypal"`

	// Subscription
	PlanCode string `json:"plan_code,omitempty"`
	Months   int    `json:"months,omitempty"`

	// Addon
	AddonType string `json:"addon_type,omitempty" validate:"omitempty,oneof=channel_slot agent_slot ai_credits"`
	Quantity  int    `json:"quantity,omitempty"`
}
