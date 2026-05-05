package controller

import (
	"context"
	"fmt"
	"user-service/internal/transport/rabbitmq/usecase"
	"viya/pkg/messaging"

	"github.com/google/uuid"
)

type UnFollowUserHandler struct {
	usecase usecase.UnFollowUserUseCase
}

func NewUnFollowUserHandler(UnFollowUserUsecase usecase.UnFollowUserUseCase) *UnFollowUserHandler {
	return &UnFollowUserHandler{
		usecase: UnFollowUserUsecase,
	}
}

func (h *UnFollowUserHandler) Handle(msg messaging.Message) error {

	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid message data format")
	}

	// Tek bir hata kontrolü ile ilerlemek için:
	val := func(key string) string {
		s, _ := data[key].(string)
		return s
	}

	follower := val("follower")
	following := val("following")

	if follower == "" || following == "" {
		return fmt.Errorf("missing required fields in message")
	}

	followerUUID, err := uuid.Parse(follower)
	if err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}
	followingUUID, err := uuid.Parse(following)
	if err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}
	return h.usecase.Execute(context.Background(), followerUUID, followingUUID)
}
