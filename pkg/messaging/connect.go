package messaging

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	DLQExchangeName = "dead_letter.exchange"
	DLQName         = "dead_letter.queue"
)

func (r *RabbitClient) connect(serviceType ServiceType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	conn, err := amqp.DialConfig(r.config.GetAMQPURL(), amqp.Config{
		Heartbeat: 10 * time.Second,
		Dial:      amqp.DefaultDial(r.config.ConnectionTimeout),
	})
	if err != nil {
		return &MessagingError{Code: "CONNECTION_FAILED", Message: "Failed to connect", Err: err}
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return &MessagingError{Code: "CHANNEL_FAILED", Message: "Failed to create channel", Err: err}
	}

	if err := r.setupExchanges(ch, serviceType); err != nil {
		ch.Close()
		conn.Close()
		return err
	}

	r.conn = conn
	r.channel = ch
	r.closed = false

	return nil
}

func (r *RabbitClient) setupExchanges(ch *amqp.Channel, serviceType ServiceType) error {

	err := ch.ExchangeDeclare(
		r.config.ExchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	if serviceType != "" {
		serviceExchangeName := fmt.Sprintf("dc_clone.%s.service", serviceType)
		err = ch.ExchangeDeclare(
			serviceExchangeName,
			"direct",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

	}

	if r.config.EnableRetry {
		// Retry exchange
		serviceExchangeName := fmt.Sprintf("dc_clone.%s.service", serviceType)
		err = ch.ExchangeDeclare(
			r.config.RetryExchangeName,
			"direct",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		// Tek bir retry kuyruğu oluştur
		retryQueueName := string(serviceType) + ".retry.queue"
		_, err = ch.QueueDeclare(
			retryQueueName,
			true,
			false,
			false,
			false,
			amqp.Table{
				"x-dead-letter-exchange":    serviceExchangeName, // retry sonrası ana kuyruğa
				"x-dead-letter-routing-key": "",                  //string(serviceType)
			},
		)
		if err != nil {
			return err
		}

		// Bind retry queue to retry exchange
		err = ch.QueueBind(
			retryQueueName,
			string(serviceType), // Tek bir routing key kullan
			r.config.RetryExchangeName,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	err = ch.ExchangeDeclare(
		DLQExchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// DLQ queue
	_, err = ch.QueueDeclare(
		DLQName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// DLQ queue bind
	err = ch.QueueBind(
		DLQName,
		"",
		DLQExchangeName,
		false,
		nil,
	)
	return err

}
