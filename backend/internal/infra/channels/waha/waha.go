package waha

import (
	"encoding/json"
	"fmt"
	"github/socialforge/internal/infra/channels"
	"strings"
	"time"
)

type webhook struct {
	Event   string  `json:"event"`
	Session string  `json:"session"`
	Payload payload `json:"payload"`
}

type payload struct {
	ID         string `json:"id"`
	Timestamp  int64  `json:"timestamp"`
	From       string `json:"from"`
	FromMe     bool   `json:"fromMe"`
	Body       string `json:"body"`
	HasMedia   bool   `json:"hasMedia"`
	NotifyName string `json:"notifyName"`
	Media      *media `json:"media"`
	Data       struct {
		NotifyName string `json:"notifyName"`
	} `json:"_data"`
}

type media struct {
	URL      string `json:"url"`
	MimeType string `json:"mimetype"`
	Filename string `json:"filename"`
}

// Normalize converts a WAHA webhook payload into a Unified Envelope.
func Normalize(body []byte) (*channels.Envelope, error) {
	var wh webhook
	if err := json.Unmarshal(body, &wh); err != nil {
		return nil, fmt.Errorf("invalid waha payload: %w", err)
	}

	env := &channels.Envelope{
		Provider:        channels.ProviderWAHA,
		Kind:            channels.KindOther,
		EventType:       wh.Event,
		ProviderEventID: wh.Payload.ID,
		Raw:             body,
	}

	switch {
	case strings.HasPrefix(wh.Event, "call."):
		env.Kind = channels.KindCall
		env.Contact.ExternalID = normalizeJID(wh.Payload.From)
		return env, nil
	case wh.Event == "message" || wh.Event == "message.any":
		if wh.Payload.FromMe {
			// Outbound echo — not an inbound message to ingest.
			return env, nil
		}
	default:
		return env, nil // acks, session status, etc. — ignored
	}

	env.Kind = channels.KindMessage
	env.ProviderMsgID = wh.Payload.ID
	env.Timestamp = time.Unix(wh.Payload.Timestamp, 0)
	env.Contact.ExternalID = normalizeJID(wh.Payload.From)

	name := wh.Payload.NotifyName
	if name == "" {
		name = wh.Payload.Data.NotifyName
	}
	if name == "" {
		name = env.Contact.ExternalID
	}
	env.Contact.DisplayName = name

	if wh.Payload.HasMedia && wh.Payload.Media != nil {
		env.ContentType = mediaContentType(wh.Payload.Media.MimeType)
		env.Text = wh.Payload.Body
		env.Media = &channels.Media{
			URL:      wh.Payload.Media.URL,
			MimeType: wh.Payload.Media.MimeType,
			FileName: wh.Payload.Media.Filename,
			Caption:  wh.Payload.Body,
		}
	} else {
		env.ContentType = "text"
		env.Text = wh.Payload.Body
	}
	return env, nil
}

// normalizeJID strips the WhatsApp JID suffix (@c.us / @s.whatsapp.net) so the
// contact external id is the bare number; group JIDs (@g.us) are kept as-is.
func normalizeJID(jid string) string {
	if strings.HasSuffix(jid, "@g.us") {
		return jid
	}
	if i := strings.IndexByte(jid, '@'); i >= 0 {
		return jid[:i]
	}
	return jid
}

func mediaContentType(mime string) string {
	switch {
	case strings.HasPrefix(mime, "image/"):
		return "image"
	case strings.HasPrefix(mime, "video/"):
		return "video"
	case strings.HasPrefix(mime, "audio/"):
		return "audio"
	default:
		return "document"
	}
}
