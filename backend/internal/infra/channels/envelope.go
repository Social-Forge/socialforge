// Package channels defines the provider-agnostic ingestion contract: a single
// Unified Message Envelope that every channel normalizer produces, so the rest
// of the pipeline (dedup, persist, realtime) is identical across providers.
package channels

import "time"

const (
	ProviderWAHA     = "waha"
	ProviderTelegram = "telegram"
	ProviderMetaWA   = "whatsapp_meta"
	ProviderMessenger = "messenger"
	ProviderInstagram = "instagram"
)

// EnvelopeKind classifies an inbound event.
type EnvelopeKind string

const (
	KindMessage EnvelopeKind = "message" // a persistable inbound message
	KindCall    EnvelopeKind = "call"    // incoming call (WAHA auto-reject)
	KindStatus  EnvelopeKind = "status"  // delivery/read receipt
	KindOther   EnvelopeKind = "other"   // anything we ignore for now
)

// ContactInfo is the provider-side identity of the sender.
type ContactInfo struct {
	ExternalID  string // phone (wa), chat id (tg), PSID (messenger), IG id
	DisplayName string
	AvatarURL   string
}

// Media describes an attached media item (resolved to MinIO later, Fase 2E).
type Media struct {
	MimeType string `json:"mime_type,omitempty"`
	URL      string `json:"url,omitempty"`
	Caption  string `json:"caption,omitempty"`
	FileName string `json:"file_name,omitempty"`
	Size     int64  `json:"size,omitempty"`
	// ProviderRef is the provider-side media id/token to download later.
	ProviderRef string `json:"provider_ref,omitempty"`
}

// Envelope is the normalized form of one inbound provider event.
type Envelope struct {
	Provider        string
	Kind            EnvelopeKind
	ProviderEventID string // for webhook-level dedup (webhook_events)
	ProviderMsgID   string // for message-level dedup (messages.provider_message_id)
	Contact         ContactInfo
	ContentType     string // text | image | video | audio | document | location | sticker
	Text            string
	Media           *Media
	ReplyToMsgID    string
	Timestamp       time.Time
	EventType       string // raw provider event name (for logging/webhook_events)
	Raw             []byte // original payload
}

// IsPersistableMessage reports whether this envelope should become a message row.
func (e *Envelope) IsPersistableMessage() bool {
	return e.Kind == KindMessage && e.Contact.ExternalID != ""
}
