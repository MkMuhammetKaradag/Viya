package database

import (
	"context"
	"database/sql"
	"fmt"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) GetWaypointByID(ctx context.Context, id uuid.UUID) (*domain.Waypoint, error) {
	query := `SELECT id, trip_id, title, description, latitude, longitude, order_index FROM waypoints WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var wp domain.Waypoint
	err := row.Scan(&wp.ID, &wp.TripID, &wp.Title, &wp.Description, &wp.Latitude, &wp.Longitude, &wp.OrderIndex)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Kayıt bulunamadığında nil döndür
		}
		return nil, fmt.Errorf("failed to get waypoint by ID: %w", err)
	}
	return &wp, nil
}
