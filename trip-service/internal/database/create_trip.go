package database

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) CreateTrip(ctx context.Context, trip *domain.Trip) (uuid.UUID, error) {
	query := `
	     INSERT INTO trips (user_id,title,description,is_active)
		VALUES ($1,$2,$3,$4)
		RETURNING id
	`

	var newID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, trip.UserID, trip.Title, trip.Description, trip.IsActive).Scan(&newID)
	if err != nil {
		return uuid.Nil, err
	}

	return newID, nil
}
