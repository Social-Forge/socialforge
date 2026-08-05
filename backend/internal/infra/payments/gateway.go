// Package payments provides a provider-agnostic payment gateway abstraction
// (Xendit, Midtrans, PayPal) for creating checkout invoices and verifying
// webhook callbacks.
package payments

import "context"

// CreateInvoiceParams is the normalized input for creating a hosted checkout.
type CreateInvoiceParams struct {
	ExternalID  string // our invoice id (round-tripped back in the webhook)
	Number      int    // human-friendly invoice number
	Amount      int    // whole IDR (or minor units for USD providers)
	Currency    string
	Description string
	PayerEmail  string
	SuccessURL  string
	FailureURL  string
}

// CreatedInvoice is what the provider returns after creating a checkout.
type CreatedInvoice struct {
	ProviderInvoiceID string
	CheckoutURL       string
}

// WebhookResult is the normalized outcome of a verified provider callback.
type WebhookResult struct {
	ProviderInvoiceID string // provider's invoice id
	ExternalID        string // our invoice id, if the provider echoes it
	EventID           string // provider event id (audit/dedup)
	EventType         string
	Paid              bool // true when this callback confirms payment
	Expired           bool // true when the checkout expired/failed
}

// Gateway is implemented by each payment provider adapter.
type Gateway interface {
	Name() string
	CreateInvoice(ctx context.Context, p CreateInvoiceParams) (*CreatedInvoice, error)
	// VerifyWebhook authenticates and parses a raw provider callback. It must
	// return an error if the signature/token is invalid.
	VerifyWebhook(headers map[string]string, body []byte) (*WebhookResult, error)
}
