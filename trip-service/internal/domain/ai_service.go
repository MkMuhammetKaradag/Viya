package domain

import "context"

type AIService interface {
	GetVector(ctx context.Context, text string) ([]float32, error)
}
