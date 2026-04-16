package usecase

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type GetExploreTripsUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.TripExploreDTO, error)
}

type getExploreTripsUseCase struct {
	tripRepo domain.TripRepository
}

func NewGetExploreTripsUseCase(tripRepo domain.TripRepository, worker domain.Worker) GetExploreTripsUseCase {
	return &getExploreTripsUseCase{tripRepo: tripRepo}
}

func (uc *getExploreTripsUseCase) Execute(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.TripExploreDTO, error) {
	trip, err := uc.tripRepo.GetExploreTrips(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return trip, nil
}
