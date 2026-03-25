// trip-service/infrastructure/worker/worker.go
package worker

import (
	"encoding/json"
	"trip-service/internal/domain"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const TaskUploadWaypointPhoto = "task:upload_waypoint_photo"

type Worker struct {
	client *asynq.Client
}

func NewWorker(client *asynq.Client) *Worker {
	return &Worker{
		client: client,
	}
}

func (w *Worker) EnqueueUploadWaypointPhoto(payload domain.UploadWaypointPhotoTaskPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TaskUploadWaypointPhoto, data, asynq.MaxRetry(5), asynq.Queue("critical"))

	_, err = w.client.Enqueue(task)
	return err
}

func (w *Worker) EnqueueIncrementTripView(tripID, userID uuid.UUID) error {
	payload := domain.IncrementTripViewPayload{
		TripID: tripID,
		UserID: userID,
	}
	data, _ := json.Marshal(payload)

	// Görüntülenme sayısı "critical" bir iş değil, "default" kuyruğuna atabiliriz.
	task := asynq.NewTask(domain.TaskIncrementTripView, data, asynq.MaxRetry(3), asynq.Queue("default"))
	_, err := w.client.Enqueue(task)
	return err
}
