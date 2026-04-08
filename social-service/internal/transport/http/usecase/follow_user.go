package usecase

import (
	"context"
	"fmt"
	"social-service/internal/domain"

	"github.com/google/uuid"
)

type FollowUserUseCase interface {
	Execute(ctx context.Context, followerID, targetUserID uuid.UUID) (string, error)
}

type followUserUseCase struct {
	repo domain.SocialRepository
}

func NewFollowUserUseCase(repo domain.SocialRepository) FollowUserUseCase {
	return &followUserUseCase{repo: repo}
}

func (uc *followUserUseCase) Execute(ctx context.Context, followerID, targetUserID uuid.UUID) (string, error) {
	if followerID == targetUserID {
		return "", fmt.Errorf("You can't keep track of yourself.")
	}

	// 2. Engel kontrolü (O beni engellemiş mi?)
	blockers, err := uc.repo.IsBlocked(ctx, targetUserID, followerID)
	if err != nil {
		return "", err
	}

	if len(blockers) > 0 {
		for _, bID := range blockers {
			if bID == targetUserID {
				return "", fmt.Errorf("user not found")
			}
		}

		for _, bID := range blockers {
			if bID == followerID {
				return "", fmt.Errorf("you have blocked this user")
			}
		}
	}

	status, err := uc.repo.CreateFollow(ctx, followerID, targetUserID)
	return status, err
}
