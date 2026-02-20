package database

import (
	"context"
	"fmt"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) AddWaypoint(ctx context.Context, wp *domain.Waypoint) (uuid.UUID, error) {
	// 1. İşlemi Başlat
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
	}

	// Herhangi bir hata durumunda tüm işlemleri geri al (Rollback)
	// Eğer işlem başarılı olursa Commit'ten sonra bu satır bir şey yapmaz.
	defer tx.Rollback()

	// 2. Durağı Ekle
	wpQuery := `INSERT INTO waypoints (trip_id, title,latitude, longitude, description,order_index) VALUES ($1, $2, $3, $4, $5,(SELECT COALESCE(MAX(order_index), 0) + 1 FROM waypoints WHERE trip_id = $1)) RETURNING id`
	var wpID uuid.UUID
	err = tx.QueryRowContext(ctx, wpQuery, wp.TripID, wp.Title, wp.Latitude, wp.Longitude, wp.Description).Scan(&wpID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("waypoint insert failed: %w", err)
	}

	// 3. Her şey tamamsa Değişiklikleri Onayla
	if err := tx.Commit(); err != nil {
		return uuid.Nil, fmt.Errorf("transaction onaylanamadı (commit): %w", err)
	}

	return wpID, nil

}
