package domain

import "github.com/google/uuid"

type Worker interface {
	EnqueueUploadWaypointPhoto(payload UploadWaypointPhotoTaskPayload) error
}

type UploadWaypointPhotoTaskPayload struct {
	WayPointID string `json:"waypoint_id"`
	FilePath   string `json:"file_path"` // Diskteki geçici yol
}

const TaskIncrementTripView = "task:increment_trip_view"

type IncrementTripViewPayload struct {
	TripID uuid.UUID `json:"trip_id"`
	UserID uuid.UUID `json:"user_id"`
}
