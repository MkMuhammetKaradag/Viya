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

func (w *Worker) EnqueueIncrementTrip(tripID, userID uuid.UUID, weight float32, action string) error {
	payload := domain.InteractionTripPayload{
		TripID: tripID,
		UserID: userID,
		Weight: weight,
		Action: action,
	}
	data, _ := json.Marshal(payload)

	task := asynq.NewTask(domain.TaskIncrementTrip, data, asynq.MaxRetry(3), asynq.Queue("default"))
	_, err := w.client.Enqueue(task)
	return err
}

func (w *Worker) EnqueueTripEmbedding(tripID uuid.UUID) error {
	payload := domain.TripEmbeddingPayload{TripID: tripID}
	data, _ := json.Marshal(payload)

	task := asynq.NewTask(domain.TaskGenerateTripEmbedding, data, asynq.MaxRetry(10))
	_, err := w.client.Enqueue(task)
	return err
}
