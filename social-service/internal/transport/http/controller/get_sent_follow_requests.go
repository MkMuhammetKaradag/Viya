package controller

import (
	"fmt"
	"social-service/internal/domain"
	"social-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetSentFollowRequestsRequest struct {
}

type GetSentFollowRequestsResponse struct {
	Requests []domain.PendingRequest `json:"requests"`
}

type GetSentFollowRequestsController struct {
	usecase usecase.GetSentFollowRequestsUseCase
}

func NewGetSentFollowRequestsController(usecase usecase.GetSentFollowRequestsUseCase) *GetSentFollowRequestsController {
	return &GetSentFollowRequestsController{
		usecase: usecase,
	}
}

func (c *GetSentFollowRequestsController) Handle(fbrCtx fiber.Ctx, req *GetSentFollowRequestsRequest) (*GetSentFollowRequestsResponse, error) {
	rawUserID := fbrCtx.Get("X-User-ID")
	myID, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}
	fmt.Printf("Fetching sent follow requests for user: %s\n", myID)

	requests, err := c.usecase.Execute(fbrCtx.Context(), myID)
	if err != nil {
		return nil, err
	}
	return &GetSentFollowRequestsResponse{Requests: requests}, nil
}
