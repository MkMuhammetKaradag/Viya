package controller

import (
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DeleteWaypointRequest struct {
	WaypointID uuid.UUID `uri:"waypoint_id" validate:"required"`
}
type DeleteWaypointResponse struct {
	Message string `json:"message"`
}

type DeleteWaypointController struct {
	usecase usecase.DeleteWaypointUseCase
}

func NewDeleteWaypointController(usecase usecase.DeleteWaypointUseCase) *DeleteWaypointController {
	return &DeleteWaypointController{
		usecase: usecase,
	}
}

func (c *DeleteWaypointController) Handle(fiberCtx fiber.Ctx, req *DeleteWaypointRequest) (*DeleteWaypointResponse, error) {

	err := c.usecase.Execute(fiberCtx.Context(), req.WaypointID)
	if err != nil {
		return nil, err
	}
	return &DeleteWaypointResponse{Message: "Waypoint deleted successfully"}, nil
}
