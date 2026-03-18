package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitClient) ConsumeMessages(handler MessageHandler) error {
	queueName := string(r.service) + ".queue"

	// Ana kuyruk tanımlama
	q, err := r.channel.QueueDeclare(
		queueName,
		r.config.QueueDurable,
		r.config.QueueAutoDelete,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange": DLQExchangeName,
		},
	)
	if err != nil {
		return &MessagingError{Code: "QUEUE_FAILED", Message: "Failed to declare queue", Err: err}
	}

	// Broadcast mesajlar için genel exchange'e bağlan
	err = r.channel.QueueBind(
		q.Name,
		"",
		r.config.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		return &MessagingError{Code: "BIND_FAILED", Message: "Failed to bind queue to broadcast exchange", Err: err}
	}

	// Servise özel mesajlar için direct exchange'e bağlan
	serviceExchangeName := fmt.Sprintf("dc_clone.%s.service", string(r.service))
	err = r.channel.QueueBind(
		q.Name,
		"",
		serviceExchangeName,
		false,
		nil,
	)
	if err != nil {
		return &MessagingError{Code: "BIND_FAILED", Message: "Failed to bind queue to service exchange", Err: err}
	}

	// Retry kuyruğunu da deklare et ve bağla
	if r.config.EnableRetry {
		retryQueueName := string(r.service) + ".retry.queue"

		// Retry kuyruğunu tanımla
		_, err = r.channel.QueueDeclare(
			retryQueueName,
			true,
			false,
			false,
			false,
			amqp.Table{
				"x-dead-letter-exchange":    serviceExchangeName, // Direkt servise gönder
				"x-dead-letter-routing-key": "",                  // Direct exchange için boş routing key
			},
		)
		if err != nil {
			return &MessagingError{Code: "RETRY_QUEUE_FAILED", Message: "Failed to declare retry queue", Err: err}
		}

		// Retry kuyruğunu servis adına göre bağla
		err = r.channel.QueueBind(
			retryQueueName,
			string(r.service), // Sadece bu servis için olan retry mesajlarını al
			r.config.RetryExchangeName,
			false,
			nil,
		)
		if err != nil {
			return &MessagingError{Code: "RETRY_BIND_FAILED", Message: "Failed to bind retry queue", Err: err}
		}
	}

	// Ana kuyruğu dinle
	msgs, err := r.channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return &MessagingError{Code: "CONSUME_FAILED", Message: "Failed to start consuming", Err: err}
	}

	go func() {
		for msg := range msgs {
			var message Message
			if err := json.Unmarshal(msg.Body, &message); err != nil {
				log.Printf("[ERROR] JSON Unmarshal Failed! Body: %s, Error: %v", string(msg.Body), err)
				r.publishToDLQ(msg.Body, err.Error())

				msg.Ack(false)
				// msg.Nack(false, false) // DLQ'ya gönder
				continue
			}

			// Mesaj bu servise ait değilse işleme
			// if message.ToService != "" && message.ToService != r.service {
			// 	log.Printf("Bu mesaj %s servisi için, atlıyoruz: %s", message.ToService, message.ID)
			// 	msg.Ack(false) // Mesajı kabul et ama işleme
			// 	continue
			// }

			if len(message.ToServices) > 0 {
				isForThisService := false
				for _, svc := range message.ToServices {
					if svc == r.service {
						isForThisService = true
						break
					}
				}

				if !isForThisService {
					log.Printf("This message is for %v services, we skip it because this service (%s) is not in the list: %s",
						message.ToServices, r.service, message.ID)
					msg.Ack(false) // Mesajı kabul et ama işleme
					continue
				}
			}

			log.Printf("Processing message [ID: %s, Type: %s, RetryCount: %d, ToServices: %v]",
				message.ID, message.Type, message.RetryCount, message.ToServices)

			err := handler(message)

			if err != nil {
				log.Printf("Message processing failed: %v", err)
				if message.Critical {
					// Kritik mesajlar için retry sayısını dikkate almadan tekrar dene
					r.handleCriticalMessageRetry(&message)
					msg.Ack(false) // Orijinal mesajı kabul et
				} else if r.shouldRetry(message) {
					log.Printf("Scheduling retry for message ID: %s", message.ID)
					r.handleRetry(&message)
					msg.Ack(false) // Orijinal mesajı kabul et, retry kuyruğunda yeni bir kopya var
				} else {
					log.Printf("Message failed permanently, sending to DLQ. ID: %s", message.ID)
					msg.Nack(false, false) // DLQ'ya gönder
				}
			} else {
				log.Printf("Message processed successfully. ID: %s", message.ID)
				msg.Ack(false)
			}
		}
	}()

	return nil
}
func (r *RabbitClient) shouldRetry(msg Message) bool {
	// Retry özelliği aktif değilse, hiç deneme yapma
	if !r.config.EnableRetry {
		return false
	}

	// Mesaj tipi retry listesinde mi kontrol et
	isRetryableType := false
	for _, t := range r.config.RetryTypes {
		if t == msg.Type { // Burada Type kontrolü yapılıyor
			isRetryableType = true
			break
		}
	}

	// Retry tipi ve sayısı uygun mu?
	if isRetryableType && msg.RetryCount < r.config.MaxRetries {
		log.Printf("Message will be retried. Current retry count: %d, Max retries: %d",
			msg.RetryCount, r.config.MaxRetries)
		return true
	}

	return false
}
func (r *RabbitClient) handleRetry(msg *Message) {
	// Retry sayısını artır
	msg.RetryCount++

	// Retry delay hesapla
	retryDelay := 5000 * msg.RetryCount

	body, err := json.Marshal(msg)
	if err != nil {
		log.Printf("handleRetry marshal error: %v", err)
		return
	}

	// Tek bir servise retry yapılıyor varsayımıyla
	var routingKey string
	if len(msg.ToServices) > 0 {
		routingKey = string(msg.ToServices[0])
	} else {
		routingKey = string(r.service)
	}

	log.Printf("Retry mesajı gönderiliyor. ID: %s, ToServices: %v, RoutingKey: %s",
		msg.ID, msg.ToServices, routingKey)

	err = r.channel.Publish(
		r.config.RetryExchangeName,
		routingKey, // Önemli: Mesajın sadece hedef servise gitmesi için routing key
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			MessageId:    msg.ID,
			Timestamp:    time.Now(),
			DeliveryMode: 2,
			Headers:      amqp.Table(msg.Headers),
			Expiration:   fmt.Sprintf("%d", retryDelay),
		},
	)

	if err != nil {
		log.Printf("handleRetry publish error: %v", err)
	} else {
		log.Printf("Message sent to retry queue for %s with %d seconds delay",
			routingKey, retryDelay/1000)
	}
}

func (r *RabbitClient) publishToDLQ(originalBody []byte, reason string) {
	// DLQ için özel header'lar oluştur
	headers := amqp.Table{
		"x-original-error": reason,
		"x-failed-at":      time.Now().Format(time.RFC3339),
		"x-service":        string(r.service),
	}

	err := r.channel.PublishWithContext(context.Background(),
		DLQExchangeName, // Senin daha önce tanımladığın DLQ Exchange
		"",              // Fanout olduğu için boş
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         originalBody,
			Headers:      headers,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		log.Printf("[CRITICAL] Failed to move broken message to DLQ: %v", err)
	} else {
		log.Printf("[INFO] Broken message successfully moved to DLQ with error header")
	}
}
