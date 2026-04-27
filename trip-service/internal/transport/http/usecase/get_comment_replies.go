package usecase

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type GetCommentRepliseUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, commentID uuid.UUID, limit, page int) ([]domain.Comment, error)
}

type getCommentRepliseUseCase struct {
	tripRepo domain.TripRepository
}

func NewGetCommentRepliseUseCase(tripRepo domain.TripRepository) GetCommentRepliseUseCase {
	return &getCommentRepliseUseCase{tripRepo: tripRepo}
}

func (uc *getCommentRepliseUseCase) Execute(ctx context.Context, userID, commentID uuid.UUID, limit, page int) ([]domain.Comment, error) {
	comments, err := uc.tripRepo.GetCommentReplies(ctx, commentID, page, limit)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
