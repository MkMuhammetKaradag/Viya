package controller

import (
	"social-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UnFollowUserRequest struct {
	TargetUserID uuid.UUID `uri:"target_user_id" validate:"required,uuid4"`
}

type UnFollowUserResponse struct {
	Message string `json:"message"`
}

type UnFollowUserController struct {
	usecase usecase.UnFollowUserUseCase
}

func NewUnFollowUserController(usecase usecase.UnFollowUserUseCase) *UnFollowUserController {
	return &UnFollowUserController{
		usecase: usecase,
	}
}

func (c *UnFollowUserController) Handle(fbrCtx fiber.Ctx, req *UnFollowUserRequest) (*UnFollowUserResponse, error) {
	rawUserID := fbrCtx.Get("X-User-ID")
	parsedUserID, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	err = c.usecase.Execute(fbrCtx.Context(), parsedUserID, req.TargetUserID)
	if err != nil {
		return nil, err
	}
	return &UnFollowUserResponse{Message: "Successfully unfollowed the user."}, nil
}
