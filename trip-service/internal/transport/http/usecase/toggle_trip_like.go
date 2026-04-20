package usecase

import (
	"context"
	"errors"
	"fmt"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type ToggleTripLikeUseCase interface {
	Execute(ctx context.Context, tripID, userID uuid.UUID) (bool, error)
}

type toggleTriplikeUseCase struct {
	tripRepo domain.TripRepository
	worker   domain.Worker
}

func NewToggleTripLikeUseCase(tripRepo domain.TripRepository, worker domain.Worker) ToggleTripLikeUseCase {
	return &toggleTriplikeUseCase{tripRepo: tripRepo, worker: worker}
}

func (uc *toggleTriplikeUseCase) Execute(ctx context.Context, tripID, userID uuid.UUID) (bool, error) {
	trip, err := uc.tripRepo.GetTripStatus(ctx, tripID)
	if err != nil {
		return false, err
	}

	isOwner := trip.UserID == userID

	if !isOwner {

		if !trip.IsPublic {
			return false, errors.New("bu gezi gizli, beğenilemez")
		}

		if trip.OwnerIsPrivate {
			isFollowing, err := uc.tripRepo.CheckFollowStatus(ctx, userID, trip.UserID)
			if err != nil || !isFollowing {
				return false, errors.New("bu gizli kullanıcının gezisini beğenmek için onu takip etmelisin")
			}
		}
	}

	isLiked, err := uc.tripRepo.ToggleTripLike(ctx, tripID, userID)
	if err != nil {
		return false, fmt.Errorf("beğeni işlemi başarısız: %w", err)
	}
	var weight float32 = 0.5
	action := "like"

	if !isLiked {
		weight = -0.5
		action = "unlike"
	}
	if err := uc.worker.EnqueueIncrementTrip(tripID, userID, weight, action); err != nil {
		fmt.Printf("Warning: Could not enqueue task: %v\n", err)
	}

	return isLiked, nil

}
