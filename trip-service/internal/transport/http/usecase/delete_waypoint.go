package usecase

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type DeleteWaypointUseCase interface {
	Execute(ctx context.Context, wayPointID uuid.UUID) error
}

type deleteWaypointUseCase struct {
	tripRepo domain.TripRepository
}

func NewDeleteWaypointUseCase(tripRepo domain.TripRepository) DeleteWaypointUseCase {
	return &deleteWaypointUseCase{tripRepo: tripRepo}
}

func (u *deleteWaypointUseCase) Execute(ctx context.Context, wayPointID uuid.UUID) error {

	err := u.tripRepo.DeleteWaypoint(ctx, wayPointID)
	if err != nil {
		return err
	}
	return nil
}
