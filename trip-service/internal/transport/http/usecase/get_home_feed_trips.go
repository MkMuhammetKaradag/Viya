package usecase

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type GetHomeFeedTripsUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.TripExploreDTO, error)
}

type getHomeFeedTripsUseCase struct {
	tripRepo domain.TripRepository
}

func NewGetHomeFeedTripsUseCase(tripRepo domain.TripRepository, worker domain.Worker) GetHomeFeedTripsUseCase {
	return &getHomeFeedTripsUseCase{tripRepo: tripRepo}
}

func (uc *getHomeFeedTripsUseCase) Execute(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.TripExploreDTO, error) {
	trip, err := uc.tripRepo.GetHomeFeedTrips(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return trip, nil
}
