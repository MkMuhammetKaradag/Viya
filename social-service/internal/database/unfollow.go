package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) UnfollowUser(ctx context.Context, followerID, targetUserID uuid.UUID) error {
	query := `
        DELETE FROM follows 
        WHERE follower_id = $1 AND following_id = $2
    `
	result, err := r.db.ExecContext(ctx, query, followerID, targetUserID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("No follow-up relationship found.")
	}

	return nil
}
