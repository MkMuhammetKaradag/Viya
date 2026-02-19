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

	// fmt.Println("id:", wayPointID)
	// var photoURLs []string
	// for i, fileHeader := range files {

	// 	file, _ := fileHeader.Open()
	// 	buf := new(bytes.Buffer)
	// 	buf.ReadFrom(file)
	// 	file.Close()
	// 	fmt.Printf("File %d: Name=%s, Size=%d bytes\n  ", i+1, fileHeader.Filename, fileHeader.Size)

	// 	url, err := uc.imgSvc.UploadImageFromBytes(ctx, buf.Bytes(), domain.UploadOptions{
	// 		Folder:     "waypoint_photos",
	// 		WayPointID: wayPointID.String(),
	// 	})

	// 	if err != nil {

	// 		return fmt.Errorf("failed to upload image %s: %w", fileHeader.Filename, err)

	// 	}
	// 	fmt.Println(url)
	// 	photoURLs = append(photoURLs, url)

	// }

	// err := uc.tripRepo.AddWaypointPhotos(ctx, wayPointID, photoURLs)
	// if err != nil {
	// 	return fmt.Errorf("failed to add photo URLs to waypoint: %w", err)
	// }

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
