package messaging

import (
	"context"
	"encoding/json"
	
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitClient, bağlantıyı kurar ve kanalı açar.
func NewRabbitClient(url string) (*RabbitClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitClient{conn: conn, channel: ch}, nil
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