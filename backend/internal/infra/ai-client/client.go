package aiclient

import (
	"context"
	"errors"
	"fmt"
	"github/socialforge/config"
	"github/socialforge/internal/infra/contextpool"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

var (
	AIClientStorage *AIClient
	aiOnce          sync.Once
)

const (
	AnthropicAPIURL = "https://api.anthropic.com/v1/messages"
	GeminiAPIURL    = "https://generativelanguage.googleapis.com/v1beta/models"
)

type AIClient struct {
	config     *config.AIConfig
	logger     *zap.Logger
	httpClient *http.Client
	isUp       bool
	mu         sync.RWMutex
}
type Message struct {
	Role    string `json:"role"`    // "user" or "assistant"
	Content string `json:"content"` // Message content
}
type AIResponse struct {
	Content   string `json:"content"`
	Model     string `json:"model"`
	TokensIn  int    `json:"tokens_in"`
	TokensOut int    `json:"tokens_out"`
}

func NewAIClient(cfg *config.AIConfig, logger *zap.Logger) (*AIClient, error) {
	var initErr error
	aiOnce.Do(func() {
		if cfg.AnthropicKey == "" && cfg.GeminiKey == "" {
			initErr = errors.New("at least one AI provider key is required (Anthropic or Gemini)")
			logger.Error("AI configuration missing: no API keys provided")
			return
		}

		AIClientStorage = &AIClient{
			config: cfg,
			logger: logger,
			httpClient: &http.Client{
				Timeout: 60 * time.Second,
			},
			isUp: true,
		}

		logger.Info("✅ AI client initialized successfully",
			zap.Bool("anthropic_enabled", cfg.AnthropicKey != ""),
			zap.Bool("gemini_enabled", cfg.GeminiKey != ""),
		)
	})

	if initErr != nil {
		return nil, initErr
	}
	return AIClientStorage, nil
}
func GetAIClient() (*AIClient, error) {
	if AIClientStorage == nil {
		return nil, errors.New("AI client not initialized: call NewAIClient first")
	}
	return AIClientStorage, nil
}
func (ai *AIClient) IsUp() bool {
	ai.mu.RLock()
	defer ai.mu.RUnlock()
	return ai.isUp
}
func (ai *AIClient) Chat(ctx context.Context, messages []Message, systemPrompt string) (*AIResponse, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 30*time.Second)
	defer cancel()

	// Try Anthropic Claude first (primary)
	if ai.config.AnthropicKey != "" {
		response, err := ai.chatWithAnthropic(ctx, messages, systemPrompt)
		if err == nil {
			return response, nil
		}
		ai.logger.Warn("Anthropic request failed, falling back to Gemini",
			zap.Error(err),
		)
	}

	// Fallback to Gemini
	if ai.config.GeminiKey != "" {
		response, err := ai.chatWithGemini(ctx, messages, systemPrompt)
		if err == nil {
			return response, nil
		}
		ai.logger.Error("Gemini request also failed",
			zap.Error(err),
		)
		return nil, fmt.Errorf("all AI providers failed: %w", err)
	}

	return nil, errors.New("no AI provider available")
}
func (ai *AIClient) GenerateAutoReply(ctx context.Context, customerMessage string, context string) (string, error) {
	systemPrompt := `You are a helpful customer service assistant. Generate a professional and friendly response to the customer's message. 
Keep the response concise (2-3 sentences max) and helpful. Use the provided context to personalize the response if available.`

	messages := []Message{
		{
			Role:    "user",
			Content: fmt.Sprintf("Customer message: %s\n\nContext: %s\n\nGenerate a helpful response:", customerMessage, context),
		},
	}

	response, err := ai.Chat(ctx, messages, systemPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate auto reply: %w", err)
	}

	return response.Content, nil
}
func (ai *AIClient) AnalyzeSentiment(ctx context.Context, message string) (string, error) {
	systemPrompt := `Analyze the sentiment of the following message and respond with only one word: "positive", "negative", or "neutral".`

	messages := []Message{
		{
			Role:    "user",
			Content: message,
		},
	}

	response, err := ai.Chat(ctx, messages, systemPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to analyze sentiment: %w", err)
	}

	return response.Content, nil
}
func (ai *AIClient) SummarizeConversation(ctx context.Context, conversationHistory []Message) (string, error) {
	systemPrompt := `Summarize the following conversation in 2-3 sentences. Focus on the main topics and outcomes.`

	// Convert conversation to a single message
	var conversationText string
	for _, msg := range conversationHistory {
		conversationText += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
	}

	messages := []Message{
		{
			Role:    "user",
			Content: conversationText,
		},
	}

	response, err := ai.Chat(ctx, messages, systemPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to summarize conversation: %w", err)
	}

	return response.Content, nil
}
func (ai *AIClient) Close() error {
	ai.mu.Lock()
	defer ai.mu.Unlock()

	ai.isUp = false
	ai.logger.Info("✅ AI client closed")
	return nil
}
