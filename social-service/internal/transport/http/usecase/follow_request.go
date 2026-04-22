package usecase

import (
	"context"
	"fmt"
	"log"
	"social-service/internal/domain"
	"viya/pkg/messaging"

	"github.com/google/uuid"
)

type FollowRequestUseCase interface {
	Execute(ctx context.Context, myID, followerID uuid.UUID, action string) error
}

type followRequestUseCase struct {
	repo         domain.SocialRepository
	rabbitClient domain.RabbitMQClient
}

func NewFollowRequestUseCase(repo domain.SocialRepository, rabbitClient domain.RabbitMQClient) FollowRequestUseCase {
	return &followRequestUseCase{repo: repo, rabbitClient: rabbitClient}
}

func (uc *followRequestUseCase) Execute(ctx context.Context, myID, followerID uuid.UUID, action string) error {
	if action == "ACCEPT" {
		err := uc.repo.UpdateFollowStatus(ctx, followerID, myID, "ACCEPTED")
		if err != nil {
			return err
		}
		updatedMessage := messaging.Message{
			Type: messaging.SocialTypes.FollowUser,
			ToServices: []messaging.ServiceType{
				messaging.TripService,
				messaging.UserService,
			},
			Data: map[string]interface{}{
				"follower":  followerID,
				"following": myID,
				"status":    "ACCEPTED",
			},
			Critical: true,
		}

		// 3. RabbitMQ üzerinden yayınla
		err = uc.rabbitClient.PublishMessage(ctx, updatedMessage)
		if err != nil {
			log.Printf("User update message could not be sent: %v", err)
		}

		return nil
	} else if action == "REJECT" {
		return uc.repo.UpdateFollowStatus(ctx, followerID, myID, "REJECTED")
	}
	return fmt.Errorf("geçersiz işlem")
}
