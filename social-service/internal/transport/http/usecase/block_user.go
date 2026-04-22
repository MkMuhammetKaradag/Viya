package usecase

import (
	"context"
	"log"
	"social-service/internal/domain"
	"viya/pkg/messaging"

	"github.com/google/uuid"
)

type BlockUserUseCase interface {
	Execute(ctx context.Context, BlockerID, targetUserID uuid.UUID) error
}

type blockUserUseCase struct {
	repo         domain.SocialRepository
	rabbitClient domain.RabbitMQClient
}

func NewBlockUserUseCase(repo domain.SocialRepository, rabbitClient domain.RabbitMQClient) BlockUserUseCase {
	return &blockUserUseCase{repo: repo, rabbitClient: rabbitClient}
}

func (uc *blockUserUseCase) Execute(ctx context.Context, BlockerID, targetUserID uuid.UUID) error {
	err := uc.repo.BlockUser(ctx, BlockerID, targetUserID)
	if err != nil {
		return err
	}

	updatedMessage := messaging.Message{
		Type: messaging.SocialTypes.BlockUser,
		ToServices: []messaging.ServiceType{
			messaging.TripService,
			messaging.UserService,
		},
		Data: map[string]interface{}{
			"blocker": BlockerID,
			"blocked": targetUserID,
		},
		Critical: true,
	}

	err = uc.rabbitClient.PublishMessage(ctx, updatedMessage)
	if err != nil {
		log.Printf("User update message could not be sent: %v", err)
	}
	return nil

}
