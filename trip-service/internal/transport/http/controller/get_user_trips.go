package controller

import (
	"fmt"
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetUserTripsRequest struct {
	Page         int       `query:"page" validate:"omitempty,min=1"`
	Limit        int       `query:"limit" validate:"omitempty,min=1,max=50"`
	TargetUserID uuid.UUID `uri:"user_id"`
}

type GetUserTripsResponse struct {
	Trips []domain.TripSummary `json:"trip"`
}
type GetUserTripsController struct {
	usecase usecase.GetUserTripsUseCase
}

func NewGetUserTripsController(usecase usecase.GetUserTripsUseCase) *GetUserTripsController {
	return &GetUserTripsController{
		usecase: usecase,
	}
}

func (c *GetUserTripsController) Handle(fbrCtx fiber.Ctx, req *GetUserTripsRequest) (*GetUserTripsResponse, error) {

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 10
	}

	rawUserID := fbrCtx.Get("X-User-ID")
	fmt.Println("rawUserID:", rawUserID)

	parsedID, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	trips, err := c.usecase.Execute(fbrCtx.Context(), parsedID, req.TargetUserID, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}
	return &GetUserTripsResponse{Trips: trips}, nil
}
