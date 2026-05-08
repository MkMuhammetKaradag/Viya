package controller

import (
	"fmt"
	"social-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetPendingCountRequest struct {
}

type GetPendingCountResponse struct {
	Count int `json:"count"`
}

type GetPendingCountController struct {
	usecase usecase.GetPendingCountUseCase
}

func NewGetPendingCountController(usecase usecase.GetPendingCountUseCase) *GetPendingCountController {
	return &GetPendingCountController{
		usecase: usecase,
	}
}

func (c *GetPendingCountController) Handle(fbrCtx fiber.Ctx, req *GetPendingCountRequest) (*GetPendingCountResponse, error) {
	rawUserID := fbrCtx.Get("X-User-ID")
	myID, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}
	fmt.Printf("Fetching sent follow requests for user: %s\n", myID)

	count, err := c.usecase.Execute(fbrCtx.Context(), myID)
	if err != nil {
		return nil, err
	}
	return &GetPendingCountResponse{Count: count}, nil
}
