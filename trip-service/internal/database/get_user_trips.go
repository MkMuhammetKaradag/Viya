package database

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) GetUserTrips(ctx context.Context, userID uuid.UUID, page, limit int) ([]domain.TripSummary, error) {
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
