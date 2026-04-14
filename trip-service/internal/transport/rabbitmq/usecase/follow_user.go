package usecase

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type FollowUserUseCase interface {
	Execute(ctx context.Context, follower, following uuid.UUID, status string) error
}
type followUserUseCase struct {
	repository domain.TripRepository
}

func NewFollowUserUseCase(repository domain.TripRepository) FollowUserUseCase {
	return &followUserUseCase{
		repository: repository,
	}
}

func (uc *followUserUseCase) Execute(ctx context.Context, follower, following uuid.UUID, status string) error {
	return uc.repository.UpsertLocalFollow(ctx, follower, following, status)
}
