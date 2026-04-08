package controller

import (
	"social-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type FollowUserRequest struct {
	TargetUserID uuid.UUID `uri:"target_user_id" validate:"required,uuid4"`
}

type FollowUserResponse struct {
	Status string `json:"status"` // "PENDING" veya "ACCEPTED"
}

type FollowUserController struct {
	usecase usecase.FollowUserUseCase
}

func NewFollowUserController(usecase usecase.FollowUserUseCase) *FollowUserController {
	return &FollowUserController{
		usecase: usecase,
	}
}

func (c *FollowUserController) Handle(fbrCtx fiber.Ctx, req *FollowUserRequest) (*FollowUserResponse, error) {
	rawUserID := fbrCtx.Get("X-User-ID")
	parsedUserID, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	status, err := c.usecase.Execute(fbrCtx.Context(), parsedUserID, req.TargetUserID)
	if err != nil {
		return nil, err
	}
	return &FollowUserResponse{Status: status}, nil
}
