package controller

import (
	"context"
	"fmt"
	"user-service/internal/transport/rabbitmq/usecase"
	"viya/pkg/messaging"

	"github.com/google/uuid"
)

type CreatedUserHandler struct {
	usecase usecase.CreateUserUseCase
}

func NewUserCreatedHandler(createdUserUsecase usecase.CreateUserUseCase) *CreatedUserHandler {
	return &CreatedUserHandler{
		usecase: createdUserUsecase,
	}
}

func (h *CreatedUserHandler) Handle(msg messaging.Message) error {

	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid message data format")
	}

	// Tek bir hata kontrolü ile ilerlemek için:
	val := func(key string) string {
		s, _ := data[key].(string)
		return s
	}

	idStr := val("id")
	email := val("email")
	username := val("username")

	if idStr == "" || email == "" {
		return fmt.Errorf("missing required fields in message")
	}

	idUUID, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}
	return h.usecase.Execute(context.Background(), idUUID, username, email)

}
