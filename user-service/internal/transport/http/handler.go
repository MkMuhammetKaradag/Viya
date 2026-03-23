package http

import (
	"user-service/internal/domain"
	"user-service/internal/transport/http/controller"
	"user-service/internal/transport/http/usecase"
)

type Handler struct {
	User *UserHandlers
}

type UserHandlers struct {
	UploadAvatar  *controller.UploadAvatarController
	UpdateProfile *controller.UpdateProfileController
}

func NewHandlers(userRepo domain.UserRepository, cloudinaryService domain.CloudinaryService) *Handler {
	return &Handler{
		User: &UserHandlers{
			UploadAvatar:  controller.NewUploadAvatarController(usecase.NewUploadAvatarUseCase(userRepo, cloudinaryService)),
			UpdateProfile: controller.NewUpdateProfileController(usecase.NewUpdateProfileUseCase(userRepo)),
		},
	}
}
