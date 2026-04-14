package rabbitmq

import (
	"trip-service/internal/domain"
	controller "trip-service/internal/transport/rabbitmq/controller"
	"trip-service/internal/transport/rabbitmq/usecase"
	"viya/pkg/messaging"
)

type Handlers struct {
	UserCreated domain.MessageHandler
	UserUpdated domain.MessageHandler
	FollowUser  domain.MessageHandler
	BlockUser   domain.MessageHandler
}

func NewMessageHandlers(repo domain.TripRepository) *Handlers {
	return &Handlers{

		UserCreated: controller.NewUserCreatedHandler(
			usecase.NewUserCreatedUseCase(repo),
		),
		UserUpdated: controller.NewUserUpdatedHandler(
			usecase.NewUserUpdateddUseCase(repo),
		),
		FollowUser: controller.NewFollowUserHandler(
			usecase.NewFollowUserUseCase(repo),
		),
		BlockUser: controller.NewBlockUserHandler(
			usecase.NewBlockUserUseCase(repo),
		),
	}
}

func SetupMessageHandlers(repo domain.TripRepository) map[messaging.MessageType]domain.MessageHandler {
	h := NewMessageHandlers(repo)

	return map[messaging.MessageType]domain.MessageHandler{
		messaging.AuthTypes.CreatedUser:  h.UserCreated,
		messaging.UserTypes.UpdatedUser:  h.UserUpdated,
		messaging.SocialTypes.FollowUser: h.FollowUser,
		messaging.SocialTypes.BlockUser:  h.BlockUser,
	}
}
