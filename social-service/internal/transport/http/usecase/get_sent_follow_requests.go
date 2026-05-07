package usecase

import (
	"context"
	"social-service/internal/domain"

	"github.com/google/uuid"
)

type GetSentFollowRequestsUseCase interface {
	Execute(ctx context.Context, myID uuid.UUID) ([]domain.PendingRequest, error)
}

type getSentFollowrequestsUseCase struct {
	repo domain.SocialRepository
}

func NewGetSentFollowRequestsUseCase(repo domain.SocialRepository) GetSentFollowRequestsUseCase {
	return &getSentFollowrequestsUseCase{repo: repo}
}

func (uc *getSentFollowrequestsUseCase) Execute(ctx context.Context, myID uuid.UUID) ([]domain.PendingRequest, error) {
	return uc.repo.GetSentFollowRequests(ctx, myID)
}
