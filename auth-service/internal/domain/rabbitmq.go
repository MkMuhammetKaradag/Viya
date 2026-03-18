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

// type RabbitMQClient interface {
// 	Consume(queueName string, handler func(data []byte) error)
// 	Publish(queueName string, body interface{}) error
// 	PublishMessage(ctx context.Context, msg messaging.Message) error
// 	ConsumeMessages(handler messaging.MessageHandler) error
// 	Close() error
// }
