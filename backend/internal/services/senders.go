package services

import (
	"github/socialforge/config"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/channels"
	"github/socialforge/internal/infra/channels/telegram"
	"github/socialforge/internal/infra/channels/waha"
)

// BuildSenders wires the provider send adapters keyed by channel type. Meta
// channels (whatsapp_meta/messenger/instagram) are added in Fase 4.
func BuildSenders(cfg *config.Config) map[string]channels.Sender {
	return map[string]channels.Sender{
		entity.ChannelTypeTelegram:     telegram.NewAdapter(),
		entity.ChannelTypeWhatsAppWaha: waha.NewAdapter(&cfg.Waha),
	}
}

// BuildConnectors wires the provider connect adapters keyed by channel type.
func BuildConnectors(cfg *config.Config) map[string]channels.Connector {
	return map[string]channels.Connector{
		entity.ChannelTypeTelegram:     telegram.NewAdapter(),
		entity.ChannelTypeWhatsAppWaha: waha.NewAdapter(&cfg.Waha),
	}
}
