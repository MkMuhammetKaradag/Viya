package domain

type RabbitMQClient interface {
	Consume(queueName string, handler func(data []byte) error)
	Publish(queueName string, body interface{}) error
}
