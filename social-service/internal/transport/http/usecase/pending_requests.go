package usecase

import (
	"context"
	"social-service/internal/domain"

	"github.com/google/uuid"
)

type PendingRequestsUseCase interface {
	Execute(ctx context.Context, myID uuid.UUID) ([]domain.PendingRequest, error)
}

type pendingRequestsUseCase struct {
	repo domain.SocialRepository
}

func NewPendingRequestsUseCase(repo domain.SocialRepository) PendingRequestsUseCase {
	return &pendingRequestsUseCase{repo: repo}
}

func (uc *pendingRequestsUseCase) Execute(ctx context.Context, myID uuid.UUID) ([]domain.PendingRequest, error) {
	return uc.repo.GetPendingFollowRequests(ctx, myID)
}
