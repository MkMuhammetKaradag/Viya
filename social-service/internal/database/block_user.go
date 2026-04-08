package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Engelleme kaydını ekle
	_, err = tx.ExecContext(ctx, "INSERT INTO blocks (blocker_id, blocked_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", blockerID, blockedID)
	if err != nil {
		return err
	}

	// 2. Karşılıklı tüm takipleri sil (Büyük Temizlik)
	deleteQuery := `
        DELETE FROM follows 
        WHERE (follower_id = $1 AND following_id = $2) 
           OR (follower_id = $2 AND following_id = $1)
    `
	_, err = tx.ExecContext(ctx, deleteQuery, blockerID, blockedID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
func (r *Repository) IsBlocked(ctx context.Context, userA, userB uuid.UUID) ([]uuid.UUID, error) {
	query := `
        SELECT blocker_id FROM blocks 
        WHERE (blocker_id = $1 AND blocked_id = $2) 
           OR (blocker_id = $2 AND blocked_id = $1)
    `
	rows, err := r.db.QueryContext(ctx, query, userA, userB)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blockerIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		blockerIDs = append(blockerIDs, id)
	}
	return blockerIDs, nil
}
func (r *Repository) UnblockUser(ctx context.Context, myID, targetUserID uuid.UUID) error {
	query := `
        DELETE FROM blocks 
        WHERE blocker_id = $1 AND blocked_id = $2
    `
	// myID: Engeli kaldıran (Sen)
	// targetUserID: Engeli kaldırılan kişi
	result, err := r.db.ExecContext(ctx, query, myID, targetUserID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("No active blocks were found for this user.")
	}

	return nil
}
