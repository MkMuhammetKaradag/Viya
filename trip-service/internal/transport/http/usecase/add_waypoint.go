package usecase

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type AddWayPointUseCase interface {
	Execute(ctx context.Context, wp *domain.Waypoint, files []*multipart.FileHeader, tags []string) (uuid.UUID, error)
}

type addWayPointUseCase struct {
	tripRepo domain.TripRepository
	worker   domain.Worker
}

func NewAddWayPointUseCase(tripRepo domain.TripRepository, worker domain.Worker) AddWayPointUseCase {
	return &addWayPointUseCase{tripRepo: tripRepo, worker: worker}
}

func (uc *addWayPointUseCase) Execute(ctx context.Context, wp *domain.Waypoint, files []*multipart.FileHeader, tags []string) (uuid.UUID, error) {
	// A. Önce Waypoint'i veritabanına ekle (ID almamız şart)
	wpID, err := uc.tripRepo.AddWaypoint(ctx, wp)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not save waypoint: %w", err)
	}

	// B. Eğer fotoğraf varsa, onları diske kaydet ve Worker'a haber ver
	if len(files) > 0 {
		for i, fileHeader := range files {
			currentFile := fileHeader

			// Güvenlik: Eğer tag sayısı dosya sayısından azsa paniklememesi için kontrol
			currentTags := ""
			if i < len(tags) {
				currentTags = tags[i]
			}

			tempFileName := fmt.Sprintf("%s_%d%s", wpID.String(), time.Now().UnixNano(), filepath.Ext(currentFile.Filename))
			tempPath := filepath.Join("tmp/uploads", tempFileName)

			if err := saveFileToDisk(currentFile, tempPath); err != nil {
				fmt.Printf("Warning: Failed to save temp file: %v\n", err)
				continue
			}
			fmt.Println("i:", i, " tags:", currentTags)

			// PAYLOAD'A TAGLARI EKLE
			payload := domain.UploadWaypointPhotoTaskPayload{
				WayPointID: wpID.String(),
				FilePath:   tempPath,
				Tags:       currentTags, // İşte burası!
			}

			if err := uc.worker.EnqueueUploadWaypointPhoto(payload); err != nil {
				fmt.Printf("Warning: Could not enqueue task: %v\n", err)
			}
		}
	}

	return wpID, nil

}
func saveFileToDisk(file *multipart.FileHeader, dest string) error {
	// 1. Hedef klasörü oluştur (Eğer yoksa)
	dir := filepath.Dir(dest)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("klasör oluşturulamadı: %w", err)
	}

	// 2. Kaynak dosyayı aç
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// 3. Hedef dosyayı oluştur
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// 4. İçeriği kopyala
	_, err = io.Copy(out, src)
	return err
}
