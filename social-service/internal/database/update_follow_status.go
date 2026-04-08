package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) UpdateFollowStatus(ctx context.Context, followerID, followingID uuid.UUID, newStatus string) error {

	query := `
        UPDATE follows 
        SET status = $3 
        WHERE follower_id = $1 AND following_id = $2 AND status = 'PENDING'
    `

	// Eğer 'REJECTED' ise kaydı tamamen silmek de bir tercihtir:
	if newStatus == "REJECTED" {
		query = "DELETE FROM follows WHERE follower_id = $1 AND following_id = $2 AND status = 'PENDING'"
		_, err := r.db.ExecContext(ctx, query, followerID, followingID)
		return err
	}

	result, err := r.db.ExecContext(ctx, query, followerID, followingID, newStatus)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("onaylanacak bekleyen bir istek bulunamadı")
	}

	return nil
}
