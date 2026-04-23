package usecase

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type CreateCommentUseCase interface {
	Execute(ctx context.Context, comment *domain.Comment) (uuid.UUID, error)
}

type createCommentUseCase struct {
	tripRepo domain.TripRepository
	worker   domain.Worker
}

func NewCreateCommentUseCase(tripRepo domain.TripRepository) CreateCommentUseCase {
	return &createCommentUseCase{tripRepo: tripRepo}
}

func (uc *createCommentUseCase) Execute(ctx context.Context, comment *domain.Comment) (uuid.UUID, error) {
	id, err := uc.tripRepo.CreateComment(ctx, comment)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
