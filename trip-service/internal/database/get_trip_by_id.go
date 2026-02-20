package database

import (
	"context"
	"encoding/json"
	"fmt"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) GetTripByID(ctx context.Context, tripID uuid.UUID) (*domain.Trip, error) {
	query := `
		SELECT 
			t.id, t.user_id, t.title, t.description, t.is_active,
			COALESCE(
				(SELECT json_agg(json_build_object(
					'id', w.id,
					'trip_id', w.trip_id,
					'title', w.title,
					'description', w.description,
					'order_index', w.order_index,
					'latitude', w.latitude,
					'longitude', w.longitude,
					'photos', COALESCE(
						(SELECT json_agg(p.url) 
						 FROM photos p 
						 WHERE p.waypoint_id = w.id), 
						'[]'::json
					)
				) ORDER BY w.order_index)
				FROM waypoints w 
				WHERE w.trip_id = t.id), 
				'[]'::json
			) as waypoints
		FROM trips t
		WHERE t.id = $1`

	var trip domain.Trip
	var waypointsJSON []byte

	err := r.db.QueryRowContext(ctx, query, tripID).Scan(
		&trip.ID,
		&trip.UserID,
		&trip.Title,
		&trip.Description,
		&trip.IsActive,
		&waypointsJSON,
	)
	if err != nil {
		return nil, err
	}

	// JSON'u Go struct'ına unmarshal ediyoruz.
	// domain.Waypoint struct'ında Photos []string alanı olduğundan emin olmalısın.
	if err := json.Unmarshal(waypointsJSON, &trip.WayPoints); err != nil {
		return nil, fmt.Errorf("waypoints unmarshal error: %w", err)
	}

	return &trip, nil
}
