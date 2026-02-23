package domain

import "github.com/google/uuid"

type UpdateWaypointInput struct {
	WayPointID  uuid.UUID
	Title       *string
	Description *string
	Latitude    *float64
	Longitude   *float64
}
