package domain

import (
	"context"
	"time"
)

type SessionData struct {
	UserID string `json:"user_id"`
	// Add more fields as needed, e.g., roles, permissions, etc.
}
type SessionRepository interface {
	CreateSession(ctx context.Context, duration time.Duration, data *SessionData) (string, error)
}
