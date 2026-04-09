package domain

import "github.com/google/uuid"

type Worker interface {
	EnqueueUploadWaypointPhoto(payload UploadWaypointPhotoTaskPayload) error
	EnqueueTripEmbedding(tripID uuid.UUID) error
	EnqueueIncrementTripView(tripID, userID uuid.UUID) error
}

type UploadWaypointPhotoTaskPayload struct {
	WayPointID string `json:"waypoint_id"`
	FilePath   string `json:"file_path"`
	Tags       string `json:"tags"` // JSON string olarak React Native'den gelen veri
}

const TaskIncrementTripView = "task:increment_trip_view"

type IncrementTripViewPayload struct {
	TripID uuid.UUID `json:"trip_id"`
	UserID uuid.UUID `json:"user_id"`
}

const TaskGenerateTripEmbedding = "task:generate_trip_embedding"

type TripEmbeddingPayload struct {
	TripID uuid.UUID `json:"trip_id"`
}
