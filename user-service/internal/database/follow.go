package database

import (
	"context"
	"fmt"

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
func (r *Repository) DeleteLocalFollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	query := `
        DELETE FROM local_follows 
        WHERE follower_id = $1 AND following_id = $2;`

	result, err := r.db.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("No follow-up relationship found.")
	}

	return nil
}
func (r *Repository) CheckFollowStatus(ctx context.Context, followerID, followedID uuid.UUID) (bool, error) {

	if followerID == followedID {
		return true, nil
	}

	query := `
        SELECT EXISTS (
            SELECT 1 FROM local_follows 
            WHERE follower_id = $1 AND following_id = $2 and status = 'ACCEPTED'	
        )`

	var isFollowing bool
	err := r.db.QueryRowContext(ctx, query, followerID, followedID).Scan(&isFollowing)
	if err != nil {

		return false, fmt.Errorf("takip durumu sorgulanırken hata: %w", err)
	}

	return isFollowing, nil
}
