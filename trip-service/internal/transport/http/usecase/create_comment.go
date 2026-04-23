package usecase

import (
	"context"
	"fmt"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type CreateCommentUseCase interface {
	Execute(ctx context.Context, comment *domain.Comment) (uuid.UUID, error)
}

type createCommentUseCase struct {
	tripRepo   domain.TripRepository
	moderation domain.ModerationService
}

func NewCreateCommentUseCase(tripRepo domain.TripRepository, moderation domain.ModerationService) CreateCommentUseCase {
	return &createCommentUseCase{tripRepo: tripRepo, moderation: moderation}
}

func (uc *createCommentUseCase) Execute(ctx context.Context, comment *domain.Comment) (uuid.UUID, error) {

	result, err := uc.moderation.Moderate(ctx, comment.Content)
	if err != nil {
		return uuid.Nil, err
	}
	if !result.IsAppropriate {
		return uuid.Nil, fmt.Errorf("comment rejected: %s", result.Reason)
	}
	id, err := uc.tripRepo.CreateComment(ctx, comment)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
