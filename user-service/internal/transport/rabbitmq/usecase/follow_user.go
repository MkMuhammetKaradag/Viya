package usecase

import (
	"context"
	"user-service/internal/domain"

	"github.com/google/uuid"
)

type FollowUserUseCase interface {
	Execute(ctx context.Context, follower, following uuid.UUID, status string) error
}
type followUserUseCase struct {
	repository domain.UserRepository
}

func NewFollowUserUseCase(repository domain.UserRepository) FollowUserUseCase {
	return &followUserUseCase{
		repository: repository,
	}
}

func (uc *followUserUseCase) Execute(ctx context.Context, follower, following uuid.UUID, status string) error {
	return uc.repository.UpsertLocalFollow(ctx, follower, following, status)
}
