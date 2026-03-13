package domain

import (
	"context"
	"time"
)

type SessionData struct {
	UserID         string    `json:"user_id"`
	FirstCreatedAt time.Time `json:"fisrt_created_at"`
	CreatedAt      time.Time `json:"created_at"`
	Device         string    `json:"device"`
	Ip             string    `json:"ip"`
	// Add more fields as needed, e.g., roles, permissions, etc.
}
type SessionRepository interface {
	CreateSession(ctx context.Context, duration time.Duration, data *SessionData) (string, error)
	IsActionLocked(ctx context.Context, key string) (bool, error)
	SetActionLock(ctx context.Context, key string, duration time.Duration) error
	DeleteSession(ctx context.Context, sessionID string) error
	GetSession(ctx context.Context, sessionID string) (*SessionData, error)
	Rotate(ctx context.Context, oldSessionID string, session *SessionData) (string, error)
	DeleteAllSession(ctx context.Context, userID string) error
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Close() error
}
