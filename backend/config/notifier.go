package config

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Notifier interface {
	SendAlert(alert AlertRequest)
	Shutdown()
}
type AlertRequest struct {
	Subject  string
	Message  string
	Metadata map[string]interface{}
}
type unifiedNotifier struct {
	telegramConfig *TelegramConfig
	workers        int
	queue          chan AlertRequest
	wg             sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
	tgBot          *tgbotapi.BotAPI
	tgEnabled      bool
}

func NewUnifiedNotifier(workers int, queueSize int, Cooldown time.Duration, cfg *TelegramConfig) (Notifier, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if cfg == nil {
		return nil, fmt.Errorf("telegram config is nil")
	}
	un := &unifiedNotifier{
		telegramConfig: cfg,
		workers:        workers,
		queue:          make(chan AlertRequest, queueSize),
		ctx:            ctx,
		cancel:         cancel,
	}

	if cfg.BotToken != "" {
		bot, err := tgbotapi.NewBotAPIWithClient(cfg.BotToken, tgbotapi.APIEndpoint, &http.Client{
			Timeout: 10 * time.Second,
		})
		if err != nil {
			Logger.Warn("Telegram notifier disabled - initialization failed",
				zap.Error(err),
				zap.String("token_prefix", safeTokenPrefix(cfg.BotToken)),
			)
			un.tgEnabled = false
		} else {
			un.tgBot = bot
			un.tgEnabled = true
			Logger.Info("Telegram notifier initialized",
				zap.String("bot_username", bot.Self.UserName),
			)
		}
	} else {
		Logger.Warn("Telegram notifier disabled - no token provided")
	}

	if un.tgEnabled {
		un.wg.Add(workers)
		for i := 0; i < workers; i++ {
			go un.worker()
		}
	} else {
		Logger.Warn("No notifiers available - alert system will be disabled")
	}

	return un, nil
}

func (un *unifiedNotifier) SendAlert(alert AlertRequest) {
	select {
	case un.queue <- alert:
	case <-un.ctx.Done():
		Logger.Debug("Alert system is shutdown")
	}
}
func (un *unifiedNotifier) Shutdown() {
	un.cancel()
	un.wg.Wait()
	Logger.Info("Alert manager closed successfully")
	Logger.Info("Telegram Instance closed successfully")
}
func (un *unifiedNotifier) worker() {
	defer un.wg.Done()

	for {
		select {
		case alert := <-un.queue:
			un.processAlert(alert)
		case <-un.ctx.Done():
			return
		}
	}
}
func (un *unifiedNotifier) processAlert(alert AlertRequest) {
	var wg sync.WaitGroup
	if un.tgEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			un.sendTelegram(alert)
		}()
	}
	wg.Wait()
}
func (un *unifiedNotifier) sendTelegram(alert AlertRequest) {
	if !un.tgEnabled || un.tgBot == nil {
		return
	}

	msgText := fmt.Sprintf("<b>%s</b>\n%s", alert.Subject, alert.Message)
	if len(alert.Metadata) > 0 {
		msgText += "\n\n<b>Metadata:</b>"
		for k, v := range alert.Metadata {
			msgText += fmt.Sprintf("\n%s: %v", k, v)
		}
	}

	chatID, err := strconv.ParseInt(un.telegramConfig.ChatID, 10, 64)
	if err != nil {
		Logger.Error("Invalid Telegram ChatID", zap.String("chat_id", un.telegramConfig.ChatID), zap.Error(err))
		return
	}
	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ParseMode = "HTML"

	if _, err := un.tgBot.Send(msg); err != nil {
		Logger.Error("Failed to send Telegram alert",
			zap.String("subject", alert.Subject),
			zap.Error(err))
	} else {
		Logger.Debug("Telegram alert sent successfully",
			zap.String("subject", alert.Subject))
	}
}
func safeTokenPrefix(token string) string {
	if len(token) > 5 {
		return token[:3] + "..."
	}
	return "[redacted]"
}
