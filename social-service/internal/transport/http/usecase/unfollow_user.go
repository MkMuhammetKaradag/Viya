package usecase

import (
	"context"
	"fmt"
	"log"
	"social-service/internal/domain"
	"viya/pkg/messaging"

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

	err := uc.repo.UnfollowUser(ctx, followerID, targetUserID)
	if err != nil {
		return fmt.Errorf("Failed to unfollow user: %w", err)
	}
	updatedMessage := messaging.Message{
		Type: messaging.SocialTypes.UnFollowUser,
		ToServices: []messaging.ServiceType{
			messaging.TripService,
			messaging.UserService,
		},
		Data: map[string]interface{}{
			"follower":  followerID,
			"following": targetUserID,
		},
		Critical: true,
	}

	// 3. RabbitMQ üzerinden yayınla
	err = uc.rabbitClient.PublishMessage(ctx, updatedMessage)
	if err != nil {
		log.Printf("User update message could not be sent: %v", err)
	}
	return nil
}
