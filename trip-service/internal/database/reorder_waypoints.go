package database

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) ReorderWaypoints(ctx context.Context, wpID uuid.UUID, index int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//1.Asım: sıralanacak durağın sırasını öğren
	var currentIndex int
	var tripID uuid.UUID
	err = tx.QueryRowContext(ctx, "SELECT order_index,trip_id FROM waypoints WHERE id  = $1", wpID).Scan(&currentIndex, &tripID)
	if err != nil {
		return err
	}

	if currentIndex == index {
		return nil
	}
	//2.Adım: Sırayı düzeltme (shift)
	if currentIndex < index {
		_, err = tx.ExecContext(ctx, "UPDATE waypoints SET order_index = order_index - 1 WHERE trip_id=$1 and  order_index > $2 and order_index <= $3", tripID, currentIndex, index)
		if err != nil {
			return err
		}
	} else {
		_, err = tx.ExecContext(ctx, "UPDATE waypoints SET order_index = order_index + 1 WHERE trip_id=$1 and  order_index >= $2 and order_index < $3", tripID, index, currentIndex)
		if err != nil {
			return err
		}
	}

	//3.Adım: Sırayı güncelleme
	_, err = tx.ExecContext(ctx, "UPDATE waypoints SET order_index = $1 WHERE id = $2", index, wpID)
	if err != nil {
		return err
	}

	return tx.Commit()

}
