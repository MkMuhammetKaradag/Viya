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
}

func NewGetTripUseCase(tripRepo domain.TripRepository) GetTripUseCase {
	return &getTripUseCase{tripRepo: tripRepo}
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

	go func() {
		// Tekil görüntülenme mantığını
		_ = uc.tripRepo.IncrementUniqueView(context.Background(), tripID, userID)
	}()

	return trip, nil
}
