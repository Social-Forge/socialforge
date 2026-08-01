package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github/socialforge/internal/entity"
	"net/http"
	"strconv"
	"time"
)

// Adapter sends messages via the Telegram Bot API. The bot token is per-channel
// (stored in channel.credentials.bot_token).
type Adapter struct {
	http *http.Client
}

func NewAdapter() *Adapter {
	return &Adapter{http: &http.Client{Timeout: 15 * time.Second}}
}

func botToken(channel *entity.Channel) (string, error) {
	if channel.Credentials == nil {
		return "", fmt.Errorf("telegram channel has no credentials")
	}
	tok, _ := (*channel.Credentials)["bot_token"].(string)
	if tok == "" {
		return "", fmt.Errorf("telegram channel missing bot_token")
	}
	return tok, nil
}

// Connect registers our webhook URL + secret token with Telegram so the bot's
// updates are delivered to POST /api/webhooks/telegram/{channel_id}.
func (a *Adapter) Connect(ctx context.Context, channel *entity.Channel, webhookURL, secret string) (map[string]interface{}, string, error) {
	token, err := botToken(channel)
	if err != nil {
		return nil, entity.ChannelStatusFailed, err
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"url":             webhookURL,
		"secret_token":    secret,
		"allowed_updates": []string{"message", "edited_message"},
	})
	url := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, entity.ChannelStatusFailed, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.http.Do(req)
	if err != nil {
		return nil, entity.ChannelStatusFailed, fmt.Errorf("telegram setWebhook failed: %w", err)
	}
	defer resp.Body.Close()

	var out struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if !out.OK {
		return nil, entity.ChannelStatusFailed, fmt.Errorf("telegram setWebhook rejected: %s", out.Description)
	}
	return map[string]interface{}{"webhook_url": webhookURL}, entity.ChannelStatusConnected, nil
}

func (a *Adapter) SendText(ctx context.Context, channel *entity.Channel, to, text string) (string, error) {
	token, err := botToken(channel)
	if err != nil {
		return "", err
	}

	payload, _ := json.Marshal(map[string]interface{}{"chat_id": to, "text": text})
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("telegram send failed: %w", err)
	}
	defer resp.Body.Close()

	var out struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
		Result      struct {
			MessageID int64 `json:"message_id"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("telegram: decode response: %w", err)
	}
	if !out.OK {
		return "", fmt.Errorf("telegram send rejected (%d): %s", resp.StatusCode, out.Description)
	}
	return strconv.FormatInt(out.Result.MessageID, 10), nil
}
