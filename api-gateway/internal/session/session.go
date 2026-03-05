package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Session struct {
	UserID string `json:"user_id"`
}

type SessionManager struct {
	client *redis.Client

	ttl time.Duration
}

func NewSessionManager(redisAddr string, password string, db int, ttl time.Duration) (*SessionManager, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis bağlantısı başarısız: %w", err)
	}

	return &SessionManager{
		client: client,
		ttl:    ttl,
	}, nil
}
func (c *SessionManager) GetSession(ctx context.Context, token string) (*Session, error) {
	// key := c.sessionKey(token)

	data, err := c.client.Get(ctx, token).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("seesion not found")
	}
	if err != nil {
		return nil, fmt.Errorf("session read error : %w", err)
	}

	var session Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("session deserialize error : %w", err)
	}

	return &session, nil
}
func (c *SessionManager) sessionKey(token string) string {
	return fmt.Sprintf("session:%s", token)
}
func (c *SessionManager) Close() error {
	return c.client.Close()
}
