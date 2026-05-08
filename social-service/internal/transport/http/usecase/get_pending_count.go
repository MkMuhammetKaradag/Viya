package usecase

import (
	"context"
	"social-service/internal/domain"

	"github.com/google/uuid"
)

type GetPendingCountUseCase interface {
	Execute(ctx context.Context, myID uuid.UUID) (int, error)
}

type getPendingCountUseCase struct {
	repo domain.SocialRepository
}

func NewGetPendingCountUseCase(repo domain.SocialRepository) GetPendingCountUseCase {
	return &getPendingCountUseCase{repo: repo}
}

func (uc *getPendingCountUseCase) Execute(ctx context.Context, myID uuid.UUID) (int, error) {
	return uc.repo.GetPendingFollowRequestsCount(ctx, myID)
}
