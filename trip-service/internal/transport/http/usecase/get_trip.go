package usecase

import (
	"context"
	"fmt"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type GetTripUseCase interface {
	Execute(ctx context.Context, tripID, userID uuid.UUID) (*domain.Trip, error)
}

type getTripUseCase struct {
	tripRepo domain.TripRepository
	worker   domain.Worker
}

func NewGetTripUseCase(tripRepo domain.TripRepository, worker domain.Worker) GetTripUseCase {
	return &getTripUseCase{tripRepo: tripRepo, worker: worker}
}

func (uc *getTripUseCase) Execute(ctx context.Context, tripID, userID uuid.UUID) (*domain.Trip, error) {
	trip, err := uc.tripRepo.GetTripWithWaypointsAndPhotos(ctx, tripID)
	if err != nil {
		return nil, err
	}
	if !trip.IsPublic && trip.UserID != userID {
		// Kullanıcı sahibi değilse ve gezi gizliyse 404
		return nil, fmt.Errorf("trip not found or access denied")
	}

	// go func() {
	// 	// Tekil görüntülenme mantığını
	// 	_ = uc.tripRepo.IncrementUniqueView(context.Background(), tripID, userID)
	// }()

	if err := uc.worker.EnqueueIncrementTripView(tripID, userID); err != nil {
		fmt.Printf("Warning: Could not enqueue task: %v\n", err)
	}

	return trip, nil
}
