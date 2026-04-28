package controller

import (
	"user-service/internal/domain"
	"user-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetUserRequest struct {
	UserID uuid.UUID `uri:"user_id"`
}

type GetUserController struct {
	usecase usecase.GetUserUseCase
}

type GetUserResponse struct {
	User *domain.UserSummary `json:"user"`
}

func NewGetUserController(usecase usecase.GetUserUseCase) *GetUserController {
	return &GetUserController{
		usecase: usecase,
	}
}

func (h *GetUserController) Handle(fbrctx fiber.Ctx, req *GetUserRequest) (*GetUserResponse, error) {

	userIDStr := fbrctx.Get("X-User-ID")
	currentUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid or missing user id")
	}
	user, err := h.usecase.Execute(fbrctx.Context(), currentUserID, req.UserID)
	if err != nil {
		return nil, err
	}
	return &GetUserResponse{User: user}, nil
}
