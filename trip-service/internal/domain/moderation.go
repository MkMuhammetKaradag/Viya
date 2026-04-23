package domain

import "context"

type ModerationResult struct {
	IsAppropriate bool
	Reason        string
}

type ModerationService interface {
	Moderate(ctx context.Context, content string) (*ModerationResult, error)
}
