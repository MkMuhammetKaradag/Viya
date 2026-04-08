package usecase

import (
	"context"
	"fmt"
	"social-service/internal/domain"

	"github.com/google/uuid"
)

type FollowRequestUseCase interface {
	Execute(ctx context.Context, myID, followerID uuid.UUID, action string) error
}

type followRequestUseCase struct {
	repo domain.SocialRepository
}

func NewFollowRequestUseCase(repo domain.SocialRepository) FollowRequestUseCase {
	return &followRequestUseCase{repo: repo}
}

func (uc *followRequestUseCase) Execute(ctx context.Context, myID, followerID uuid.UUID, action string) error {
	if action == "ACCEPT" {
		return uc.repo.UpdateFollowStatus(ctx, followerID, myID, "ACCEPTED")
	} else if action == "REJECT" {
		return uc.repo.UpdateFollowStatus(ctx, followerID, myID, "REJECTED")
	}
	return fmt.Errorf("geçersiz işlem")
}
