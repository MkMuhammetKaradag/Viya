package session

import (
	"auth-service/internal/config"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(cfg *config.Config) (*SessionRepository, error) {
	client, err := newRedisDB(cfg)
	if err != nil {
		return nil, err
	}

	return &SessionRepository{
		client: client,
	}, nil
}
func (sr *SessionRepository) userSessionsKey(userID string) string {
	return fmt.Sprintf("auth-service:user_sessions:%s", userID)
}
func (sr *SessionRepository) Close() error {
	if sr.client != nil {
		err := sr.client.Close()
		if err != nil {
			return fmt.Errorf("redis connection close error: %w", err)
		}
		fmt.Println("📡 Redis connection closed successfully.")
	}
	return nil
}
