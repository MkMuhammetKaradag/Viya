package session

import (
	"context"
	"fmt"
	"time"
)

// IsActionLocked, belirtilen anahtarın Redis'te olup olmadığını kontrol eder.
func (sr *SessionRepository) IsActionLocked(ctx context.Context, key string) (bool, error) {
	exists, err := sr.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check action lock: %w", err)
	}

	return exists > 0, nil
}

// SetActionLock, belirli bir anahtarı verilen süre boyunca kilitli tutar.
func (sr *SessionRepository) SetActionLock(ctx context.Context, key string, duration time.Duration) error {

	err := sr.client.Set(ctx, key, "1", duration).Err()
	if err != nil {
		return fmt.Errorf("failed to set action lock: %w", err)
	}
	return nil
}
