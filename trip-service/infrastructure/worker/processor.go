// trip-service/infrastructure/worker/processor.go
package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"trip-service/internal/domain"

	"os"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type TaskProcessor struct {
	server        *asynq.Server
	repo          domain.TripRepository
	cloudinarySvc domain.ImageService
	ai            domain.AIService
}

func NewTaskProcessor(redisOpt asynq.RedisClientOpt, repo domain.TripRepository, cloudinarySvc domain.ImageService, ai domain.AIService) *TaskProcessor {
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
		ai:            ai,
	}
}
func (p *TaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskUploadWaypointPhoto, p.ProcessWaypointUploadTask)
	mux.HandleFunc(domain.TaskIncrementTrip, p.ProcessIncrementTripTask)
	mux.HandleFunc(domain.TaskGenerateTripEmbedding, p.ProcessTripEmbeddingTask)

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
func (p *TaskProcessor) ProcessIncrementTripTask(ctx context.Context, t *asynq.Task) error {
	var payload domain.InteractionTripPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	// 1. Mevcut İşlem: Görüntülenme sayısını artır
	switch payload.Action {
	case "view":
		p.repo.IncrementUniqueView(ctx, payload.TripID, payload.UserID)
	case "like":

	}

	if err := p.repo.UpdateUserInterest(ctx, payload.UserID, payload.TripID, payload.Weight); err != nil {
		// AI güncellemesi kritik değilse sadece logla, görevi tamamen iptal etme
		fmt.Printf("User interest update failed: %v\n", err)
	}

	return nil
}
func (p *TaskProcessor) ProcessTripEmbeddingTask(ctx context.Context, t *asynq.Task) error {

	var payload domain.TripEmbeddingPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	// 1. Trip detaylarını DB'den çek (Prompt için gerekli tüm bilgilerle)
	trip, err := p.repo.GetTripByIDForAI(ctx, payload.TripID)
	if err != nil {
		return fmt.Errorf("trip not found for embedding: %w", err)
	}

	// 2. Senin o meşhur "Zengin Prompt"u hazırla
	prompt := buildRichPrompt(trip)

	// 3. Ollama'ya git ve vektörü al
	vector, err := p.ai.GetVector(ctx, prompt)
	if err != nil {
		return fmt.Errorf("ollama vector failed: %w", err)
	}

	// 4. DB'deki content_vector kolonunu güncelle
	return p.repo.UpdateTripEmbedding(ctx, payload.TripID, vector)
}

func buildRichPrompt(trip *domain.Trip) string {
	// 1. Temel Bilgiler
	prompt := fmt.Sprintf("Travel Trip: %s. Description: %s. ", trip.Title, trip.Description)

	// 2. Kategori Bilgisi (Çok önemli!)
	// Not: Buraya category_id yerine category_name (isim) gelmeli.
	// Eğer sadece ID varsa, DB'den isimleri join ile çekip buraya eklemelisin.
	if len(trip.CategoryNames) > 0 {
		prompt += "Categories: " + strings.Join(trip.CategoryNames, ", ") + ". "
	}

	// 3. Duraklar ve Yerel Bilgiler (Waypoint Context)
	if len(trip.Waypoints) > 0 {
		prompt += "Route details: "
		for _, wp := range trip.Waypoints {
			// Durak ismi ve oradaki aktiviteyi ekle
			prompt += fmt.Sprintf("Stopped at %s. Activities: %s. ", wp.Title, wp.Description)
			if wp.Note != "" {
				prompt += fmt.Sprintf("Note: %s. ", wp.Note)
			}
		}
	}

	// 4. Konum Bilgisi
	if trip.LocationName != nil {
		prompt += fmt.Sprintf("Location: %s.", trip.LocationName)
	}

	return prompt
}
func (p *TaskProcessor) Close() error {
	log.Println("Worker Processor durduruluyor...")
	p.server.Shutdown()

	return nil
}
