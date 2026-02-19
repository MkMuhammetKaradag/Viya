package domain

type Worker interface {
	EnqueueUploadWaypointPhoto(payload UploadWaypointPhotoTaskPayload) error
}

type UploadWaypointPhotoTaskPayload struct {
	WayPointID string `json:"waypoint_id"`
	FilePath   string `json:"file_path"` // Diskteki geçici yol
}
