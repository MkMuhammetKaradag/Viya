// trip-service/infrastructure/worker/processor.go
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
		return fmt.Errorf("json unmarshal failed: %w", err)
	}

	// 1. Dosyayı diskten oku
	fileBytes, err := os.ReadFile(payload.FilePath)
	if err != nil {
		// Eğer dosya diskte yoksa retry etmenin anlamı yok (SkipRetry)
		return fmt.Errorf("file not found at %s: %w", payload.FilePath, asynq.SkipRetry)
	}

	// 2. Cloudinary'ye yükle
	url, err := p.cloudinarySvc.UploadImageFromBytes(ctx, fileBytes, domain.UploadOptions{
		WayPointID: payload.WayPointID,
		Folder:     "waypoint_photos",
	})
	if err != nil {
		return fmt.Errorf("cloudinary upload failed: %w", err)
	}

	// 3. DB'ye kaydet
	wpID, err := uuid.Parse(payload.WayPointID)
	if err != nil {
		return fmt.Errorf("invalid uuid %s: %w", payload.WayPointID, asynq.SkipRetry)
	}

	if err := p.repo.AddWaypointPhotos(ctx, wpID, []string{url}); err != nil {
		return fmt.Errorf("db persistence failed: %w", err)
	}

	// 4. İşlem başarıyla bittiğinde dosyayı temizle
	if err := os.Remove(payload.FilePath); err != nil {
		log.Printf("Warning: temporary file could not be removed: %v", err)
	}

	log.Printf("Successfully processed photo for waypoint: %s", payload.WayPointID)
	return nil
}

func (p *TaskProcessor) Stop() {
	log.Println("Worker Processor durduruluyor...")
	p.server.Shutdown()
}
