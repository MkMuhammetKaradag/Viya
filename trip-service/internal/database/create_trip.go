package database

import (
	"context"
	"fmt"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) CreateTrip(ctx context.Context, trip *domain.Trip) (uuid.UUID, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()

	// 1. Ana Trip Kaydı
	query := `INSERT INTO trips (user_id, title, description, cover_image_url, is_public, published_at, is_active)
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var tripID uuid.UUID
	err = tx.QueryRowContext(ctx, query, trip.UserID, trip.Title, trip.Description,
		trip.CoverImageURL, trip.IsPublic, trip.PublishedAt, trip.IsActive).Scan(&tripID)
	if err != nil {
		return uuid.Nil, err
	}

	// 2. Waypoints Döngüsü
	for _, wp := range trip.Waypoints {
		var wpID uuid.UUID
		wpQuery := `INSERT INTO waypoints (trip_id, title, description, order_index, latitude, longitude, note)
                    VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
		err = tx.QueryRowContext(ctx, wpQuery, tripID, wp.Title, wp.Description,
			wp.OrderIndex, wp.Latitude, wp.Longitude, wp.Note).Scan(&wpID)
		if err != nil {
			return uuid.Nil, err
		}

		// 3. Fotoğrafları İşle
		for _, photo := range wp.Photos {
			var photoID uuid.UUID
			// Fotoğrafı ekle ve ID'sini dön (RETURNING id)
			photoQuery := `INSERT INTO photos (waypoint_id, url) VALUES ($1, $2) RETURNING id`
			err = tx.QueryRowContext(ctx, photoQuery, wpID, photo.URL).Scan(&photoID)
			if err != nil {
				return uuid.Nil, err
			}

			// 4. Eğer bu fotoğrafın ETİKETLERİ varsa onları TOPLU ekle
			if len(photo.Tags) > 0 {
				tagQuery := "INSERT INTO photo_tags (photo_id, label, x_pos, y_pos) VALUES "
				var tagVals []interface{}

				for i, tag := range photo.Tags {
					// Her etiket için 4 parametre ($1, $2, $3, $4)
					pIdx := i * 4
					tagQuery += fmt.Sprintf("($%d, $%d, $%d, $%d),", pIdx+1, pIdx+2, pIdx+3, pIdx+4)
					tagVals = append(tagVals, photoID, tag.Label, tag.XPos, tag.YPos)
				}

				// Sondaki virgülü temizle
				tagQuery = tagQuery[:len(tagQuery)-1]

				_, err = tx.ExecContext(ctx, tagQuery, tagVals...)
				if err != nil {
					return uuid.Nil, err
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}
	return tripID, nil
}
