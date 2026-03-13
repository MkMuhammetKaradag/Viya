package session

import (
	"context"
	"time"
)

func (r *SessionRepository) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}
func (r *SessionRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}
func (r *SessionRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
