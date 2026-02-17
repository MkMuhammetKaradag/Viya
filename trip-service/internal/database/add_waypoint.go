package database

import (
	"context"
	"fmt"
	"trip-service/internal/domain"
)

func (r *Repository) AddWaypoint(ctx context.Context, wp *domain.Waypoint, photos []string) error {
	// 1. İşlemi Başlat
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Herhangi bir hata durumunda tüm işlemleri geri al (Rollback)
	// Eğer işlem başarılı olursa Commit'ten sonra bu satır bir şey yapmaz.
	defer tx.Rollback()

	// 2. Durağı Ekle
	wpQuery := `INSERT INTO waypoints (trip_id, lat, lon, note) VALUES ($1, $2, $3, $4) RETURNING id`
	var wpID string
	err = tx.QueryRowContext(ctx, wpQuery, wp.TripID, wp.Lat, wp.Lon, wp.Note).Scan(&wpID)
	if err != nil {
		return fmt.Errorf("waypoint insert failed: %w", err)
	}

	// 3. Fotoğrafları Ekle (Döngü ile)
	if len(photos) > 0 {
		photoQuery := `INSERT INTO photos (waypoint_id, url) VALUES ($1, $2)`
		for _, url := range photos {
			_, err = tx.ExecContext(ctx, photoQuery, wpID, url)
			if err != nil {
				return fmt.Errorf("fotoğraf eklenirken hata: %w", err)
			}
		}
	}

	// 4. Her şey tamamsa Değişiklikleri Onayla
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction onaylanamadı (commit): %w", err)
	}

	return nil

}
