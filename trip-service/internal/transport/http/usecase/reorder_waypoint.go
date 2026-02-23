package usecase

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type ReorderUseCase interface {
	Execute(ctx context.Context, wayPointID uuid.UUID, index int) error
}

type reorderUseCase struct {
	tripRepo domain.TripRepository
}

func NewReorderUseCase(tripRepo domain.TripRepository) ReorderUseCase {
	return &reorderUseCase{tripRepo: tripRepo}
}

func (uc *reorderUseCase) Execute(ctx context.Context, wayPointID uuid.UUID, index int) error {

	err := uc.tripRepo.ReorderWaypoints(ctx, wayPointID, index)
	if err != nil {
		return err
	}

	return nil

}
