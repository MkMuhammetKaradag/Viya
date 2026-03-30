package domain

import (
	"time"

	"github.com/google/uuid"
)

//(gis) PostGIS   unutma! not

type Trip struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	Title         string    `json:"title" db:"title"`
	Description   string    `json:"description" db:"description"`
	CoverImageURL *string   `json:"cover_image_url" db:"cover_image_url"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	IsPublic      bool      `json:"is_public" db:"is_public"`
	PublishedAt   time.Time `json:"published_at" db:"published_at"`
	ViewCount     int       `json:"view_count" db:"view_count"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`

	Waypoints []Waypoint `json:"waypoints,omitempty"`
}
type Photo struct {
	ID         uuid.UUID `json:"id" db:"id"`
	WaypointID uuid.UUID `json:"waypoint_id" db:"waypoint_id"`
	URL        string    `json:"url" db:"url"`
	Tags       []Tag     `json:"tags,omitempty"`
}

type Tag struct {
	ID      uuid.UUID `json:"id"`
	PhotoID uuid.UUID `json:"photo_id"`
	Label   string    `json:"label"` // Etiket metni
	XPos    float64   `json:"x_pos"` // Resmin genişliğinin % kaçında?
	YPos    float64   `json:"y_pos"` // Resmin yüksekliğinin % kaçında?
}
type Waypoint struct {
	ID          uuid.UUID `json:"id" db:"id"`
	TripID      uuid.UUID `json:"trip_id" db:"trip_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	OrderIndex  int       `json:"order_index" db:"order_index"`
	Latitude    float64   `json:"latitude" db:"latitude"`
	Longitude   float64   `json:"longitude" db:"longitude"`
	Note        string    `json:"note" db:"note"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`

	Photos []Photo `json:"photos,omitempty"`
}
type TripSummary struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	CoverImageURL *string   `json:"cover_image_url"`
	IsPublic      bool      `json:"is_public"`
	ViewCount     int       `json:"view_count"`
	WaypointCount int       `json:"waypoint_count"`
	CreatedAt     time.Time `json:"created_at"`
}
