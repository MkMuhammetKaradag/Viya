package controller

import (
	"social-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type FollowRequestRequest struct {
	FollowerID uuid.UUID `uri:"follower_id" validate:"required,uuid4"`
	Action     string    `json:"action" validate:"required,oneof=ACCEPT REJECT"`
}

type FollowRequestResponse struct {
	Status string `json:"status"` // "PENDING" veya "ACCEPTED"
}

type FollowRequestController struct {
	usecase usecase.FollowRequestUseCase
}

func NewFollowRequestController(usecase usecase.FollowRequestUseCase) *FollowRequestController {
	return &FollowRequestController{
		usecase: usecase,
	}
}

func (c *FollowRequestController) Handle(fbrCtx fiber.Ctx, req *FollowRequestRequest) (*FollowRequestResponse, error) {
	rawUserID := fbrCtx.Get("X-User-ID")
	myID, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	err = c.usecase.Execute(fbrCtx.Context(), myID, req.FollowerID, req.Action)
	if err != nil {
		return nil, err
	}
	return &FollowRequestResponse{Status: "SUCCESS"}, nil
}
