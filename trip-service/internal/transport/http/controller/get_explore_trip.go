package controller

import (
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetExploreTripsRequest struct {
	Page  int `query:"page" validate:"omitempty,min=1"`
	Limit int `query:"limit" validate:"omitempty,min=1,max=50"`
}

type GetExploreTripsResponse struct {
	Trips []domain.TripExploreDTO `json:"trips"`
}
type GetExploreTripsController struct {
	usecase usecase.GetExploreTripsUseCase
}

func NewGetExploreTripsController(usecase usecase.GetExploreTripsUseCase) *GetExploreTripsController {
	return &GetExploreTripsController{
		usecase: usecase,
	}
}

func (c *GetExploreTripsController) Handle(fbrCtx fiber.Ctx, req *GetExploreTripsRequest) (*GetExploreTripsResponse, error) {
	rawUserID := fbrCtx.Get("X-User-ID")

	var currentUserID uuid.UUID
	if rawUserID != "" {
		parsedID, err := uuid.Parse(rawUserID)
		if err == nil {
			currentUserID = parsedID
		}
	}
	offset := (req.Page - 1) * req.Limit
	trips, err := c.usecase.Execute(fbrCtx.Context(), currentUserID, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	return &GetExploreTripsResponse{Trips: trips}, nil
}
