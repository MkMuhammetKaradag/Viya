package controller

import (
	"fmt"
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetMeTripsRequest struct {
	Page  int `query:"page" validate:"omitempty,min=1"`
	Limit int `query:"limit" validate:"omitempty,min=1,max=50"`
}

type GetMeTripsResponse struct {
	Trips []domain.TripSummary `json:"trip"`
}
type GetMeTripsController struct {
	usecase usecase.GetMeTripsUseCase
}

func NewGetMeTripsController(usecase usecase.GetMeTripsUseCase) *GetMeTripsController {
	return &GetMeTripsController{
		usecase: usecase,
	}
}

func (c *GetMeTripsController) Handle(fbrCtx fiber.Ctx, req *GetMeTripsRequest) (*GetMeTripsResponse, error) {

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

	trips, err := c.usecase.Execute(fbrCtx.Context(), parsedID, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}
	return &GetMeTripsResponse{Trips: trips}, nil
}
