package usecase

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type ForkTripUseCase interface {
	Execute(ctx context.Context, tripID, userID uuid.UUID) (uuid.UUID, error)
}

type forkTripUseCase struct {
	tripRepo domain.TripRepository
	worker   domain.Worker
}

func NewForkTripUseCase(tripRepo domain.TripRepository, worker domain.Worker) ForkTripUseCase {
	return &forkTripUseCase{tripRepo: tripRepo, worker: worker}
}

func (uc *forkTripUseCase) Execute(ctx context.Context, tripID, userID uuid.UUID) (uuid.UUID, error) {
	tripId, err := uc.tripRepo.ForkTrip(ctx, tripID, userID)
	if err != nil {
		return uuid.Nil, err
	}
	// if !trip.IsPublic && trip.UserID != userID {
	// 	// Kullanıcı sahibi değilse ve gezi gizliyse 404
	// 	return nil, fmt.Errorf("trip not found or access denied")
	// }

	// if err := uc.worker.EnqueueIncrementTrip(tripID, userID, 0.1, "fork"); err != nil {
	// 	fmt.Printf("Warning: Could not enqueue task: %v\n", err)
	// }

	return tripId, nil
}
