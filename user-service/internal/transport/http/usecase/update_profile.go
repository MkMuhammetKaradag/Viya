package usecase

import (
	"context"
	"log"
	"user-service/internal/domain"
	"viya/pkg/messaging"

	"github.com/google/uuid"
)

type UpdateProfileUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, params domain.UpdateProfileParams) error
}
type updateProfileUseCase struct {
	userRepository domain.UserRepository
	rabbitClient   domain.RabbitMQClient
}

func NewUpdateProfileUseCase(userRepository domain.UserRepository, rabbitClient domain.RabbitMQClient) UpdateProfileUseCase {
	return &updateProfileUseCase{
		userRepository: userRepository,
		rabbitClient:   rabbitClient,
	}
}

func (uc *updateProfileUseCase) Execute(ctx context.Context, userID uuid.UUID, params domain.UpdateProfileParams) error {

	err := uc.userRepository.UpdateProfile(ctx, userID, params)
	if err != nil {
		return err
	}
	if params.IsPrivate != nil {

		updatedMessage := messaging.Message{
			Type: messaging.UserTypes.UpdatedUser,
			ToServices: []messaging.ServiceType{
				messaging.SocialService,
				// messaging.TripService,
			},
			Data: map[string]interface{}{
				"id":         userID,
				"is_private": params.IsPrivate,
			},
			Critical: true,
		}

		// 3. RabbitMQ üzerinden yayınla
		err = uc.rabbitClient.PublishMessage(ctx, updatedMessage)
		if err != nil {
			log.Printf("User update message could not be sent: %v", err)
		}
	}
	return nil
}
