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

	query := `INSERT INTO trips (user_id, title, description, cover_image_url, is_public, published_at, is_active)
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var tripID uuid.UUID
	err = tx.QueryRowContext(ctx, query, trip.UserID, trip.Title, trip.Description,
		trip.CoverImageURL, trip.IsPublic, trip.PublishedAt, trip.IsActive).Scan(&tripID)
	if err != nil {
		return uuid.Nil, err
	}

	// 2. YENİ: Trip Kategorilerini İşle (Many-to-Many)
	if len(trip.CategoryIDs) > 0 {
		catQuery := "INSERT INTO trip_categories (trip_id, category_id) VALUES "
		var catVals []interface{}
		for i, catID := range trip.CategoryIDs {
			pIdx := i * 2
			catQuery += fmt.Sprintf("($%d, $%d),", pIdx+1, pIdx+2)
			catVals = append(catVals, tripID, catID)
		}
		catQuery = catQuery[:len(catQuery)-1]
		_, err = tx.ExecContext(ctx, catQuery, catVals...)
		if err != nil {
			return uuid.Nil, fmt.Errorf("trip categories insert failed: %w", err)
		}
	}

	// 3. Waypoints Döngüsü
	for _, wp := range trip.Waypoints {
		var wpID uuid.UUID

		wpQuery := `INSERT INTO waypoints (trip_id, category_id, title, description, order_index, latitude, longitude, note)
                    VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
		err = tx.QueryRowContext(ctx, wpQuery, tripID, wp.CategoryID, wp.Title, wp.Description,
			wp.OrderIndex, wp.Latitude, wp.Longitude, wp.Note).Scan(&wpID)
		if err != nil {
			return uuid.Nil, err
		}

		// 4. Fotoğrafları ve Etiketleri İşle (Aynı kalıyor...)
		for _, photo := range wp.Photos {
			var photoID uuid.UUID
			photoQuery := `INSERT INTO photos (waypoint_id, url) VALUES ($1, $2) RETURNING id`
			err = tx.QueryRowContext(ctx, photoQuery, wpID, photo.URL).Scan(&photoID)
			if err != nil {
				return uuid.Nil, err
			}

			if len(photo.Tags) > 0 {
				tagQuery := "INSERT INTO photo_tags (photo_id, label, x_pos, y_pos) VALUES "
				var tagVals []interface{}
				for i, tag := range photo.Tags {
					pIdx := i * 4
					tagQuery += fmt.Sprintf("($%d, $%d, $%d, $%d),", pIdx+1, pIdx+2, pIdx+3, pIdx+4)
					tagVals = append(tagVals, photoID, tag.Label, tag.XPos, tag.YPos)
				}
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
