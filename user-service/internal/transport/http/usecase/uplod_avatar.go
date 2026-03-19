package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"user-service/internal/domain"

	"github.com/google/uuid"
)

type UploadAvatarUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, file multipart.File) error
}
type uploadAvatarUseCase struct {
	cloudinaryService domain.CloudinaryService
	userRepository    domain.UserRepository
}

func NewUploadAvatarUseCase(userRepository domain.UserRepository, cldSvc domain.CloudinaryService) UploadAvatarUseCase {
	return &uploadAvatarUseCase{
		cloudinaryService: cldSvc,
		userRepository:    userRepository,
	}
}

func (uc *uploadAvatarUseCase) Execute(ctx context.Context, userID uuid.UUID, file multipart.File) error {
	uploadRes, err := uc.cloudinaryService.UploadAvatar(ctx, file, userID.String())

	if err != nil {
		return err
	}
	err = uc.userRepository.UpdateAvatar(ctx, userID, uploadRes)
	if err != nil {
		return err
	}
	fmt.Println(uploadRes)
	return nil
}
