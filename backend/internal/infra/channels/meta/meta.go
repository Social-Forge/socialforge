// Package meta implements the Meta (Facebook) channels: WhatsApp Business Cloud
// API (WABA), Messenger and Instagram. They share the Graph API + webhook shape.
package meta

import (
	"encoding/json"
	"fmt"
	"github/socialforge/internal/infra/channels"
	"strconv"
	"time"
)

type webhook struct {
	Object string  `json:"object"`
	Entry  []entry `json:"entry"`
}

type entry struct {
	ID        string      `json:"id"`
	Changes   []change    `json:"changes"`   // WABA
	Messaging []messaging `json:"messaging"` // Messenger / Instagram
}

type change struct {
	Field string `json:"field"`
	Value struct {
		Contacts []struct {
			Profile struct {
				Name string `json:"name"`
			} `json:"profile"`
			WaID string `json:"wa_id"`
		} `json:"contacts"`
		Messages []struct {
			From      string `json:"from"`
			ID        string `json:"id"`
			Timestamp string `json:"timestamp"`
			Type      string `json:"type"`
			Text      struct {
				Body string `json:"body"`
			} `json:"text"`
		} `json:"messages"`
	} `json:"value"`
}

type messaging struct {
	Sender    struct{ ID string } `json:"sender"`
	Timestamp int64               `json:"timestamp"`
	Message   struct {
		Mid    string `json:"mid"`
		Text   string `json:"text"`
		IsEcho bool   `json:"is_echo"`
	} `json:"message"`
}

// Normalize converts a Meta webhook payload (WABA / Messenger / Instagram) into
// a Unified Envelope. The provider is derived from the payload object type.
func Normalize(body []byte) (*channels.Envelope, error) {
	var wh webhook
	if err := json.Unmarshal(body, &wh); err != nil {
		return nil, fmt.Errorf("invalid meta payload: %w", err)
	}

	env := &channels.Envelope{Kind: channels.KindOther, Raw: body, EventType: wh.Object}
	if len(wh.Entry) == 0 {
		return env, nil
	}
	e := wh.Entry[0]

	switch wh.Object {
	case "whatsapp_business_account":
		env.Provider = channels.ProviderMetaWA
		if len(e.Changes) == 0 || len(e.Changes[0].Value.Messages) == 0 {
			return env, nil
		}
		v := e.Changes[0].Value
		m := v.Messages[0]
		env.Kind = channels.KindMessage
		env.ProviderMsgID = m.ID
		env.ProviderEventID = m.ID
		env.Contact.ExternalID = m.From
		if len(v.Contacts) > 0 {
			env.Contact.DisplayName = v.Contacts[0].Profile.Name
		}
		if ts, err := strconv.ParseInt(m.Timestamp, 10, 64); err == nil {
			env.Timestamp = time.Unix(ts, 0)
		}
		env.ContentType = "text"
		env.Text = m.Text.Body

	case "page":
		env.Provider = channels.ProviderMessenger
		return fillMessaging(env, e)
	case "instagram":
		env.Provider = channels.ProviderInstagram
		return fillMessaging(env, e)
	}

	if env.Contact.DisplayName == "" {
		env.Contact.DisplayName = env.Contact.ExternalID
	}
	return env, nil
}

func fillMessaging(env *channels.Envelope, e entry) (*channels.Envelope, error) {
	if len(e.Messaging) == 0 {
		return env, nil
	}
	msg := e.Messaging[0]
	if msg.Message.IsEcho || msg.Message.Mid == "" {
		return env, nil // outbound echo or non-message event
	}
	env.Kind = channels.KindMessage
	env.ProviderMsgID = msg.Message.Mid
	env.ProviderEventID = msg.Message.Mid
	env.Contact.ExternalID = msg.Sender.ID
	env.Contact.DisplayName = msg.Sender.ID
	env.Timestamp = time.Unix(msg.Timestamp/1000, 0)
	env.ContentType = "text"
	env.Text = msg.Message.Text
	return env, nil
}
