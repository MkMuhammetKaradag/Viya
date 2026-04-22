package usecase

import (
	"context"
	"user-service/internal/domain"

	"github.com/google/uuid"
)

type BlockUserUseCase interface {
	Execute(ctx context.Context, blocker, blocked uuid.UUID) error
}
type blockUserUseCase struct {
	repository domain.UserRepository
}

func NewBlockUserUseCase(repository domain.UserRepository) BlockUserUseCase {
	return &blockUserUseCase{
		repository: repository,
	}
}

func (uc *blockUserUseCase) Execute(ctx context.Context, blocker, blocked uuid.UUID) error {
	return uc.repository.UpsertLocalBlock(ctx, blocker, blocked)
}
