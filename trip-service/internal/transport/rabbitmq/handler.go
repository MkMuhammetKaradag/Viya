package rabbitmq

import (
	"trip-service/internal/domain"
	controller "trip-service/internal/transport/rabbitmq/controller"
	"trip-service/internal/transport/rabbitmq/usecase"
	"viya/pkg/messaging"
)

type Handlers struct {
	UserCreated domain.MessageHandler
}

func NewMessageHandlers(repo domain.TripRepository) *Handlers {
	return &Handlers{

		UserCreated: controller.NewUserCreatedHandler(
			usecase.NewUserCreatedUseCase(repo),
		),
	}
}

func SetupMessageHandlers(repo domain.TripRepository) map[messaging.MessageType]domain.MessageHandler {
	h := NewMessageHandlers(repo)

	return map[messaging.MessageType]domain.MessageHandler{
		messaging.AuthTypes.CreatedUser: h.UserCreated,
	}
}
