package controller

import (
	"fmt"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type AddWayPointPhotosRequest struct {
	WayPointID uuid.UUID `uri:"waypoint_id" validate:"required"`
}

type AddWayPointPhotosResponse struct {
	Message string `json:"message"`
}

type AddWayPointPhotosController struct {
	usecase usecase.AddWayPointPhotosUseCase
}

func NewAddWayPointPhotosController(usecase usecase.AddWayPointPhotosUseCase) *AddWayPointPhotosController {
	return &AddWayPointPhotosController{
		usecase: usecase,
	}
}

func (c *AddWayPointPhotosController) Handle(fiberCtx fiber.Ctx, req *AddWayPointPhotosRequest) (*AddWayPointPhotosResponse, error) {

	form, err := fiberCtx.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf("form error: %w", err)
	}

	files := form.File["images"]
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one image is required")
	}
	fmt.Println(files)

	err = c.usecase.Execute(fiberCtx, req.WayPointID, files)
	if err != nil {
		return nil, err
	}

	return &AddWayPointPhotosResponse{Message: "Photos added successfully"}, nil
}
