package usecase

import (
	"context"
	"fmt"
	"social-service/internal/domain"

	"github.com/google/uuid"
)

type RemoveFollowerUseCase interface {
	Execute(ctx context.Context, myID, followerID uuid.UUID) error
}

type removeFollowerUseCase struct {
	repo domain.SocialRepository
}

func NewRemoveFollowerUseCase(repo domain.SocialRepository) RemoveFollowerUseCase {
	return &removeFollowerUseCase{repo: repo}
}

func (uc *removeFollowerUseCase) Execute(ctx context.Context, myID, followerID uuid.UUID) error {
	if myID == followerID {
		return fmt.Errorf("Invalid transaction")
	}

	return uc.repo.RemoveFollower(ctx, myID, followerID)
}
