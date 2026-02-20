package database

import (
	"context"
	"fmt"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) AddWaypoint(ctx context.Context, wp *domain.Waypoint) (uuid.UUID, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()

	// 1. ADIM: Sıralamayı Belirle
	// Eğer wp.OrderIndex 0 ise (gelmemişse), en sona ekle.
	// Eğer 0'dan büyükse, araya yer aç (Shift).

	var finalOrderIndex int

	if wp.OrderIndex > 0 {
		// Araya ekleme: Mevcutları birer kaydır
		shiftQuery := `
            UPDATE waypoints 
            SET order_index = order_index + 1 
            WHERE trip_id = $1 AND order_index >= $2`
		_, err = tx.ExecContext(ctx, shiftQuery, wp.TripID, wp.OrderIndex)
		if err != nil {
			return uuid.Nil, fmt.Errorf("shifting failed: %w", err)
		}
		finalOrderIndex = wp.OrderIndex
	} else {
		// En sona ekleme: Mevcut MAX + 1 değerini hesapla
		calcQuery := `SELECT COALESCE(MAX(order_index), 0) + 1 FROM waypoints WHERE trip_id = $1`
		err = tx.QueryRowContext(ctx, calcQuery, wp.TripID).Scan(&finalOrderIndex)
		if err != nil {
			return uuid.Nil, fmt.Errorf("calc order failed: %w", err)
		}
	}

	// 2. ADIM: Kaydı Gerçekleştir
	insertQuery := `
        INSERT INTO waypoints (trip_id, title, latitude, longitude, description, order_index) 
        VALUES ($1, $2, $3, $4, $5, $6) 
        RETURNING id`

	var wpID uuid.UUID
	err = tx.QueryRowContext(ctx, insertQuery,
		wp.TripID, wp.Title, wp.Latitude, wp.Longitude, wp.Description, finalOrderIndex,
	).Scan(&wpID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("insert failed: %w", err)
	}

	return wpID, tx.Commit()
}
