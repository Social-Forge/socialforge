package entity

import (
	"time"

	"github.com/google/uuid"
)

// Invoice statuses + providers (mirror the DB CHECK constraints).
const (
	InvoiceStatusPending = "pending"
	InvoiceStatusPaid    = "paid"
	InvoiceStatusExpired = "expired"
	InvoiceStatusFailed  = "failed"

	PaymentProviderXendit   = "xendit"
	PaymentProviderMidtrans = "midtrans"
	PaymentProviderPaypal   = "paypal"
)

// Invoice is a payment request for a plan subscription or an addon top-up. The
// `purpose` JSONB carries what the payment is for, e.g.
// {"kind":"subscription","plan_code":"pro"} or
// {"kind":"addon","addon_type":"ai_credits","quantity":10000}.
type Invoice struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	TenantID          uuid.UUID  `json:"tenant_id" db:"tenant_id"`
	Number            int        `json:"number" db:"number"`
	Status            string     `json:"status" db:"status"`
	Amount            int        `json:"amount" db:"amount"`
	Currency          string     `json:"currency" db:"currency"`
	Description       string     `json:"description" db:"description"`
	Purpose           JSONMap    `json:"purpose" db:"purpose"`
	Provider          string     `json:"provider" db:"provider"`
	ProviderInvoiceID NullString `json:"provider_invoice_id" db:"provider_invoice_id"`
	CheckoutURL       NullString `json:"checkout_url" db:"checkout_url"`
	PaidAt            NullTime   `json:"paid_at" db:"paid_at"`
	ExpiresAt         NullTime   `json:"expires_at" db:"expires_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// Invoice purpose kinds.
const (
	InvoicePurposeSubscription = "subscription"
	InvoicePurposeAddon        = "addon"
)
