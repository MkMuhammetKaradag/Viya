package controller

import (
	"context"
	"fmt"
	"user-service/internal/transport/rabbitmq/usecase"
	"viya/pkg/messaging"

	"github.com/google/uuid"
)

type BlockUserHandler struct {
	usecase usecase.BlockUserUseCase
}

func NewBlockUserHandler(BlockUserUsecase usecase.BlockUserUseCase) *BlockUserHandler {
	return &BlockUserHandler{
		usecase: BlockUserUsecase,
	}
}

func (h *BlockUserHandler) Handle(msg messaging.Message) error {

	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid message data format")
	}

	// Tek bir hata kontrolü ile ilerlemek için:
	val := func(key string) string {
		s, _ := data[key].(string)
		return s
	}

	blocker := val("blocker")
	blocked := val("blocked")

	if blocked == "" || blocker == "" {
		return fmt.Errorf("missing required fields in message")
	}

	blockerUUID, err := uuid.Parse(blocker)
	if err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}
	blockedUUID, err := uuid.Parse(blocked)
	if err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}
	return h.usecase.Execute(context.Background(), blockerUUID, blockedUUID)
}
