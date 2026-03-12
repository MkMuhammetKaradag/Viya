package session

import (
	"auth-service/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (sr *SessionRepository) CreateSession(ctx context.Context, duration time.Duration, data *domain.SessionData) (string, error) {
	// Implement session creation logic, e.g., generate a session token and store it in Redis
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	pipe := sr.client.TxPipeline()
	sessionID := generateSessionID()
	pipe.Set(ctx, "session:"+sessionID, jsonData, duration)
	pipe.SAdd(ctx, sr.userSessionsKey(data.UserID), sessionID)
	pipe.Expire(ctx, sr.userSessionsKey(data.UserID), duration)
	if _, err := pipe.Exec(ctx); err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return sessionID, nil
}
func generateSessionID() string {
	//return "sess:" + uuid.New().String()
	return uuid.New().String()
}
