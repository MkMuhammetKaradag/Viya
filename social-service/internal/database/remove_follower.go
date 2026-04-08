package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) RemoveFollower(ctx context.Context, myID, followerID uuid.UUID) error {
	query := `
        DELETE FROM follows 
        WHERE follower_id = $1 AND following_id = $2
    `
	// followerID: Listeden atılacak kişi
	// myID: Sen (Takip edilen kişi)
	result, err := r.db.ExecContext(ctx, query, followerID, myID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("takipçi bulunamadı")
	}

	return nil
}
