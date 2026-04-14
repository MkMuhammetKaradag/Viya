package usecase

import (
	"context"
	"fmt"
	"social-service/internal/domain"

	"github.com/google/uuid"
)

type UnFollowUserUseCase interface {
	Execute(ctx context.Context, followerID, targetUserID uuid.UUID) error
}

type UnfollowUserUseCase struct {
	repo         domain.SocialRepository
	rabbitClient domain.RabbitMQClient
}

func NewUnFollowUserUseCase(repo domain.SocialRepository, rabbitClient domain.RabbitMQClient) UnFollowUserUseCase {
	return &UnfollowUserUseCase{repo: repo, rabbitClient: rabbitClient}
}

func (uc *UnfollowUserUseCase) Execute(ctx context.Context, followerID, targetUserID uuid.UUID) error {
	if followerID == targetUserID {
		return fmt.Errorf("Invalid transaction")
	}

	return uc.repo.UnfollowUser(ctx, followerID, targetUserID)
}
