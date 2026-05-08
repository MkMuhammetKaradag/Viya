package controller

import (
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetHomeFeedTripsRequest struct {
	Page  int `query:"page" validate:"omitempty,min=1"`
	Limit int `query:"limit" validate:"omitempty,min=1,max=50"`
}

type GetHomeFeedTripsResponse struct {
	Trips []domain.TripExploreDTO `json:"trips"`
}
type GetHomeFeedTripsController struct {
	usecase usecase.GetHomeFeedTripsUseCase
}

func NewGetHomeFeedTripsController(usecase usecase.GetHomeFeedTripsUseCase) *GetHomeFeedTripsController {
	return &GetHomeFeedTripsController{
		usecase: usecase,
	}
}

func (c *GetHomeFeedTripsController) Handle(fbrCtx fiber.Ctx, req *GetHomeFeedTripsRequest) (*GetHomeFeedTripsResponse, error) {
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

	return &GetHomeFeedTripsResponse{Trips: trips}, nil
}
