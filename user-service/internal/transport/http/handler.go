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
	UploadBanner  *controller.UploadBannerController
	UpdateProfile *controller.UpdateProfileController
	GetMe         *controller.GetMeController
	SearchUsers   *controller.SearchUsersController
}

func NewHandlers(userRepo domain.UserRepository, cloudinaryService domain.CloudinaryService, rabbitClient domain.RabbitMQClient) *Handler {
	return &Handler{
		User: &UserHandlers{
			UploadAvatar:  controller.NewUploadAvatarController(usecase.NewUploadAvatarUseCase(userRepo, cloudinaryService, rabbitClient)),
			UploadBanner:  controller.NewUploadBannerController(usecase.NewUploadBannerUseCase(userRepo, cloudinaryService)),
			UpdateProfile: controller.NewUpdateProfileController(usecase.NewUpdateProfileUseCase(userRepo, rabbitClient)),
			GetMe:         controller.NewGetMeController(usecase.NewGetMeUseCase(userRepo)),
			SearchUsers:   controller.NewSearchUsersController(usecase.NewSearchUsersUseCase(userRepo)),
		},
	}
}
