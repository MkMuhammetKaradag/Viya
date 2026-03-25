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

	// 1. Mevcut en yüksek index'i öğrenelim
	var currentMax int
	maxQuery := `SELECT COALESCE(MAX(order_index), 0) FROM waypoints WHERE trip_id = $1`
	err = tx.QueryRowContext(ctx, maxQuery, wp.TripID).Scan(&currentMax)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not get max order: %w", err)
	}

	var finalOrderIndex int

	// 2. MANTIK KONTROLÜ: Kullanıcının istediği index mantıklı mı?
	// Eğer wp.OrderIndex 0'dan büyükse ve mevcut MAX+1'den çok daha büyükse,
	// onu en sona (MAX+1) zorluyoruz (Boşluk oluşmasın diye).
	if wp.OrderIndex > 0 {
		if wp.OrderIndex > currentMax+1 {
			// Kullanıcı uçuk bir rakam gönderdi, biz onu "en son + 1" yapıyoruz.
			finalOrderIndex = currentMax + 1
		} else {
			// Araya ekleme: Mevcutları birer kaydır (Shift)
			shiftQuery := `
                UPDATE waypoints 
                SET order_index = order_index + 1 
                WHERE trip_id = $1 AND order_index >= $2`
			_, err = tx.ExecContext(ctx, shiftQuery, wp.TripID, wp.OrderIndex)
			if err != nil {
				return uuid.Nil, fmt.Errorf("shifting failed: %w", err)
			}
			finalOrderIndex = wp.OrderIndex
		}
	} else {
		// Kullanıcı 0 gönderdi veya hiç göndermedi: Doğrudan sona ekle
		finalOrderIndex = currentMax + 1
	}

	// 3. ADIM: Kaydı Gerçekleştir
	insertQuery := `
        INSERT INTO waypoints (trip_id, title, latitude, longitude, description, order_index, note) 
        VALUES ($1, $2, $3, $4, $5, $6, $7) 
        RETURNING id`

	var wpID uuid.UUID
	err = tx.QueryRowContext(ctx, insertQuery,
		wp.TripID, wp.Title, wp.Latitude, wp.Longitude, wp.Description, finalOrderIndex, wp.Note,
	).Scan(&wpID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("insert failed: %w", err)
	}

	return wpID, tx.Commit()
}
