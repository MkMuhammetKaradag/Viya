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

	PendingRequests       *controller.PendingRequestsController
	GetSentFollowRequests *controller.GetSentFollowRequestsController
	BlockUser             *controller.BlockUserController
	UnblockUser           *controller.UnblockUserController
}

func NewHandlers(repo domain.SocialRepository, rabbitClient domain.RabbitMQClient) *Handlers {
	socialHandlers := &socialHandlers{
		FollowUser:            controller.NewFollowUserController(usecase.NewFollowUserUseCase(repo, rabbitClient)),
		UnFollowUser:          controller.NewUnFollowUserController(usecase.NewUnFollowUserUseCase(repo, rabbitClient)),
		RemoveFollower:        controller.NewRemoveFollowerController(usecase.NewRemoveFollowerUseCase(repo, rabbitClient)),
		BlockUser:             controller.NewBlockUserController(usecase.NewBlockUserUseCase(repo, rabbitClient)),
		UnblockUser:           controller.NewUnblockUserController(usecase.NewUnblockUserUseCase(repo, rabbitClient)),
		FollowRequest:         controller.NewFollowRequestController(usecase.NewFollowRequestUseCase(repo, rabbitClient)),
		PendingRequests:       controller.NewPendingRequestsController(usecase.NewPendingRequestsUseCase(repo)),
		GetSentFollowRequests: controller.NewGetSentFollowRequestsController(usecase.NewGetSentFollowRequestsUseCase(repo)),
	}

	return &Handlers{
		Social: socialHandlers,
	}
}
