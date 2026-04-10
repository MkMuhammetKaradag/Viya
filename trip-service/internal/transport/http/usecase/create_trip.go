package usecase

import (
	"context"
	"fmt"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type CreateTripUseCase interface {
	Execute(ctx context.Context, trip *domain.Trip) (uuid.UUID, error)
}

type createTripUseCase struct {
	tripRepo domain.TripRepository
	worker   domain.Worker
}

func NewCreateTripUseCase(tripRepo domain.TripRepository, worker domain.Worker) CreateTripUseCase {
	return &createTripUseCase{tripRepo: tripRepo, worker: worker}
}

func (uc *createTripUseCase) Execute(ctx context.Context, trip *domain.Trip) (uuid.UUID, error) {
	id, err := uc.tripRepo.CreateTrip(ctx, trip)
	if err != nil {
		return uuid.Nil, err
	}

	if err := uc.worker.EnqueueTripEmbedding(id); err != nil {
		fmt.Printf("Warning: Could not enqueue task: %v\n", err)
	}

	return id, nil
}
