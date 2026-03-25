package usecase

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type GetUserTripsUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, page, limit int) ([]domain.TripSummary, error)
}

type getUserTripsUseCase struct {
	tripRepo domain.TripRepository
}

func NewGetUserTripsUseCase(tripRepo domain.TripRepository) GetUserTripsUseCase {
	return &getUserTripsUseCase{tripRepo: tripRepo}
}

func (uc *getUserTripsUseCase) Execute(ctx context.Context, userID uuid.UUID, page, limit int) ([]domain.TripSummary, error) {
	trips, err := uc.tripRepo.GetUserTrips(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	return trips, nil
}
