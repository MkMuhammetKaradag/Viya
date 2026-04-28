package usecase

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type GetMeTripsUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, page, limit int) ([]domain.TripSummary, error)
}

type getMeTripsUseCase struct {
	tripRepo domain.TripRepository
}

func NewGetMeTripsUseCase(tripRepo domain.TripRepository) GetMeTripsUseCase {
	return &getMeTripsUseCase{tripRepo: tripRepo}
}

func (uc *getMeTripsUseCase) Execute(ctx context.Context, userID uuid.UUID, page, limit int) ([]domain.TripSummary, error) {
	trips, err := uc.tripRepo.GetMeTrips(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	return trips, nil
}
