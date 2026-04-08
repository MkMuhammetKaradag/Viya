package controller

import (
	"social-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UnblockUserRequest struct {
	TargetUserID uuid.UUID `uri:"target_user_id" validate:"required,uuid4"`
}

type UnblockUserResponse struct {
	Message string `json:"message"`
}

type UnblockUserController struct {
	usecase usecase.UnblockUserUseCase
}

func NewUnblockUserController(usecase usecase.UnblockUserUseCase) *UnblockUserController {
	return &UnblockUserController{
		usecase: usecase,
	}
}

func (c *UnblockUserController) Handle(fbrCtx fiber.Ctx, req *UnblockUserRequest) (*UnblockUserResponse, error) {
	rawUserID := fbrCtx.Get("X-User-ID")
	parsedUserID, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	err = c.usecase.Execute(fbrCtx.Context(), parsedUserID, req.TargetUserID)
	if err != nil {
		return nil, err
	}
	return &UnblockUserResponse{Message: "User unblocked successfully"}, nil
}
