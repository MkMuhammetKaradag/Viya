package database

import (
	"context"
	"fmt"
	"trip-service/internal/domain"

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
func (r *Repository) AddWaypointPhotoWithTags(ctx context.Context, wpID uuid.UUID, photoURL string, tags []domain.Tag) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Önce fotoğrafı ekle ve ID'sini al
	var photoID uuid.UUID
	err = tx.QueryRowContext(ctx, "INSERT INTO photos (waypoint_id, url) VALUES ($1, $2) RETURNING id", wpID, photoURL).Scan(&photoID)
	if err != nil {
		return err
	}

	// 2. Eğer etiket varsa, senin CreateTrip'teki bulk insert mantığını buraya yapıştır
	if len(tags) > 0 {
		query := "INSERT INTO photo_tags (photo_id, label, x_pos, y_pos) VALUES "
		var vals []interface{}
		for i, t := range tags {
			pIdx := i * 4
			query += fmt.Sprintf("($%d, $%d, $%d, $%d),", pIdx+1, pIdx+2, pIdx+3, pIdx+4)
			vals = append(vals, photoID, t.Label, t.XPos, t.YPos)
		}
		query = query[:len(query)-1]
		_, err = tx.ExecContext(ctx, query, vals...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
