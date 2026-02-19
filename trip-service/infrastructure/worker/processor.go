// internal/product-service/infrastructure/worker/processor.go
package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"trip-service/internal/domain"

	"os"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type TaskProcessor struct {
	server        *asynq.Server
	repo          domain.TripRepository
	cloudinarySvc domain.ImageService
}

func NewTaskProcessor(redisOpt asynq.RedisClientOpt, repo domain.TripRepository, cloudinarySvc domain.ImageService) *TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 5,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
		},
	})

	return &TaskProcessor{
		server:        server,
		repo:          repo,
		cloudinarySvc: cloudinarySvc,
	}
}
func (p *TaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskUploadWaypointPhoto, p.ProcessWaypointUploadTask)

	log.Println("Worker Processor başlatılıyor...")
	return p.server.Run(mux)
}

func (p *TaskProcessor) ProcessWaypointUploadTask(ctx context.Context, t *asynq.Task) error {
	var payload domain.UploadWaypointPhotoTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	// 1. Dosyayı diskten aç
	fileBytes, err := os.ReadFile(payload.FilePath)
	if err != nil {
		return fmt.Errorf("dosya okunamadı (belki silindi?): %w", err)
	}

	// 2. Cloudinary'ye yükle
	url, err := p.cloudinarySvc.UploadImageFromBytes(ctx, fileBytes, domain.UploadOptions{
		WayPointID: payload.WayPointID,
		Folder:     "waypoint_photos",
	})
	if err != nil {
		return err // Retry mekanizması çalışır
	}

	// 3. DB'ye kaydet
	wpID, _ := uuid.Parse(payload.WayPointID)
	if err := p.repo.AddWaypointPhotos(ctx, wpID, []string{url}); err != nil {
		return err
	}

	// 4. İŞLEM TAMAM: Geçici dosyayı diskten sil (Cleanup)
	_ = os.Remove(payload.FilePath)

	return nil
}
