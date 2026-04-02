package domain

import (
	"context"

	"github.com/google/uuid"
)

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *Trip) (uuid.UUID, error)
	AddWaypoint(ctx context.Context, wp *Waypoint) (uuid.UUID, error)
	AddWaypointPhotos(ctx context.Context, waypointID uuid.UUID, photoURLs []string) error
	AddWaypointPhotoWithTags(ctx context.Context, wpID uuid.UUID, photoURL string, tags []Tag) error
	// GetTripByID(ctx context.Context, tripID uuid.UUID) (*Trip, error)

	GetTripWithWaypointsAndPhotos(ctx context.Context, tripID uuid.UUID) (*Trip, error)
	IncrementUniqueView(ctx context.Context, tripID, userID uuid.UUID) error

	GetUserTrips(ctx context.Context, userID uuid.UUID, page, limit int) ([]TripSummary, error)

	DeleteWaypoint(ctx context.Context, waypointID uuid.UUID) error
	ReorderWaypoints(ctx context.Context, wpID uuid.UUID, index int) error
	GetWaypointByID(ctx context.Context, id uuid.UUID) (*Waypoint, error)
	UpdateWaypoint(ctx context.Context, wp *Waypoint) error

	CreateUser(ctx context.Context, id uuid.UUID, username, email string) error
	Close() error
}
