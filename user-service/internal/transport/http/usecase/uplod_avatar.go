package usecase

import (
	"context"
	"log"
	"mime/multipart"
	"user-service/internal/domain"
	"viya/pkg/messaging"

	"github.com/google/uuid"
)

type UploadAvatarUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, file multipart.File) (string, error)
}
type uploadAvatarUseCase struct {
	cloudinaryService domain.CloudinaryService
	userRepository    domain.UserRepository
	rabbitClient      domain.RabbitMQClient
}

func NewUploadAvatarUseCase(userRepository domain.UserRepository, cldSvc domain.CloudinaryService, rabbitClient domain.RabbitMQClient) UploadAvatarUseCase {
	return &uploadAvatarUseCase{
		cloudinaryService: cldSvc,
		userRepository:    userRepository,
		rabbitClient:      rabbitClient,
	}
}

func (uc *uploadAvatarUseCase) Execute(ctx context.Context, userID uuid.UUID, file multipart.File) (string, error) {
	uploadRes, err := uc.cloudinaryService.UploadAvatar(ctx, file, userID.String())

	if err != nil {
		return "", err
	}
	err = uc.userRepository.UpdateAvatar(ctx, userID, uploadRes)
	if err != nil {
		return "", err
	}

	updatedMessage := messaging.Message{
		Type: messaging.UserTypes.UpdatedUser,
		ToServices: []messaging.ServiceType{
			messaging.SocialService,
			// messaging.TripService,
		},
		Data: map[string]interface{}{
			"id":         userID,
			"avatar_url": uploadRes,
		},
		Critical: true,
	}

	// 3. RabbitMQ üzerinden yayınla
	err = uc.rabbitClient.PublishMessage(ctx, updatedMessage)
	if err != nil {
		log.Printf("User update message could not be sent: %v", err)
	}

	return uploadRes, nil
}
