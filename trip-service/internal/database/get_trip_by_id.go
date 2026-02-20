package database

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) GetTripByID(ctx context.Context, tripID uuid.UUID) (*domain.Trip, error) {
	query := `
		SELECT id, user_id, title, description, is_active
		FROM trips
		WHERE id = $1`

	var trip domain.Trip
	err := r.db.QueryRowContext(ctx, query, tripID).Scan(&trip.ID, &trip.UserID, &trip.Title, &trip.Description, &trip.IsActive)
	if err != nil {
		return nil, err

	}
	return &trip, nil
}
