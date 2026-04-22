package database

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) UpsertLocalBlock(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, "INSERT INTO local_blocks (blocker_id, blocked_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", blockerID, blockedID)
	if err != nil {
		return err
	}

	deleteQuery := `
        DELETE FROM local_follows 
        WHERE (follower_id = $1 AND following_id = $2) 
           OR (follower_id = $2 AND following_id = $1)
    `
	_, err = tx.ExecContext(ctx, deleteQuery, blockerID, blockedID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
