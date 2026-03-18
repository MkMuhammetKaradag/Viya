package rabbitmq

import (
	"log"
	"user-service/internal/domain"
	controller "user-service/internal/transport/rabbitmq/controller"
	"user-service/internal/transport/rabbitmq/usecase"
	"viya/pkg/messaging"
)

type RabbitRouter struct {
	handlers map[messaging.MessageType]domain.MessageHandler
}

func NewRabbitRouter(repo domain.UserRepository) *RabbitRouter {
	// Handler'ları burada kaydet
	userCreatedUseCase := usecase.NewUserCreatedUseCase(repo)
	userCreatedHandler := controller.NewUserCreatedHandler(userCreatedUseCase)

	return &RabbitRouter{
		handlers: map[messaging.MessageType]domain.MessageHandler{
			messaging.AuthTypes.CreatedUser: userCreatedHandler,
		},
	}
}

// Route metodu, gelen mesajı doğru handler'a iletir
func (r *RabbitRouter) Route(msg messaging.Message) error {
	handler, ok := r.handlers[msg.Type]
	if !ok {
		// Log basıp geçebiliriz, hata döndürmek mesajın retry'a girmesine neden olabilir
		log.Printf("[WARN] No handler registered for message type: %s", msg.Type)
		return nil
	}
	return handler.Handle(msg)
}
