package controller

import (
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type AddWayPointRequest struct {
	TripID uuid.UUID `json:"trip_id"`
	Lat    float64   `json:"lat" validate:"required"`
	Lon    float64   `json:"lon" validate:"required"`
	Note   string    `json:"note,omitempty"`
}

type AddWayPointResponse struct {
	Message    string    `json:"message"`
	WayPointID uuid.UUID `json:"waypoint_id"`
}

type AddWayPointController struct {
	usecase usecase.AddWayPointUseCase
}

func NewAddWaypointController(usecase usecase.AddWayPointUseCase) *AddWayPointController {
	return &AddWayPointController{
		usecase: usecase,
	}
}

func (c *AddWayPointController) Handle(fiberCtx fiber.Ctx, req *AddWayPointRequest) (*AddWayPointResponse, error) {

	wayPointModel := &domain.Waypoint{
		TripID: req.TripID,
		Lat:    req.Lat,
		Lon:    req.Lon,
		Note:   req.Note,
	}

	wpID, err := c.usecase.Execute(fiberCtx.Context(), wayPointModel)
	if err != nil {
		return nil, err
	}
	return &AddWayPointResponse{Message: "Waypoint added successfully with ID: " + wpID.String(), WayPointID: wpID}, nil
}
