package controller

import (
	"social-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type BlockUserRequest struct {
	TargetUserID uuid.UUID `uri:"target_user_id" validate:"required,uuid4"`
}

type BlockUserResponse struct {
	Message string `json:"message"`
}

type BlockUserController struct {
	usecase usecase.BlockUserUseCase
}

func NewBlockUserController(usecase usecase.BlockUserUseCase) *BlockUserController {
	return &BlockUserController{
		usecase: usecase,
	}
}

func (c *BlockUserController) Handle(fbrCtx fiber.Ctx, req *BlockUserRequest) (*BlockUserResponse, error) {
	rawUserID := fbrCtx.Get("X-User-ID")
	parsedUserID, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	err = c.usecase.Execute(fbrCtx.Context(), parsedUserID, req.TargetUserID)
	if err != nil {
		return nil, err
	}
	return &BlockUserResponse{Message: "User blocked successfully"}, nil
}
