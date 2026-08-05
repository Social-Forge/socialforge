package aiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	// OpenAIEmbeddingsURL is OpenAI's direct embeddings endpoint.
	OpenAIEmbeddingsURL = "https://api.openai.com/v1/embeddings"
	// OpenRouterEmbeddingsURL is OpenRouter's OpenAI-compatible embeddings endpoint.
	OpenRouterEmbeddingsURL = "https://openrouter.ai/api/v1/embeddings"
	// EmbeddingDim must match ai_knowledge.embedding_vec (vector(1536)) —
	// text-embedding-3-small emits 1536 dims.
	EmbeddingDim = 1536
)

type embeddingsRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type embeddingsResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// EmbeddingsEnabled reports whether embeddings can be generated. Prefers
// OpenRouter (same account/credit as chat), falls back to a direct OpenAI key.
func (ai *AIClient) EmbeddingsEnabled() bool {
	return ai.config.OpenRouterKey != "" || ai.config.OpenAIKey != ""
}

// Embed returns the embedding vector for a piece of text. It routes through
// OpenRouter when available (which has the working credit), else OpenAI direct.
// Callers should treat any error as "fall back to lexical search".
func (ai *AIClient) Embed(ctx context.Context, text string) ([]float32, error) {
	model := ai.config.EmbeddingModel
	if model == "" {
		model = "text-embedding-3-small"
	}

	var url, key string
	switch {
	case ai.config.OpenRouterKey != "":
		url, key = OpenRouterEmbeddingsURL, ai.config.OpenRouterKey
		// OpenRouter needs a provider-prefixed slug.
		if !strings.Contains(model, "/") {
			model = "openai/" + model
		}
	case ai.config.OpenAIKey != "":
		url, key = OpenAIEmbeddingsURL, ai.config.OpenAIKey
		// OpenAI direct wants the bare model id.
		model = strings.TrimPrefix(model, "openai/")
	default:
		return nil, errors.New("embeddings unavailable: no OpenRouter/OpenAI key")
	}

	jsonData, err := json.Marshal(embeddingsRequest{Model: model, Input: text})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal embeddings request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create embeddings request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)

	resp, err := ai.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embeddings request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read embeddings response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embeddings API error (status %d): %s", resp.StatusCode, string(body))
	}

	var er embeddingsResponse
	if err := json.Unmarshal(body, &er); err != nil {
		return nil, fmt.Errorf("failed to unmarshal embeddings response: %w", err)
	}
	if er.Error != nil && er.Error.Message != "" {
		return nil, fmt.Errorf("embeddings error: %s", er.Error.Message)
	}
	if len(er.Data) == 0 || len(er.Data[0].Embedding) == 0 {
		return nil, errors.New("empty embedding from provider")
	}
	return er.Data[0].Embedding, nil
}
