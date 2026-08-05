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
		if cfg.AnthropicKey == "" && cfg.GeminiKey == "" && cfg.OpenRouterKey == "" {
			initErr = errors.New("at least one AI provider key is required (OpenRouter, Anthropic or Gemini)")
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
			zap.String("default_provider", cfg.DefaultProvider),
			zap.Bool("openrouter_enabled", cfg.OpenRouterKey != ""),
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
// provider is a named LLM backend in the fallback chain.
type provider struct {
	name string
	key  string
	fn   func(context.Context, []Message, string) (*AIResponse, error)
}

// providerChain returns the enabled providers ordered so DefaultProvider is
// tried first, then the rest. Only providers with a configured key are included.
// This makes the client provider-agnostic (roadmap goal): swapping the primary
// is a config change, not a code change.
func (ai *AIClient) providerChain() []provider {
	all := []provider{
		{name: "openrouter", key: ai.config.OpenRouterKey, fn: ai.chatWithOpenRouter},
		{name: "anthropic", key: ai.config.AnthropicKey, fn: ai.chatWithAnthropic},
		{name: "gemini", key: ai.config.GeminiKey, fn: ai.chatWithGemini},
	}
	// Normalize DefaultProvider aliases to a chain name.
	preferred := ""
	switch ai.config.DefaultProvider {
	case "openrouter":
		preferred = "openrouter"
	case "claude", "anthropic":
		preferred = "anthropic"
	case "google", "gemini":
		preferred = "gemini"
	}

	chain := make([]provider, 0, len(all))
	// Preferred first (if enabled).
	for _, p := range all {
		if p.name == preferred && p.key != "" {
			chain = append(chain, p)
		}
	}
	// Then the remaining enabled providers in default order.
	for _, p := range all {
		if p.name != preferred && p.key != "" {
			chain = append(chain, p)
		}
	}
	return chain
}

func (ai *AIClient) Chat(ctx context.Context, messages []Message, systemPrompt string) (*AIResponse, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 30*time.Second)
	defer cancel()

	chain := ai.providerChain()
	if len(chain) == 0 {
		return nil, errors.New("no AI provider available")
	}

	var lastErr error
	for _, p := range chain {
		resp, err := p.fn(ctx, messages, systemPrompt)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		ai.logger.Warn("AI provider failed, trying next", zap.String("provider", p.name), zap.Error(err))
	}
	return nil, fmt.Errorf("all AI providers failed: %w", lastErr)
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
