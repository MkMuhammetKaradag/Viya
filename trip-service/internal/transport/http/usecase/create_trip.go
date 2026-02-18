package usecase

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type CreateTripUseCase interface {
	Execute(ctx context.Context, trip *domain.Trip) (uuid.UUID, error)
}

type createTripUseCase struct {
	tripRepo domain.TripRepository
}

func NewCreateTripUseCase(tripRepo domain.TripRepository) CreateTripUseCase {
	return &createTripUseCase{tripRepo: tripRepo}
}

func (uc *createTripUseCase) Execute(ctx context.Context, trip *domain.Trip) (uuid.UUID, error) {
	id, err := uc.tripRepo.CreateTrip(ctx, trip)
	if err != nil {
		return id, err
	}
	return id, nil
}
