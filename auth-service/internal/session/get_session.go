package session

import (
	"auth-service/internal/domain"
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func (sr *SessionRepository) GetSession(ctx context.Context, sessionID string) (*domain.SessionData, error) {
	getSessionID := "session:" + sessionID
	val, err := sr.client.Get(ctx, getSessionID).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("session data not finde")
	}
	if err != nil {
		return nil, err
	}
	var session domain.SessionData

	err = json.Unmarshal([]byte(val), &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}
