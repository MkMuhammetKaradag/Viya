package usecase

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type GetTripUseCase interface {
	Execute(ctx context.Context, tripID uuid.UUID) (*domain.Trip, error)
}

type getTripUseCase struct {
	tripRepo domain.TripRepository
}

func NewGetTripUseCase(tripRepo domain.TripRepository) GetTripUseCase {
	return &getTripUseCase{tripRepo: tripRepo}
}

func (uc *getTripUseCase) Execute(ctx context.Context, tripID uuid.UUID) (*domain.Trip, error) {
	trip, err := uc.tripRepo.GetTripByID(ctx, tripID)
	if err != nil {
		return nil, err
	}
	return trip, nil
}
