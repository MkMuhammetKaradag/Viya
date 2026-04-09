package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"trip-service/internal/domain"
)

type OllamaService struct {
	BaseURL string
	Model   string
	client  *http.Client
}

func NewOllamaService() domain.AIService {
	return &OllamaService{

		BaseURL: "http://localhost:11434/api/embeddings",
		Model:   "nomic-embed-text",
		client: &http.Client{
			Timeout: 10 * time.Second, // 10 saniyeden fazla beklerse kes
		},
	}
}

func (o *OllamaService) GetVector(ctx context.Context, text string) ([]float32, error) {
	payload := map[string]interface{}{
		"model":  o.Model,
		"prompt": text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	// NewRequestWithContext kullanarak context desteği ekledik
	req, err := http.NewRequestWithContext(ctx, "POST", o.BaseURL, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Ollama request create error: %v\n", err)
		return nil, fmt.Errorf("request create error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		fmt.Printf("Ollama request error: %v\n", err)
		return nil, fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Ollama response status: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("ollama error: status %d", resp.StatusCode)
	}

	var res struct {
		Embedding []float32 `json:"embedding"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		fmt.Printf("Ollama decode error: %v\n", err)
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return res.Embedding, nil
}
