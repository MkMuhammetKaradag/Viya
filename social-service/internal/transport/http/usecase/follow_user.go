package usecase

import (
	"context"
	"fmt"
	"log"
	"social-service/internal/domain"
	"viya/pkg/messaging"

	"github.com/google/uuid"
)

type FollowUserUseCase interface {
	Execute(ctx context.Context, followerID, targetUserID uuid.UUID) (string, error)
}

type followUserUseCase struct {
	repo         domain.SocialRepository
	rabbitClient domain.RabbitMQClient
}

func NewFollowUserUseCase(repo domain.SocialRepository, rabbitClient domain.RabbitMQClient) FollowUserUseCase {
	return &followUserUseCase{repo: repo, rabbitClient: rabbitClient}
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
	updatedMessage := messaging.Message{
		Type: messaging.SocialTypes.FollowUser,
		ToServices: []messaging.ServiceType{
			messaging.TripService,
		},
		Data: map[string]interface{}{
			"follower":  followerID,
			"following": targetUserID,
			"status":    status,
		},
		Critical: true,
	}

	// 3. RabbitMQ üzerinden yayınla
	err = uc.rabbitClient.PublishMessage(ctx, updatedMessage)
	if err != nil {
		log.Printf("User update message could not be sent: %v", err)
	}
	return status, err
}
