package database

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) GetUserTrips(ctx context.Context, currentUserID, targetUserID uuid.UUID, page, limit int) ([]domain.TripSummary, error) {
	// Offset hesabı: (1. sayfa için 0 atla, 2. sayfa için limit kadar atla)
	offset := (page - 1) * limit

	query := `
        SELECT 
            t.id, 
            t.title, 
            COALESCE(
                t.cover_image_url, 
                (SELECT p.url FROM photos p 
                 JOIN waypoints w ON p.waypoint_id = w.id 
                 WHERE w.trip_id = t.id 
                 ORDER BY w.order_index ASC, p.created_at ASC 
                 LIMIT 1)
            ) as effective_cover_url,
            t.is_public, 
            t.view_count, 
            t.created_at,
            (SELECT COUNT(*) FROM waypoints WHERE trip_id = t.id) as waypoint_count
        FROM trips t
        JOIN users u ON t.user_id = u.id
        WHERE t.user_id = $2 -- Hedef kullanıcının tripleri
          AND t.deleted_at IS NULL
          -- 🛡️ GİZLİLİK VE BLOK KONTROLÜ
          AND (
              -- 1. Durum: Profil açık ve kullanıcı bloklu değil
              (u.is_private = FALSE AND NOT EXISTS (
                  SELECT 1 FROM local_blocks 
                  WHERE (blocker_id = $1 AND blocked_id = $2) OR (blocker_id = $2 AND blocked_id = $1)
              ))
              OR 
              -- 2. Durum: Profil gizli ama istek atan kişi takipçi (ve bloklu değil)
              (u.is_private = TRUE AND EXISTS (
                  SELECT 1 FROM local_follows 
                  WHERE follower_id = $1 AND following_id = $2 AND status = 'ACCEPTED'
              ) AND NOT EXISTS (
                  SELECT 1 FROM local_blocks 
                  WHERE (blocker_id = $1 AND blocked_id = $2) OR (blocker_id = $2 AND blocked_id = $1)
              ))
              OR
              -- 3. Durum: Kişi kendi profiline bakıyorsa her şeyi görür
              ($1 = $2)
          )
        ORDER BY t.created_at DESC
        LIMIT $3 OFFSET $4
    `

	rows, err := r.db.QueryContext(ctx, query, currentUserID, targetUserID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []domain.TripSummary
	for rows.Next() {
		var s domain.TripSummary
		err := rows.Scan(
			&s.ID, &s.Title, &s.CoverImageURL, &s.IsPublic, &s.ViewCount, &s.CreatedAt, &s.WaypointCount,
		)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}

	return summaries, nil
}
func (r *Repository) GetMeTrips(ctx context.Context, userID uuid.UUID, page, limit int) ([]domain.TripSummary, error) {
	// Offset hesabı: (1. sayfa için 0 atla, 2. sayfa için limit kadar atla)
	offset := (page - 1) * limit

	query := `
        SELECT 
            t.id, 
            t.title, 
            COALESCE(
                t.cover_image_url, 
                (SELECT p.url FROM photos p 
                 JOIN waypoints w ON p.waypoint_id = w.id 
                 WHERE w.trip_id = t.id 
                 ORDER BY w.order_index ASC, p.created_at ASC 
                 LIMIT 1)
            ) as effective_cover_url,
            t.is_public, 
            t.view_count, 
            t.created_at,
            (SELECT COUNT(*) FROM waypoints WHERE trip_id = t.id) as waypoint_count
        FROM trips t
        WHERE t.user_id = $1
        ORDER BY t.created_at DESC
        LIMIT $2 OFFSET $3 -- Sayfalama sihrimiz burada
    `

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []domain.TripSummary
	for rows.Next() {
		var s domain.TripSummary
		err := rows.Scan(
			&s.ID, &s.Title, &s.CoverImageURL, &s.IsPublic, &s.ViewCount, &s.CreatedAt, &s.WaypointCount,
		)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}

	return summaries, nil
}
