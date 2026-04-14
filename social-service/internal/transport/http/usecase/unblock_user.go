package usecase

import (
	"context"
	"social-service/internal/domain"

	"github.com/google/uuid"
)

type UnblockUserUseCase interface {
	Execute(ctx context.Context, myID, targetUserID uuid.UUID) error
}

type unblockUserUseCase struct {
	repo         domain.SocialRepository
	rabbitClient domain.RabbitMQClient
}

func NewUnblockUserUseCase(repo domain.SocialRepository, rabbitClient domain.RabbitMQClient) UnblockUserUseCase {
	return &unblockUserUseCase{repo: repo, rabbitClient: rabbitClient}
}

func (uc *unblockUserUseCase) Execute(ctx context.Context, myID, targetUserID uuid.UUID) error {
	err := uc.repo.UnblockUser(ctx, myID, targetUserID)
	return err
}
