package domain

import (
	"time"

	"github.com/google/uuid"
)

type Trip struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`

	WayPoints []Waypoint `json:"waypoints,omitempty"`
}

type Waypoint struct {
	ID          uuid.UUID `json:"id"`
	TripID      uuid.UUID `json:"trip_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OrderIndex  int       `json:"order_index"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	CreatedAt   time.Time `json:"created_at"`

	Photos []string `json:"photos"`
}
