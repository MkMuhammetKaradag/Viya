package session

import (
	"context"
	"fmt"
)

func (sr *SessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	session, err := sr.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not finde")
	}
	pipe := sr.client.Pipeline()

	pipe.Del(ctx, sessionID)
	pipe.SRem(ctx, sr.userSessionsKey(session.UserID), sessionID)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
