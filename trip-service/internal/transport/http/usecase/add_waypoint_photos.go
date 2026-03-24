package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type AddWayPointPhotosUseCase interface {
	Execute(ctx context.Context, wayPointID uuid.UUID, files []*multipart.FileHeader) error
}

type addWayPointPhotosUseCase struct {
	tripRepo domain.TripRepository
	imgSvc   domain.ImageService
	worker   domain.Worker
}

func NewAddWayPointPhotosUseCase(tripRepo domain.TripRepository, imgSvc domain.ImageService, worker domain.Worker) AddWayPointPhotosUseCase {
	return &addWayPointPhotosUseCase{tripRepo: tripRepo, imgSvc: imgSvc, worker: worker}
}

func (uc *addWayPointPhotosUseCase) Execute(ctx context.Context, wayPointID uuid.UUID, files []*multipart.FileHeader) error {

	for _, fileHeader := range files {
		// 1. Geçici bir dosya ismi oluştur
		tempFileName := fmt.Sprintf("%s_%d%s", wayPointID.String(), time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
		tempPath := filepath.Join("tmp/uploads", tempFileName)

		// 2. Dosyayı diske kaydet (Helper fonksiyon kullandığını varsayıyorum)
		if err := saveFileToDisk(fileHeader, tempPath); err != nil {
			return fmt.Errorf("dosya geçici olarak kaydedilemedi: %w", err)
		}

		// 3. Worker'a sadece dosya yolunu gönder
		payload := domain.UploadWaypointPhotoTaskPayload{
			WayPointID: wayPointID.String(),
			FilePath:   tempPath,
		}

		if err := uc.worker.EnqueueUploadWaypointPhoto(payload); err != nil {
			return fmt.Errorf("iş kuyruğa alınamadı: %w", err)
		}
	}
	return nil
}
