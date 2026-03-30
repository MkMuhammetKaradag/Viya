package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"user-service/internal/domain"

	"github.com/google/uuid"
)

type UploadBannerUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, file multipart.File) (string, error)
}
type uploadBannerUseCase struct {
	cloudinaryService domain.CloudinaryService
	userRepository    domain.UserRepository
}

func NewUploadBannerUseCase(userRepository domain.UserRepository, cldSvc domain.CloudinaryService) UploadBannerUseCase {
	return &uploadBannerUseCase{
		cloudinaryService: cldSvc,
		userRepository:    userRepository,
	}
}

func (uc *uploadBannerUseCase) Execute(ctx context.Context, userID uuid.UUID, file multipart.File) (string, error) {
	uploadRes, err := uc.cloudinaryService.UploadBanner(ctx, file, userID.String())

	if err != nil {
		return "", err
	}
	err = uc.userRepository.UpdateBanner(ctx, userID, uploadRes)
	if err != nil {
		return "", err
	}
	fmt.Println(uploadRes)
	return uploadRes, nil
}
