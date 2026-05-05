package usecase

import (
	"context"
	"user-service/internal/domain"

	"github.com/google/uuid"
)

type UnFollowUserUseCase interface {
	Execute(ctx context.Context, follower, following uuid.UUID) error
}
type unfollowUserUseCase struct {
	repository domain.UserRepository
}

func NewUnFollowUserUseCase(repository domain.UserRepository) UnFollowUserUseCase {
	return &unfollowUserUseCase{
		repository: repository,
	}
}

func (uc *unfollowUserUseCase) Execute(ctx context.Context, follower, following uuid.UUID) error {
	return uc.repository.DeleteLocalFollow(ctx, follower, following)
}
