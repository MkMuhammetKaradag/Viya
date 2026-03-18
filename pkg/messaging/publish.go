package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitClient) PublishMessage(ctx context.Context, msg Message) error {
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}

	if msg.Created.IsZero() {
		msg.Created = time.Now()
	}
	// fmt.Println(msg)

	msg.FromService = r.service

	// ToServices boş ise broadcast için normal exchange kullan
	if len(msg.ToServices) == 0 {
		body, err := json.Marshal(msg)
		if err != nil {
			return &MessagingError{Code: "MARSHAL_FAILED", Message: "Failed to marshal message", Err: err}
		}

		return r.channel.PublishWithContext(ctx,
			r.config.ExchangeName,
			"",
			true,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				MessageId:    msg.ID,
				Timestamp:    msg.Created,
				Priority:     uint8(msg.Priority),
				Headers:      amqp.Table(msg.Headers),
				DeliveryMode: 2,
			},
		)
	}

	// Birden fazla servise gönderim için
	var publishErrors []error
	var successServices []ServiceType
	// Her bir hedef servis için mesajı gönder
	for _, toService := range msg.ToServices {

		select {
		case <-ctx.Done():
			publishErrors = append(publishErrors, fmt.Errorf("publish cancelled by context for service %s: %w", toService, ctx.Err()))
			break // Daha fazla servise göndermeyi deneme
		default:
			// İşleme devam et
		}
		// Kendine mesaj göndermeyi engelle
		if toService == r.service {
			continue
		}

		// Bu mesaj tipi bu servis için izin veriliyor mu kontrol et
		if !isAllowedMessageType(toService, msg.Type) {
			log.Printf("[WHITELIST] Blocking message '%s' for service '%s' - Not in allowed list", msg.Type, toService)
			publishErrors = append(publishErrors, &MessagingError{
				Code:    "INVALID_TYPE",
				Message: fmt.Sprintf("Message type '%s' is not allowed for service '%s'", msg.Type, toService),
			})

			continue
		}

		// Tek servis için mesaj kopyasını hazırla
		singleMsg := msg
		singleMsg.ToServices = []ServiceType{toService} // Sadece bu servisi hedefle

		body, err := json.Marshal(singleMsg)
		if err != nil {
			publishErrors = append(publishErrors, &MessagingError{Code: "MARSHAL_FAILED", Message: "Failed to marshal message", Err: err})

			continue
		}

		serviceExchangeName := fmt.Sprintf("dc_clone.%s.service", toService)
		fmt.Printf("Mesaj gönderiliyor: %s -> %s\n", r.service, toService)

		err = r.channel.PublishWithContext(ctx,
			serviceExchangeName,
			"",
			true,  // Mandatory
			false, // Immediate
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				MessageId:    msg.ID,
				Timestamp:    msg.Created,
				Priority:     uint8(msg.Priority),
				Headers:      amqp.Table(msg.Headers),
				DeliveryMode: 2, // Persistent
			},
		)

		if err != nil {
			publishErrors = append(publishErrors, err)

			// Kritik mesajlar için özel işlem: hata durumunda kalıcı depolama
			if msg.Critical {
				singleMsg.ToServices = []ServiceType{toService} // Sadece başarısız servisi hedefle
				r.saveCriticalMessageToStorage(&singleMsg)

				// Kritik mesajlar için retry mekanizmasını da tetikle
				r.handleCriticalMessageRetry(&singleMsg)
			}
		} else {
			successServices = append(successServices, toService)
		}
	}

	// Hata kontrolü
	if len(publishErrors) > 0 {
		// Bazı servislere başarılı gönderim varsa bunu logla
		if len(successServices) > 0 {
			log.Printf("Message sent successfully to services: %v", successServices)
		}

		// Tüm servislere gönderim başarısız oldu mu?
		if len(successServices) == 0 {
			// Hiçbir servise gönderilemediyse, birleşik hata mesajı döndür
			errorMsg := fmt.Sprintf("Failed to publish message to any service: %v", publishErrors)
			return &MessagingError{Code: "PUBLISH_FAILED", Message: errorMsg}
		}

		// Kısmi başarı durumunda: bazı servislere gönderildi, bazılarına gönderilmedi
		log.Printf("Message partially sent. Errors: %v", publishErrors)
		return nil
	}

	return nil
}
func (r *RabbitClient) saveCriticalMessageToStorage(msg *Message) {
	fmt.Println("critical msg:", msg)
}

func (r *RabbitClient) handleCriticalMessageRetry(msg *Message) {
	// Retry sayısını artır
	msg.RetryCount++

	// Üstel artışla bekleme süresi (backoff strategy)
	retryDelay := int(math.Min(float64(1000*math.Pow(2, float64(msg.RetryCount))), 30000))

	body, err := json.Marshal(msg)
	if err != nil {
		log.Printf("handleCriticalMessageRetry marshal error: %v", err)
		r.saveCriticalMessageToStorage(msg)
		return
	}

	for _, toService := range msg.ToServices {
		routingKey := string(toService)

		log.Printf("Kritik retry mesajı gönderiliyor. ID: %s, ToService: %s",
			msg.ID, toService)

		err = r.channel.Publish(
			r.config.RetryExchangeName,
			routingKey,
			false,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				MessageId:    msg.ID,
				Timestamp:    time.Now(),
				DeliveryMode: 2, // Persistent
				Headers:      amqp.Table(msg.Headers),
				Expiration:   fmt.Sprintf("%d", retryDelay),
			},
		)

		if err != nil {
			log.Printf("handleCriticalMessageRetry publish error for service %s: %v", toService, err)

			// Tek servisi hedefleyen bir mesaj oluştur ve depola
			singleMsg := *msg
			singleMsg.ToServices = []ServiceType{toService}
			r.saveCriticalMessageToStorage(&singleMsg)
		} else {
			log.Printf("Kritik mesaj %s servisi için retry kuyruğuna %d saniye gecikmeyle gönderildi",
				toService, retryDelay/1000)
		}
	}
}
