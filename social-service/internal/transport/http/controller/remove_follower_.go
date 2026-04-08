package controller

import (
	"social-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type RemoveFollowerRequest struct {
	FollowerID uuid.UUID `uri:"follower_id" validate:"required,uuid4"`
}

type RemoveFollowerResponse struct {
	Message string `json:"message"`
}

type RemoveFollowerController struct {
	usecase usecase.RemoveFollowerUseCase
}

func NewRemoveFollowerController(usecase usecase.RemoveFollowerUseCase) *RemoveFollowerController {
	return &RemoveFollowerController{
		usecase: usecase,
	}
}

func (c *RemoveFollowerController) Handle(fbrCtx fiber.Ctx, req *RemoveFollowerRequest) (*RemoveFollowerResponse, error) {
	rawUserID := fbrCtx.Get("X-User-ID")
	parsedUserID, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	err = c.usecase.Execute(fbrCtx.Context(), parsedUserID, req.FollowerID)
	if err != nil {
		return nil, err
	}
	return &RemoveFollowerResponse{Message: "Successfully remove follower the user."}, nil
}
