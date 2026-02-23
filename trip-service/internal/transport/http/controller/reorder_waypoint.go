package controller

import (
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ReorderRequest struct {
	WayPointID uuid.UUID `uri:"waypoint_id" validate:"required"`
	Index      int       `json:"index" validate:"required" min:"1"`
}

type ReorderResponse struct {
	Message string `json:"message"`
}

type ReorderController struct {
	usecase usecase.ReorderUseCase
}

func NewReorderController(usecase usecase.ReorderUseCase) *ReorderController {
	return &ReorderController{
		usecase: usecase,
	}
}

func (c *ReorderController) Handle(fiberCtx fiber.Ctx, req *ReorderRequest) (*ReorderResponse, error) {

	err := c.usecase.Execute(fiberCtx, req.WayPointID, req.Index)
	if err != nil {
		return nil, err
	}

	return &ReorderResponse{Message: "WayPoint reordered successfully"}, nil
}
