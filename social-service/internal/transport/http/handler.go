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
	FollowUser     *controller.FollowUserController
	UnFollowUser   *controller.UnFollowUserController
	RemoveFollower *controller.RemoveFollowerController
	FollowRequest  *controller.FollowRequestController

	PendingRequests *controller.PendingRequestsController
	BlockUser       *controller.BlockUserController
}

func NewHandlers(repo domain.SocialRepository) *Handlers {
	socialHandlers := &socialHandlers{
		FollowUser:      controller.NewFollowUserController(usecase.NewFollowUserUseCase(repo)),
		UnFollowUser:    controller.NewUnFollowUserController(usecase.NewUnFollowUserUseCase(repo)),
		RemoveFollower:  controller.NewRemoveFollowerController(usecase.NewRemoveFollowerUseCase(repo)),
		BlockUser:       controller.NewBlockUserController(usecase.NewBlockUserUseCase(repo)),
		FollowRequest:   controller.NewFollowRequestController(usecase.NewFollowRequestUseCase(repo)),
		PendingRequests: controller.NewPendingRequestsController(usecase.NewPendingRequestsUseCase(repo)),
	}

	return &Handlers{
		Social: socialHandlers,
	}
}
