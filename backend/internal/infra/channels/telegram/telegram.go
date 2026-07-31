package telegram

import (
	"encoding/json"
	"fmt"
	"github/socialforge/internal/infra/channels"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Normalize converts a Telegram webhook Update payload into a Unified Envelope.
func Normalize(body []byte) (*channels.Envelope, error) {
	var update tgbotapi.Update
	if err := json.Unmarshal(body, &update); err != nil {
		return nil, fmt.Errorf("invalid telegram update: %w", err)
	}

	env := &channels.Envelope{
		Provider:        channels.ProviderTelegram,
		Kind:            channels.KindOther,
		ProviderEventID: strconv.Itoa(update.UpdateID),
		EventType:       "update",
		Raw:             body,
	}

	msg := update.Message
	if msg == nil {
		msg = update.EditedMessage
	}
	if msg == nil {
		return env, nil // non-message update (callback, etc.) — ignored for now
	}

	env.Kind = channels.KindMessage
	env.EventType = "message"
	env.ProviderMsgID = strconv.Itoa(msg.MessageID)
	env.Timestamp = time.Unix(int64(msg.Date), 0)

	if msg.Chat != nil {
		env.Contact.ExternalID = strconv.FormatInt(msg.Chat.ID, 10)
	}
	if msg.From != nil {
		name := msg.From.FirstName
		if msg.From.LastName != "" {
			name += " " + msg.From.LastName
		}
		if name == "" {
			name = msg.From.UserName
		}
		env.Contact.DisplayName = name
	}
	if env.Contact.DisplayName == "" {
		env.Contact.DisplayName = env.Contact.ExternalID
	}

	switch {
	case msg.Text != "":
		env.ContentType = "text"
		env.Text = msg.Text
	case len(msg.Photo) > 0:
		env.ContentType = "image"
		env.Text = msg.Caption
		env.Media = &channels.Media{ProviderRef: msg.Photo[len(msg.Photo)-1].FileID, Caption: msg.Caption}
	case msg.Video != nil:
		env.ContentType = "video"
		env.Text = msg.Caption
		env.Media = &channels.Media{ProviderRef: msg.Video.FileID, Caption: msg.Caption, MimeType: msg.Video.MimeType}
	case msg.Voice != nil:
		env.ContentType = "audio"
		env.Media = &channels.Media{ProviderRef: msg.Voice.FileID, MimeType: msg.Voice.MimeType}
	case msg.Audio != nil:
		env.ContentType = "audio"
		env.Media = &channels.Media{ProviderRef: msg.Audio.FileID, MimeType: msg.Audio.MimeType}
	case msg.Document != nil:
		env.ContentType = "document"
		env.Text = msg.Caption
		env.Media = &channels.Media{ProviderRef: msg.Document.FileID, FileName: msg.Document.FileName, MimeType: msg.Document.MimeType, Caption: msg.Caption}
	case msg.Location != nil:
		env.ContentType = "location"
		env.Text = fmt.Sprintf("%f,%f", msg.Location.Latitude, msg.Location.Longitude)
	case msg.Sticker != nil:
		env.ContentType = "sticker"
		env.Media = &channels.Media{ProviderRef: msg.Sticker.FileID}
	default:
		// Unsupported content -> keep as empty text so the conversation still updates.
		env.ContentType = "text"
	}

	if msg.ReplyToMessage != nil {
		env.ReplyToMsgID = strconv.Itoa(msg.ReplyToMessage.MessageID)
	}
	return env, nil
}
