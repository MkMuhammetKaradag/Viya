package messaging

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	config    Config
	mu        sync.Mutex
	service   ServiceType
	conn      *amqp.Connection
	channel   *amqp.Channel
	reconnect chan bool
	closed    bool
}

// NewRabbitClient, bağlantıyı kurar ve kanalı açar.
func NewRabbitClient(config Config, serviceType ServiceType) (*RabbitClient, error) {
	r := &RabbitClient{
		config:    config,
		service:   serviceType,
		reconnect: make(chan bool),
	}
	if err := r.connect(serviceType); err != nil {
		return nil, err
	}
	go r.monitorConnection(serviceType)
	return r, nil
}

// Publish, senin Kafka'daki gibi mesajı JSON'a çevirip kuyruğa atar.
func (r *RabbitClient) Publish(queueName string, body interface{}) error {
	// Kuyruğun var olduğundan emin ol (Durable: true yaparak mesaj kaybını önleriz)
	q, err := r.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(body)
	return r.channel.PublishWithContext(context.Background(), "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent, // Disk'e kaydet (Kafka'daki RequiredAcks gibi)
		Body:         data,
	})
}

// Consume, mesajları dinler ve senin verdiğin handler fonksiyonuna paslar.
func (r *RabbitClient) Consume(queueName string, handler func(data []byte) error) {
	q, _ := r.channel.QueueDeclare(queueName, true, false, false, false, nil)

	msgs, _ := r.channel.Consume(q.Name, "", false, false, false, false, nil)

	go func() {
		for d := range msgs {
			err := handler(d.Body)
			if err == nil {
				d.Ack(false) // Başarılıysa onayla
			} else {
				// Hata varsa 5 saniye bekleyip tekrar denemesi için Nack yapabilirsin
				time.Sleep(5 * time.Second)
				d.Nack(false, true)
			}
		}
	}()
}

func (r *RabbitClient) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.closed = true // monitorConnection'ın durması için

	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			log.Printf("RabbitMQ channel close error: %v", err)
		}
	}

	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			return err
		}
	}

	return nil
}
