package database

import (
	"context"
	"fmt"
	"time"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) GetTripWithWaypointsAndPhotos(ctx context.Context, tripID, currentUserID uuid.UUID) (*domain.Trip, error) {
	query := `
    SELECT 
        t.id, t.user_id, t.title, t.description, t.cover_image_url, t.is_active, t.is_public, t.published_at, t.view_count,t.total_likes, t.total_comments, t.created_at,
        
        EXISTS(SELECT 1 FROM trip_likes WHERE trip_id = t.id AND user_id = $2) as is_liked,
        w.id as wp_id, w.title as wp_title, w.description as wp_desc, w.order_index, w.latitude, w.longitude, w.note, w.created_at as wp_created_at,
        p.id as photo_id, p.url as photo_url,
        pt.id as tag_id, pt.label as tag_label, pt.x_pos as tag_x, pt.y_pos as tag_y
    FROM trips t
    LEFT JOIN waypoints w ON t.id = w.trip_id
    LEFT JOIN photos p ON w.id = p.waypoint_id
    LEFT JOIN photo_tags pt ON p.id = pt.photo_id
    WHERE t.id = $1
    ORDER BY w.order_index ASC, p.id ASC, pt.id ASC
`

	rows, err := r.db.QueryContext(ctx, query, tripID, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var trip *domain.Trip
	// Waypoint ve Photo'ları ID'lerine göre takip etmek için map kullanıyoruz (Duplicate önlemek için)
	waypointMap := make(map[uuid.UUID]*domain.Waypoint)
	photoMap := make(map[uuid.UUID]*domain.Photo)

	for rows.Next() {
		var (
			// Trip alanları
			tID, tUserID                          uuid.UUID
			tTitle                                string
			tCover, tDesc                         *string
			tIsActive, tIsPublic, tIsLiked        bool
			tPublished, tCreated                  time.Time
			tViewCount, tLikeCount, tCommentCount int
			wpID                                  *uuid.UUID
			wpTitle, wpDesc, wpNote               *string
			wpOrder                               *int
			wpLat, wpLon                          *float64
			wpCreated                             *time.Time
			pID                                   *uuid.UUID
			pURL                                  *string
			// 🆕 Etiket alanları
			tagID      *uuid.UUID
			tagLabel   *string
			tagX, tagY *float64
		)

		err := rows.Scan(
			&tID, &tUserID, &tTitle, &tDesc, &tCover, &tIsActive, &tIsPublic, &tPublished, &tViewCount, &tLikeCount, &tCommentCount, &tCreated,
			&tIsLiked, &wpID, &wpTitle, &wpDesc, &wpOrder, &wpLat, &wpLon, &wpNote, &wpCreated,
			&pID, &pURL,
			&tagID, &tagLabel, &tagX, &tagY,
		)
		if err != nil {
			return nil, err
		}

		// 1. Trip nesnesini sadece ilk satırda oluştur
		if trip == nil {
			trip = &domain.Trip{
				ID: tID, UserID: tUserID, Title: tTitle, Description: tDesc,
				CoverImageURL: tCover, IsActive: tIsActive, IsPublic: tIsPublic,
				PublishedAt: tPublished, ViewCount: tViewCount, LikeCount: tLikeCount, CommentCount: tCommentCount, CreatedAt: tCreated,
				Waypoints: []domain.Waypoint{},
				IsLiked:   tIsLiked,
			}
		}

		// 2. Eğer Waypoint varsa işle (LEFT JOIN olduğu için nil gelebilir)
		if wpID != nil {
			if _, exists := waypointMap[*wpID]; !exists {
				wp := &domain.Waypoint{
					ID: *wpID, TripID: tID, Title: *wpTitle, Description: *wpDesc,
					OrderIndex: *wpOrder, Latitude: *wpLat, Longitude: *wpLon,
					Note: *wpNote, CreatedAt: *wpCreated,
					Photos: []domain.Photo{},
				}
				waypointMap[*wpID] = wp
				trip.Waypoints = append(trip.Waypoints, *wp)
			}

			// 3. Photo İşleme
			if pID != nil {
				// Trip.Waypoints içindeki doğru waypoint'i referans alalım
				var currentWp *domain.Waypoint
				for i := range trip.Waypoints {
					if trip.Waypoints[i].ID == *wpID {
						currentWp = &trip.Waypoints[i]
						break
					}
				}

				if currentWp != nil {
					photo, photoExists := photoMap[*pID]
					if !photoExists {
						photo = &domain.Photo{
							ID: *pID, WaypointID: *wpID, URL: *pURL,
							Tags: []domain.Tag{}, // 👈 Boş tag listesi ile başlat
						}
						photoMap[*pID] = photo
						currentWp.Photos = append(currentWp.Photos, *photo)
					}

					// 4. Tag İşleme (Eğer varsa)
					if tagID != nil {
						// Mevcut fotoğrafın içine tag'i ekle
						// Not: append yaparken pointer üzerinden gitmek için
						// photoMap'teki nesneyi veya Waypoint içindeki nesneyi güncellemeliyiz
						for j := range currentWp.Photos {
							if currentWp.Photos[j].ID == *pID {
								// Aynı tag'in tekrar eklenmesini önle
								alreadyAdded := false
								for _, existingTag := range currentWp.Photos[j].Tags {
									if existingTag.ID == *tagID {
										alreadyAdded = true
										break
									}
								}
								if !alreadyAdded {
									currentWp.Photos[j].Tags = append(currentWp.Photos[j].Tags, domain.Tag{
										ID: *tagID, Label: *tagLabel, XPos: *tagX, YPos: *tagY,
									})
								}
							}
						}
					}
				}
			}
		}
	}

	if trip == nil {
		return nil, fmt.Errorf("trip not found")
	}

	return trip, nil
}
func (r *Repository) IncrementUniqueView(ctx context.Context, tripID, userID uuid.UUID) error {
	// Sorgu Mantığı:
	// 1. trip_views tablosuna (trip_id, user_id) eklemeye çalış.
	// 2. Eğer zaten varsa (ON CONFLICT) hiçbir şey yapma.
	// 3. Eğer yeni bir satır eklendiyse (RETURNING), gidip trips tablosundaki sayacı artır.

	query := `
		WITH inserted_view AS (
			INSERT INTO trip_views (trip_id, user_id)
			VALUES ($1, $2)
			ON CONFLICT (trip_id, user_id) DO NOTHING
			RETURNING trip_id
		)
		UPDATE trips
		SET view_count = view_count + 1
		WHERE id IN (SELECT trip_id FROM inserted_view);
	`

	_, err := r.db.ExecContext(ctx, query, tripID, userID)
	if err != nil {
		return fmt.Errorf("failed to increment unique view: %w", err)
	}

	return nil
}
