package usecase

import (
	"context"
	"fmt"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type AddWayPointUseCase interface {
	Execute(ctx context.Context, wayPoint *domain.Waypoint) (uuid.UUID, error)
}

type addWayPointUseCase struct {
	tripRepo domain.TripRepository
}

func NewAddWayPointUseCase(tripRepo domain.TripRepository) AddWayPointUseCase {
	return &addWayPointUseCase{tripRepo: tripRepo}
}

func (uc *addWayPointUseCase) Execute(ctx context.Context, wayPoint *domain.Waypoint) (uuid.UUID, error) {
	wpID, err := uc.tripRepo.AddWaypoint(ctx, wayPoint)
	if err != nil {
		return uuid.Nil, err
	}
	fmt.Println("points:", wayPoint, "ID:", wpID)
	return wpID, nil
}
