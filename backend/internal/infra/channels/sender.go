package channels

import (
	"context"
	"github/socialforge/internal/entity"
)

// Sender delivers outbound messages to a provider. Each channel type has one
// implementation (Telegram Bot API, WAHA HTTP API, Meta Graph API, ...).
type Sender interface {
	// SendText sends a plain-text message to the given provider-side recipient
	// (contact external id) and returns the provider message id on success.
	SendText(ctx context.Context, channel *entity.Channel, to, text string) (providerMessageID string, err error)
}

// Connector registers/starts a channel with its provider (webhook registration,
// WAHA session start, etc.) and returns provider-specific connect info (e.g. a
// QR url for WhatsApp). status is the resulting channel status.
type Connector interface {
	Connect(ctx context.Context, channel *entity.Channel, webhookURL, secret string) (info map[string]interface{}, status string, err error)
}
