package database

import (
	"context"
	"fmt"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

// func (r *Repository) CreateTrip(ctx context.Context, trip *domain.Trip) (uuid.UUID, error) {
// 	// 1. Transaction başlat (Çünkü içinde noktalar olabilir, veri bütünlüğü şart)
// 	tx, err := r.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return uuid.Nil, err
// 	}
// 	defer tx.Rollback()

// 	// 2. Ana Trip kaydını yap
// 	query := `
// 		INSERT INTO trips (user_id, title, description, cover_image_url, is_public, published_at, is_active)
// 		VALUES ($1, $2, $3, $4, $5, $6, $7)
// 		RETURNING id
// 	`
// 	var tripID uuid.UUID
// 	err = tx.QueryRowContext(ctx, query,
// 		trip.UserID, trip.Title, trip.Description,
// 		trip.CoverImageURL, trip.IsPublic, trip.PublishedAt, trip.IsActive,
// 	).Scan(&tripID)

// 	if err != nil {
// 		return uuid.Nil, err
// 	}

// 	// 3. HİBRİT KONTROL: Eğer içinde Waypoints varsa onları da işle
// 	if len(trip.Waypoints) > 0 {
// 		for _, wp := range trip.Waypoints {
// 			var wpID uuid.UUID
// 			wpQuery := `
// 				INSERT INTO waypoints (trip_id, title, description, order_index, latitude, longitude, note)
// 				VALUES ($1, $2, $3, $4, $5, $6, $7)
// 				RETURNING id
// 			`
// 			err = tx.QueryRowContext(ctx, wpQuery,
// 				tripID, wp.Title, wp.Description, wp.OrderIndex,
// 				wp.Latitude, wp.Longitude, wp.Note,
// 			).Scan(&wpID)

// 			if err != nil {
// 				return uuid.Nil, err
// 			}

// 			// 4. Noktanın fotoğrafları varsa onları da işle
// 			for _, photo := range wp.Photos {
// 				_, err = tx.ExecContext(ctx, "INSERT INTO photos (waypoint_id, url) VALUES ($1, $2)", wpID, photo.URL)
// 				if err != nil {
// 					return uuid.Nil, err
// 				}
// 			}
// 		}
// 	}

// 	// 5. Her şey tamamsa Commit et
// 	if err := tx.Commit(); err != nil {
// 		return uuid.Nil, err
// 	}

//		return tripID, nil
//	}
func (r *Repository) CreateTrip(ctx context.Context, trip *domain.Trip) (uuid.UUID, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()

	// 1. Ana Trip Kaydı (Aynı kalıyor)
	query := `INSERT INTO trips (user_id, title, description, cover_image_url, is_public, published_at, is_active)
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var tripID uuid.UUID
	err = tx.QueryRowContext(ctx, query, trip.UserID, trip.Title, trip.Description,
		trip.CoverImageURL, trip.IsPublic, trip.PublishedAt, trip.IsActive).Scan(&tripID)
	if err != nil {
		return uuid.Nil, err
	}

	// 2. Waypoint ve Photo'ları optimize edelim
	if len(trip.Waypoints) > 0 {
		for _, wp := range trip.Waypoints {
			var wpID uuid.UUID

			wpQuery := `INSERT INTO waypoints (trip_id, title, description, order_index, latitude, longitude, note)
                        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
			err = tx.QueryRowContext(ctx, wpQuery, tripID, wp.Title, wp.Description,
				wp.OrderIndex, wp.Latitude, wp.Longitude, wp.Note).Scan(&wpID)
			if err != nil {
				return uuid.Nil, err
			}

			// --- KRİTİK DÜZELTME: Fotoğrafları tek tek değil, tek sorguda ekle ---
			if len(wp.Photos) > 0 {
				photoQuery := "INSERT INTO photos (waypoint_id, url) VALUES "
				var vals []interface{}
				for i, p := range wp.Photos {
					photoQuery += fmt.Sprintf("($%d, $%d),", i*2+1, i*2+2)
					vals = append(vals, wpID, p.URL)
				}
				photoQuery = photoQuery[:len(photoQuery)-1]
				_, err = tx.ExecContext(ctx, photoQuery, vals...)
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
