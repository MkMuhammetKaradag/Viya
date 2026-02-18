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
}

func NewAddWayPointPhotosUseCase(tripRepo domain.TripRepository) AddWayPointPhotosUseCase {
	return &addWayPointPhotosUseCase{tripRepo: tripRepo}
}

func (uc *addWayPointPhotosUseCase) Execute(ctx context.Context, wayPointID uuid.UUID, files []*multipart.FileHeader) error {

	fmt.Println("id:", wayPointID)
	for i, fileHeader := range files {

		file, _ := fileHeader.Open()
		buf := new(bytes.Buffer)
		buf.ReadFrom(file)
		file.Close()
		fmt.Printf(" File %d: Name=%s, Size=%d bytes\n  ", i+1, fileHeader.Filename, fileHeader.Size)

	}
	return nil
}
