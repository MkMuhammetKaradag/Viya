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
		return fmt.Errorf("err")
	}

	email, ok := data["email"].(string)
	if !ok {
		return fmt.Errorf("err")
	}

	id, ok := data["id"].(string)
	if !ok {
		return fmt.Errorf("err")
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("err")
	}

	userName, ok := data["username"].(string)
	if !ok {
		return fmt.Errorf("err")
	}
	ctx := context.Background()
	return h.usecase.Execute(ctx, idUUID, userName, email)
}
