package domain

import (
	"context"

	"github.com/google/uuid"
)

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *Trip) (uuid.UUID, error)
	AddWaypoint(ctx context.Context, wp *Waypoint) (uuid.UUID, error)
	AddWaypointPhotos(ctx context.Context, waypointID uuid.UUID, photoURLs []string) error
	GetTripByID(ctx context.Context, tripID uuid.UUID) (*Trip, error)
	DeleteWaypoint(ctx context.Context, waypointID uuid.UUID) error
	ReorderWaypoints(ctx context.Context, wpID uuid.UUID, index int) error
	GetWaypointByID(ctx context.Context, id uuid.UUID) (*Waypoint, error)
	UpdateWaypoint(ctx context.Context, wp *Waypoint) error
	Close() error
}
