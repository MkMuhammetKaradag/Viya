package database

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) UpsertLocalFollow(ctx context.Context, followerID, followingID uuid.UUID, status string) error {
	query := `
        INSERT INTO local_follows (follower_id, following_id, status)
        VALUES ($1, $2, $3)
        ON CONFLICT (follower_id, following_id) 
        DO UPDATE SET status = EXCLUDED.status;`

	_, err := r.db.ExecContext(ctx, query, followerID, followingID, status)
	return err
}
