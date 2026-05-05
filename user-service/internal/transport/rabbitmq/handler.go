package rabbitmq

import (
	"user-service/internal/domain"
	controller "user-service/internal/transport/rabbitmq/controller"
	"user-service/internal/transport/rabbitmq/usecase"
	"viya/pkg/messaging"
)

type Handlers struct {
	UserCreated  domain.MessageHandler
	FollowUser   domain.MessageHandler
	UnFollowUser domain.MessageHandler
	BlockUser    domain.MessageHandler
}

func NewMessageHandlers(repo domain.UserRepository) *Handlers {
	return &Handlers{

		UserCreated: controller.NewUserCreatedHandler(
			usecase.NewUserCreatedUseCase(repo),
		),
		FollowUser: controller.NewFollowUserHandler(
			usecase.NewFollowUserUseCase(repo),
		),
		BlockUser: controller.NewBlockUserHandler(
			usecase.NewBlockUserUseCase(repo),
		),
		UnFollowUser: controller.NewUnFollowUserHandler(
			usecase.NewUnFollowUserUseCase(repo),
		),
	}
}

func SetupMessageHandlers(repo domain.UserRepository) map[messaging.MessageType]domain.MessageHandler {
	h := NewMessageHandlers(repo)

	return map[messaging.MessageType]domain.MessageHandler{
		messaging.AuthTypes.CreatedUser:    h.UserCreated,
		messaging.SocialTypes.FollowUser:   h.FollowUser,
		messaging.SocialTypes.BlockUser:    h.BlockUser,
		messaging.SocialTypes.UnFollowUser: h.UnFollowUser,
	}
}
