package messaging

import "fmt"

type MessagingError struct {
	config    Config
	Code    string
	Message string
	Err     error
}

func (e *MessagingError) Error() string {
	return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
}

var (
	ErrConnectionFailed = &MessagingError{Code: "CONNECTION_FAILED", Message: "Failed to connect to RabbitMQ"}
	ErrPublishFailed    = &MessagingError{Code: "PUBLISH_FAILED", Message: "Failed to publish message"}
	ErrConsumeFailed    = &MessagingError{Code: "CONSUME_FAILED", Message: "Failed to consume messages"}
)
