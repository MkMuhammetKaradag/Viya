package domain

import (
	"context"

	"github.com/google/uuid"
)

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *Trip) (uuid.UUID, error)
	AddWaypoint(ctx context.Context, wp *Waypoint) (uuid.UUID, error)
	Close() error
}
