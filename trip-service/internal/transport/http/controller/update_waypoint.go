package controller

import (
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UpdateWaypointRequest struct {
	WayPointID  uuid.UUID `uri:"waypoint_id" validate:"required"`
	Title       *string   `json:"title" `
	Description *string   `json:"description" `
	Latitude    *float64  `json:"latitude"  `
	Longitude   *float64  `json:"longitude" `
}

type UpdateWaypointResponse struct {
	Message string `json:"message"`
}

type UpdateWaypointController struct {
	usecase usecase.UpdateWaypointUseCase
}

func NewUpdateWaypointController(usecase usecase.UpdateWaypointUseCase) *UpdateWaypointController {
	return &UpdateWaypointController{
		usecase: usecase,
	}
}

func (c *UpdateWaypointController) Handle(fiberCtx fiber.Ctx, req *UpdateWaypointRequest) (*UpdateWaypointResponse, error) {

	err := c.usecase.Execute(fiberCtx.Context(), &domain.UpdateWaypointInput{
		WayPointID:  req.WayPointID,
		Title:       req.Title,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateWaypointResponse{Message: "Waypoint updated successfully"}, nil
}
