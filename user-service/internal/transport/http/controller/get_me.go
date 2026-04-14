package controller

import (
	"fmt"
	"user-service/internal/domain"
	"user-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetMeRequest struct {
}

type GetMeController struct {
	usecase usecase.GetMeUseCase
}

type GetMeResponse struct {
	User *domain.User `json:"user"`
}

func NewGetMeController(usecase usecase.GetMeUseCase) *GetMeController {
	return &GetMeController{
		usecase: usecase,
	}
}

func (h *GetMeController) Handle(fbrctx fiber.Ctx, req *GetMeRequest) (*GetMeResponse, error) {
	fmt.Println("geldi:")

	userIDStr := fbrctx.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid or missing user id")
	}
	user, err := h.usecase.Execute(fbrctx.Context(), userID)
	if err != nil {
		return nil, err
	}
	return &GetMeResponse{User: user}, nil
}
