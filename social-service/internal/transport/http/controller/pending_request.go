package controller

import (
	"social-service/internal/domain"
	"social-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type PendingRequestsRequest struct {
}

type PendingRequestsResponse struct {
	Requests []domain.PendingRequest `json:"requests"`
}

type PendingRequestsController struct {
	usecase usecase.PendingRequestsUseCase
}

func NewPendingRequestsController(usecase usecase.PendingRequestsUseCase) *PendingRequestsController {
	return &PendingRequestsController{
		usecase: usecase,
	}
}

func (c *PendingRequestsController) Handle(fbrCtx fiber.Ctx, req *PendingRequestsRequest) (*PendingRequestsResponse, error) {
	rawUserID := fbrCtx.Get("X-User-ID")
	myID, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	requests, err := c.usecase.Execute(fbrCtx.Context(), myID)
	if err != nil {
		return nil, err
	}
	return &PendingRequestsResponse{Requests: requests}, nil
}
