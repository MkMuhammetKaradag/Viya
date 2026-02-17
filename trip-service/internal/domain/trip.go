package domain

import "github.com/google/uuid"

type Trip struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   string    `json:"created_at"`
}

type Waypoint struct {
	ID        uuid.UUID `json:"id"`
	TripID    uuid.UUID `json:"trip_id"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	Note      string    `json:"note,omitempty"`
	CreatedAt string    `json:"created_at"`
}
