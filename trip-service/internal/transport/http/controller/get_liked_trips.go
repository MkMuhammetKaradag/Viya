package controller

import (
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetLikedTripsRequest struct {
	Page  int `query:"page" validate:"omitempty,min=1"`
	Limit int `query:"limit" validate:"omitempty,min=1,max=50"`
}

type GetLikedTripsResponse struct {
	Trips []domain.TripSummary `json:"trips"`
}
type GetLikedTripsController struct {
	usecase usecase.GetLikedTripsUseCase
}

func NewGetLikedTripsController(usecase usecase.GetLikedTripsUseCase) *GetLikedTripsController {
	return &GetLikedTripsController{
		usecase: usecase,
	}
}

func (c *GetLikedTripsController) Handle(fbrCtx fiber.Ctx, req *GetLikedTripsRequest) (*GetLikedTripsResponse, error) {
	rawUserID := fbrCtx.Get("X-User-ID")

	var currentUserID uuid.UUID
	if rawUserID != "" {
		parsedID, err := uuid.Parse(rawUserID)
		if err == nil {
			currentUserID = parsedID
		}
	}

	trips, err := c.usecase.Execute(fbrCtx.Context(), currentUserID, req.Limit, req.Page)
	if err != nil {
		return nil, err
	}

	return &GetLikedTripsResponse{Trips: trips}, nil
}
