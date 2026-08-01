package waha

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github/socialforge/config"
	"github/socialforge/internal/entity"
	"io"
	"net/http"
	"strings"
	"time"
)

// Adapter sends messages via a WAHA server. WAHA runs one server per engine
// (webjs/noweb/gows); the channel picks the engine and session name, and the
// per-engine API key comes from config.
type Adapter struct {
	cfg  *config.WahaConfig
	http *http.Client
}

func NewAdapter(cfg *config.WahaConfig) *Adapter {
	return &Adapter{cfg: cfg, http: &http.Client{Timeout: 20 * time.Second}}
}

// baseURL returns the in-network WAHA URL for an engine. In docker-compose the
// engines are reachable at waha-{engine}:3000.
func baseURL(engine string) string {
	switch strings.ToUpper(engine) {
	case "GOWS":
		return "http://waha-gows:3000"
	case "NOWEB":
		return "http://waha-noweb:3000"
	default:
		return "http://waha-webjs:3000"
	}
}

func (a *Adapter) apiKey(engine string) string {
	switch strings.ToUpper(engine) {
	case "GOWS":
		return a.cfg.GowsAPiKey
	case "NOWEB":
		return a.cfg.NowebAPiKey
	default:
		return a.cfg.WebJsAPiKey
	}
}

// toChatID converts a bare number/JID to a WhatsApp chat id.
func toChatID(to string) string {
	if strings.Contains(to, "@") {
		return to
	}
	return to + "@c.us"
}

// Connect creates/starts the WAHA session with our webhook configured, and
// returns the QR endpoint to scan for WhatsApp login. status is "connecting"
// until the QR is scanned (WAHA then emits session.status events).
func (a *Adapter) Connect(ctx context.Context, channel *entity.Channel, webhookURL, secret string) (map[string]interface{}, string, error) {
	if !channel.WahaSessionName.Valid || channel.WahaSessionName.String == "" {
		return nil, entity.ChannelStatusFailed, fmt.Errorf("waha channel has no session name")
	}
	engine := "WEBJS"
	if channel.WahaEngine.Valid {
		engine = channel.WahaEngine.String
	}
	base := baseURL(engine)
	session := channel.WahaSessionName.String

	payload, _ := json.Marshal(map[string]interface{}{
		"name":  session,
		"start": true,
		"config": map[string]interface{}{
			"webhooks": []map[string]interface{}{
				{
					"url":    webhookURL,
					"events": []string{"message", "call.received", "session.status"},
					"hmac":   map[string]interface{}{"key": secret},
				},
			},
		},
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/api/sessions", bytes.NewReader(payload))
	if err != nil {
		return nil, entity.ChannelStatusFailed, err
	}
	req.Header.Set("Content-Type", "application/json")
	if k := a.apiKey(engine); k != "" {
		req.Header.Set("X-Api-Key", k)
	}

	resp, err := a.http.Do(req)
	if err != nil {
		return nil, entity.ChannelStatusFailed, fmt.Errorf("waha start session failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	// 201 created, or 422/409 if it already exists — both are acceptable to proceed.
	if resp.StatusCode >= 500 {
		return nil, entity.ChannelStatusFailed, fmt.Errorf("waha start session error (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return map[string]interface{}{
		"session": session,
		"qr_url":  fmt.Sprintf("%s/api/%s/auth/qr?format=image", base, session),
	}, entity.ChannelStatusConnecting, nil
}

func (a *Adapter) SendText(ctx context.Context, channel *entity.Channel, to, text string) (string, error) {
	if !channel.WahaSessionName.Valid || channel.WahaSessionName.String == "" {
		return "", fmt.Errorf("waha channel has no session")
	}
	engine := "WEBJS"
	if channel.WahaEngine.Valid {
		engine = channel.WahaEngine.String
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"session": channel.WahaSessionName.String,
		"chatId":  toChatID(to),
		"text":    text,
	})

	url := baseURL(engine) + "/api/sendText"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if k := a.apiKey(engine); k != "" {
		req.Header.Set("X-Api-Key", k)
	}

	resp, err := a.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("waha send failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("waha send rejected (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var out struct {
		ID string `json:"id"`
	}
	_ = json.Unmarshal(body, &out)
	return out.ID, nil
}
