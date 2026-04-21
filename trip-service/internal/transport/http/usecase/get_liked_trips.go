package usecase

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type GetLikedTripsUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, limit, page int) ([]domain.TripSummary, error)
}

type getLikedTripsUseCase struct {
	tripRepo domain.TripRepository
}

func NewGetLikedTripsUseCase(tripRepo domain.TripRepository) GetLikedTripsUseCase {
	return &getLikedTripsUseCase{tripRepo: tripRepo}
}

func (uc *getLikedTripsUseCase) Execute(ctx context.Context, userID uuid.UUID, limit, page int) ([]domain.TripSummary, error) {
	trip, err := uc.tripRepo.GetLikedTrips(ctx, userID, limit, page)
	if err != nil {
		return nil, err
	}

	return trip, nil
}
