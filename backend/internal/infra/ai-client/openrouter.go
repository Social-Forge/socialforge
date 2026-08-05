package aiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// OpenRouterAPIURL is the OpenAI-compatible chat-completions endpoint.
const OpenRouterAPIURL = "https://openrouter.ai/api/v1/chat/completions"

type openRouterRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type openRouterResponse struct {
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// chatWithOpenRouter calls OpenRouter's OpenAI-compatible API. The system prompt
// is sent as a leading `system` message (OpenAI convention).
func (ai *AIClient) chatWithOpenRouter(ctx context.Context, messages []Message, systemPrompt string) (*AIResponse, error) {
	all := make([]Message, 0, len(messages)+1)
	if systemPrompt != "" {
		all = append(all, Message{Role: "system", Content: systemPrompt})
	}
	all = append(all, messages...)

	model := ai.config.OpenRouterModel
	if model == "" {
		model = "anthropic/claude-3.5-sonnet"
	}
	jsonData, err := json.Marshal(openRouterRequest{Model: model, Messages: all})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", OpenRouterAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ai.config.OpenRouterKey)
	// Optional attribution headers recommended by OpenRouter.
	req.Header.Set("HTTP-Referer", "https://socialforge.io")
	req.Header.Set("X-Title", "Social Forge")

	resp, err := ai.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var orResp openRouterResponse
	if err := json.Unmarshal(body, &orResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if orResp.Error != nil && orResp.Error.Message != "" {
		return nil, fmt.Errorf("openrouter error: %s", orResp.Error.Message)
	}
	if len(orResp.Choices) == 0 {
		return nil, errors.New("empty response from OpenRouter")
	}

	model = orResp.Model
	if model == "" {
		model = ai.config.OpenRouterModel
	}
	return &AIResponse{
		Content:   orResp.Choices[0].Message.Content,
		Model:     model,
		TokensIn:  orResp.Usage.PromptTokens,
		TokensOut: orResp.Usage.CompletionTokens,
	}, nil
}
