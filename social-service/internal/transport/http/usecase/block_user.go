package usecase

import (
	"context"
	"social-service/internal/domain"

	"github.com/google/uuid"
)

type BlockUserUseCase interface {
	Execute(ctx context.Context, BlockerID, targetUserID uuid.UUID) error
}

type blockUserUseCase struct {
	repo         domain.SocialRepository
	rabbitClient domain.RabbitMQClient
}

func NewBlockUserUseCase(repo domain.SocialRepository, rabbitClient domain.RabbitMQClient) BlockUserUseCase {
	return &blockUserUseCase{repo: repo, rabbitClient: rabbitClient}
}

func (uc *blockUserUseCase) Execute(ctx context.Context, BlockerID, targetUserID uuid.UUID) error {
	err := uc.repo.BlockUser(ctx, BlockerID, targetUserID)
	return err
}
