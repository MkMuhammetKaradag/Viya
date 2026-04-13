package controller

import (
	"context"
	"fmt"
	"trip-service/internal/transport/rabbitmq/usecase"
	"viya/pkg/messaging"

	"github.com/google/uuid"
)

type UpdatedUserHandler struct {
	usecase usecase.UpdatedUserUseCase
}

func NewUserUpdatedHandler(UpdatedUserUsecase usecase.UpdatedUserUseCase) *UpdatedUserHandler {
	return &UpdatedUserHandler{
		usecase: UpdatedUserUsecase,
	}
}

func (h *UpdatedUserHandler) Handle(msg messaging.Message) error {
	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid message data format")
	}

	idStr, _ := data["id"].(string)
	idUUID, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}

	// Değişkenleri pointer olarak hazırlıyoruz ki "yoksa nil kalsın" diyebilelim
	var isPrivate *bool
	var avatarURL *string

	// Eğer map'te "is_private" anahtarı varsa değerini al
	if val, exists := data["is_private"]; exists {
		p := val.(bool) // Dikkat: JSON'dan bool olarak gelir
		isPrivate = &p
	}

	// Eğer map'te "avatar_url" anahtarı varsa değerini al
	if val, exists := data["avatar_url"]; exists {
		s := val.(string)
		avatarURL = &s
	}

	// UseCase artık bu pointer'ları almalı
	return h.usecase.Execute(context.Background(), idUUID, isPrivate, avatarURL)
}
