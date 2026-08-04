package meta

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github/socialforge/config"
	"github/socialforge/internal/entity"
	"io"
	"net/http"
	"time"
)

// Adapter sends messages via the Meta Graph API for WABA, Messenger and
// Instagram channels. Credentials (access_token + phone_number_id/page_id) live
// on the channel.
type Adapter struct {
	cfg  *config.MetaConfig
	http *http.Client
}

func NewAdapter(cfg *config.MetaConfig) *Adapter {
	return &Adapter{cfg: cfg, http: &http.Client{Timeout: 15 * time.Second}}
}

func (a *Adapter) version() string {
	if a.cfg != nil && a.cfg.GraphAPIVersion != "" {
		return a.cfg.GraphAPIVersion
	}
	return "v18.0"
}

func cred(channel *entity.Channel, key string) string {
	if channel.Credentials == nil {
		return ""
	}
	v, _ := (*channel.Credentials)[key].(string)
	return v
}

func (a *Adapter) SendText(ctx context.Context, channel *entity.Channel, to, text string) (string, error) {
	token := cred(channel, "access_token")
	if token == "" {
		return "", fmt.Errorf("meta channel missing access_token")
	}

	var senderID string
	var payload map[string]interface{}
	switch channel.Type {
	case entity.ChannelTypeWhatsAppMeta:
		senderID = cred(channel, "phone_number_id")
		payload = map[string]interface{}{
			"messaging_product": "whatsapp",
			"to":                to,
			"type":              "text",
			"text":              map[string]interface{}{"body": text},
		}
	case entity.ChannelTypeMessenger, entity.ChannelTypeInstagram:
		senderID = cred(channel, "page_id")
		if senderID == "" {
			senderID = cred(channel, "ig_id")
		}
		payload = map[string]interface{}{
			"recipient":     map[string]interface{}{"id": to},
			"messaging_type": "RESPONSE",
			"message":       map[string]interface{}{"text": text},
		}
	default:
		return "", fmt.Errorf("unsupported meta channel type: %s", channel.Type)
	}
	if senderID == "" {
		return "", fmt.Errorf("meta channel missing sender id (phone_number_id/page_id)")
	}

	body, _ := json.Marshal(payload)
	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/messages", a.version(), senderID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := a.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("meta send failed: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("meta send rejected (%d): %s", resp.StatusCode, string(raw))
	}

	var out struct {
		Messages  []struct{ ID string `json:"id"` } `json:"messages"`
		MessageID string                            `json:"message_id"`
	}
	_ = json.Unmarshal(raw, &out)
	if len(out.Messages) > 0 {
		return out.Messages[0].ID, nil
	}
	return out.MessageID, nil
}
