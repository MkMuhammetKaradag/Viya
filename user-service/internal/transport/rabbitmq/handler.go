package rabbitmq

import (
	"user-service/internal/domain"
	controller "user-service/internal/transport/rabbitmq/handler"
	"user-service/internal/transport/rabbitmq/usecase"
	"viya/pkg/messaging"
)

type Handlers struct {
	UserCreated domain.MessageHandler
}

func NewMessageHandlers(repo domain.UserRepository) *Handlers {
	return &Handlers{

		UserCreated: controller.NewUserCreatedHandler(
			usecase.NewUserCreatedUseCase(repo),
		),
	}
}

func SetupMessageHandlers(repo domain.UserRepository) map[messaging.MessageType]domain.MessageHandler {
	h := NewMessageHandlers(repo)

	return map[messaging.MessageType]domain.MessageHandler{
		messaging.AuthTypes.CreatedUser: h.UserCreated,
	}
}
