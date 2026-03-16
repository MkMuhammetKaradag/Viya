package rabbitmq

import (
	"fmt"
	"user-service/internal/domain"
)

// ConfigureRabbitMQ, tüm kuyrukları ve handler'ları eşleştirir
func ConfigureRabbitMQ(rabbit domain.RabbitMQClient, userRepo domain.UserRepository) {

	// Kullanıcı oluşturma kuyruğunu dinle
	rabbit.Consume("user_created", func(data []byte) error {
		// Burada gelen mesajı işle
		fmt.Printf("Mesaj alındı ve işleniyor: %s\n", string(data))
		// İleride buraya: return HandleUserCreated(data, userRepo)
		return nil
	})

	// Başka kuyruklar varsa altına ekleyebilirsin
	// rabbit.Consume("order_created_queue", ...)
}
