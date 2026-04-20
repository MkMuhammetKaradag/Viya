package domain

import "github.com/google/uuid"

type Worker interface {
	EnqueueUploadWaypointPhoto(payload UploadWaypointPhotoTaskPayload) error
	EnqueueTripEmbedding(tripID uuid.UUID) error
	EnqueueIncrementTrip(tripID, userID uuid.UUID, weight float32, action string) error
}

type UploadWaypointPhotoTaskPayload struct {
	WayPointID string `json:"waypoint_id"`
	FilePath   string `json:"file_path"`
	Tags       string `json:"tags"` // JSON string olarak React Native'den gelen veri
}

const TaskIncrementTrip = "task:increment_trip"

type InteractionTripPayload struct {
	UserID uuid.UUID `json:"user_id"`
	TripID uuid.UUID `json:"trip_id"`
	Weight float32   `json:"weight"` // 0.05, 0.5, -0.3 vb.
	Action string    `json:"action"` // "view", "like", "unlike", "comment"
}

const TaskGenerateTripEmbedding = "task:generate_trip_embedding"

type TripEmbeddingPayload struct {
	TripID uuid.UUID `json:"trip_id"`
}
