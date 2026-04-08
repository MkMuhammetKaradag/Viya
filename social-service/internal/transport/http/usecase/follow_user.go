package usecase

import (
	"context"
	"social-service/internal/domain"

	"github.com/google/uuid"
)

type FollowUserUseCase interface {
	Execute(ctx context.Context, followerID, targetUserID uuid.UUID) (string, error)
}

type followUserUseCase struct {
	repo domain.SocialRepository
}

func NewFollowUserUseCase(repo domain.SocialRepository) FollowUserUseCase {
	return &followUserUseCase{repo: repo}
}

func (uc *followUserUseCase) Execute(ctx context.Context, followerID, targetUserID uuid.UUID) (string, error) {
	status, err := uc.repo.CreatwFollow(ctx, followerID, targetUserID)
	return status, err
}
