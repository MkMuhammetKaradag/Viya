package usecase

import (
	"context"
	"fmt"
	"trip-service/internal/domain"
)

type UpdateWaypointUseCase interface {
	Execute(ctx context.Context, wp *domain.UpdateWaypointInput) error
}

type updateWaypointUseCase struct {
	repo domain.TripRepository
}

func NewUpdateWaypointUseCase(repo domain.TripRepository) UpdateWaypointUseCase {
	return &updateWaypointUseCase{
		repo: repo,
	}
}

func (u *updateWaypointUseCase) Execute(ctx context.Context, wp *domain.UpdateWaypointInput) error {
	existingWP, err := u.repo.GetWaypointByID(ctx, wp.WayPointID)
	if err != nil {
		return err
	}
	if existingWP == nil {
		return fmt.Errorf("waypoint not found")
	}
	if wp.Title != nil {
		existingWP.Title = *wp.Title
	}
	if wp.Description != nil {
		existingWP.Description = *wp.Description
	}
	if wp.Latitude != nil {
		existingWP.Latitude = *wp.Latitude
	}
	if wp.Longitude != nil {
		existingWP.Longitude = *wp.Longitude
	}

	return u.repo.UpdateWaypoint(ctx, existingWP)
}
