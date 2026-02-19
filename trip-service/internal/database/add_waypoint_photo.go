package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) AddWaypointPhotos(ctx context.Context, waypointID uuid.UUID, photoURLs []string) error {
	// 1. İşlemi Başlat
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 2. Fotoğrafları döngü ile ekle

	query := `INSERT INTO photos (waypoint_id, url) VALUES ($1, $2)`

	for _, url := range photoURLs {
		_, err := tx.ExecContext(ctx, query, waypointID, url)
		if err != nil {
			return fmt.Errorf("photo insert failed for url %s: %w", url, err)
		}
	}

	// 3. Değişiklikleri Onayla
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}
