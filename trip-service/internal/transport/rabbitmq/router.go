package rabbitmq

import (
	"log"
	"trip-service/internal/domain"
	"viya/pkg/messaging"
)

type RabbitRouter struct {
	handlers map[messaging.MessageType]domain.MessageHandler
}

func NewRabbitRouter(repo domain.TripRepository) *RabbitRouter {
	return &RabbitRouter{
		handlers: SetupMessageHandlers(repo),
	}
}

func (r *RabbitRouter) Route(msg messaging.Message) error {
	handler, ok := r.handlers[msg.Type]
	if !ok {
		log.Printf("[WARN] No handler registered for message type: %s", msg.Type)
		return nil
	}
	return handler.Handle(msg)
}
