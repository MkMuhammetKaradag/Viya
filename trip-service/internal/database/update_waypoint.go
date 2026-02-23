package database

import (
	"context"
	"fmt"
	"trip-service/internal/domain"
)

func (r *Repository) UpdateWaypoint(ctx context.Context, wp *domain.Waypoint) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `
        UPDATE waypoints 
        SET title = $1, 
            description = $2, 
            latitude = $3, 
            longitude = $4
        WHERE id = $5`
	result, err := tx.ExecContext(ctx, query,
		wp.Title,
		wp.Description,
		wp.Latitude,
		wp.Longitude,
		wp.ID,
	)
	if err != nil {
		return fmt.Errorf("waypoint update failed: %w", err)
	}

	// Gerçekten bir satır güncellendi mi kontrolü (Opsiyonel ama güvenli)
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("waypoint not found")
	}

	return tx.Commit()

}
