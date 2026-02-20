package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) DeleteWaypoint(ctx context.Context, waypointID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//1.Asım: silinecek durağın sırasını öğren
	var orderIndex int
	var tripID uuid.UUID

	err = tx.QueryRowContext(ctx, "SELECT order_index,trip_id FROM waypoints WHERE id  = $1", waypointID).Scan(&orderIndex, &tripID)
	if err != nil {
		return fmt.Errorf("waypoint order index fetch failed: %w", err)
	}
	//2.Adım: Durağa silme
	_, err = tx.ExecContext(ctx, "DELETE FROM waypoints WHERE id = $1", waypointID)
	if err != nil {
		return fmt.Errorf("waypoint deletion failed: %w", err)
	}

	//3.Adım: Sırayı düzeltme (shift)
	_, err = tx.ExecContext(ctx, "UPDATE waypoints SET order_index = order_index - 1 WHERE trip_id=$1 and  order_index > $2", tripID, orderIndex)
	if err != nil {
		return fmt.Errorf("order index shift failed: %w", err)
	}

	return tx.Commit()

}
