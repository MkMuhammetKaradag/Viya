package usecase

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type AddWayPointPhotosUseCase interface {
	Execute(ctx context.Context, wayPointID uuid.UUID, files []*multipart.FileHeader) error
}

type addWayPointPhotosUseCase struct {
	tripRepo domain.TripRepository
	imgSvc   domain.ImageService
}

func NewAddWayPointPhotosUseCase(tripRepo domain.TripRepository, imgSvc domain.ImageService) AddWayPointPhotosUseCase {
	return &addWayPointPhotosUseCase{tripRepo: tripRepo, imgSvc: imgSvc}
}

func (uc *addWayPointPhotosUseCase) Execute(ctx context.Context, wayPointID uuid.UUID, files []*multipart.FileHeader) error {

	fmt.Println("id:", wayPointID)
	var photoURLs []string
	for i, fileHeader := range files {

		file, _ := fileHeader.Open()
		buf := new(bytes.Buffer)
		buf.ReadFrom(file)
		file.Close()
		fmt.Printf("File %d: Name=%s, Size=%d bytes\n  ", i+1, fileHeader.Filename, fileHeader.Size)

		url, err := uc.imgSvc.UploadImageFromBytes(ctx, buf.Bytes(), domain.UploadOptions{
			Folder:     "waypoint_photos",
			WayPointID: wayPointID.String(),
		})

		if err != nil {

			return fmt.Errorf("failed to upload image %s: %w", fileHeader.Filename, err)

		}
		fmt.Println(url)
		photoURLs = append(photoURLs, url)

	}

	err := uc.tripRepo.AddWaypointPhotos(ctx, wayPointID, photoURLs)
	if err != nil {
		return fmt.Errorf("failed to add photo URLs to waypoint: %w", err)
	}
	return nil
}
