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
	mux.HandleFunc(domain.TaskIncrementTripView, p.ProcessIncrementViewTask)

	log.Println("Worker Processor başlatılıyor...")
	return p.server.Run(mux)
}

func (p *TaskProcessor) ProcessWaypointUploadTask(ctx context.Context, t *asynq.Task) error {
	var payload domain.UploadWaypointPhotoTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json unmarshal failed: %w", err)
	}

	// 1. Etiketleri hemen çöz (Unmarshal)
	var tags []domain.Tag
	if payload.Tags != "" {
		// React Native'den gelen string'i struct listesine çeviriyoruz
		if err := json.Unmarshal([]byte(payload.Tags), &tags); err != nil {
			return fmt.Errorf("tags unmarshal failed: %w", err)
		}
	}

	// 2. Dosyayı oku
	fileBytes, err := os.ReadFile(payload.FilePath)
	if err != nil {
		return fmt.Errorf("file not found: %w", asynq.SkipRetry)
	}

	// 3. Cloudinary'ye yükle
	url, err := p.cloudinarySvc.UploadImageFromBytes(ctx, fileBytes, domain.UploadOptions{
		WayPointID: payload.WayPointID,
		Folder:     "waypoint_photos",
	})
	if err != nil {
		return fmt.Errorf("cloudinary upload failed: %w", err)
	}

	// 4. DB'ye kaydet (Yeni metodumuzla)
	wpID, _ := uuid.Parse(payload.WayPointID)
	// AddWaypointPhotos yerine AddWaypointPhotoWithTags kullanıyoruz
	if err := p.repo.AddWaypointPhotoWithTags(ctx, wpID, url, tags); err != nil {
		return fmt.Errorf("db persistence failed: %w", err)
	}

	// 5. Temizlik
	os.Remove(payload.FilePath)
	return nil
}
func (p *TaskProcessor) ProcessIncrementViewTask(ctx context.Context, t *asynq.Task) error {
	var payload domain.IncrementTripViewPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	// Repository'deki o meşhur ON CONFLICT'li fonksiyonu burada çağırıyoruz
	// Bu sayede DB işlemleri arka planda, kullanıcıyı bekletmeden hallolur.
	return p.repo.IncrementUniqueView(ctx, payload.TripID, payload.UserID)
}
func (p *TaskProcessor) Close() error {
	log.Println("Worker Processor durduruluyor...")
	p.server.Shutdown()

	return nil
}
