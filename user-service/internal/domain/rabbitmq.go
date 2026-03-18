package domain

import (
	"context"
	"viya/pkg/messaging"
)

type MessageProducer interface {
	PublishMessage(ctx context.Context, msg messaging.Message) error
	Close() error
}

// Sadece mesaj dinleyebilenler için
type MessageConsumer interface {
	ConsumeMessages(handler messaging.MessageHandler) error
	Close() error
}

// Her ikisini de yapabilenler için (Mevcut yapın)
type RabbitMQClient interface {
	MessageProducer
	MessageConsumer
}
type MessageHandler interface {
	Handle(msg messaging.Message) error
}
type RabbitRouter interface {
	Route(msg messaging.Message) error
}
