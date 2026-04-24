package usecase

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type GetTripCommentsUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, tripID uuid.UUID, limit, page int) ([]domain.Comment, error)
}

type getTripCommentsUseCase struct {
	tripRepo domain.TripRepository
}

func NewGetTripCommentsUseCase(tripRepo domain.TripRepository) GetTripCommentsUseCase {
	return &getTripCommentsUseCase{tripRepo: tripRepo}
}

func (uc *getTripCommentsUseCase) Execute(ctx context.Context, userID, tripID uuid.UUID, limit, page int) ([]domain.Comment, error) {
	comments, err := uc.tripRepo.GetTripComments(ctx, userID, tripID, page, limit)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
