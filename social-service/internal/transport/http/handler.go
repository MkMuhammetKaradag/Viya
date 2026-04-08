package http

import (
	"social-service/internal/domain"
	"social-service/internal/transport/http/controller"
	"social-service/internal/transport/http/usecase"
)

type Handlers struct {
	Social *socialHandlers
}

type socialHandlers struct {
	FollowUser *controller.FollowUserController
	BlockUser  *controller.BlockUserController
}

func NewHandlers(repo domain.SocialRepository) *Handlers {
	socialHandlers := &socialHandlers{
		FollowUser: controller.NewFollowUserController(usecase.NewFollowUserUseCase(repo)),
		BlockUser:  controller.NewBlockUserController(usecase.NewBlockUserUseCase(repo)),
	}

	return &Handlers{
		Social: socialHandlers,
	}
}
