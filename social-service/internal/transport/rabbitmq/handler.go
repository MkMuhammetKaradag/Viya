package rabbitmq

import (
	"social-service/internal/domain"
	controller "social-service/internal/transport/rabbitmq/controller"
	"social-service/internal/transport/rabbitmq/usecase"
	"viya/pkg/messaging"
)

type Handlers struct {
	UserCreated domain.MessageHandler
}

func NewMessageHandlers(repo domain.SocialRepository) *Handlers {
	return &Handlers{

		UserCreated: controller.NewUserCreatedHandler(
			usecase.NewUserCreatedUseCase(repo),
		),
	}
}

func SetupMessageHandlers(repo domain.SocialRepository) map[messaging.MessageType]domain.MessageHandler {
	h := NewMessageHandlers(repo)

	return map[messaging.MessageType]domain.MessageHandler{
		messaging.AuthTypes.CreatedUser: h.UserCreated,
	}
}
