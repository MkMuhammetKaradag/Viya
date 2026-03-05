package controller

import (
	"fmt"
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetTripRequest struct {
	TripID uuid.UUID `uri:"trip_id" validate:"required"`
}

type GetTripResponse struct {
	Trip *domain.Trip `json:"trip"`
}
type GetTripController struct {
	usecase usecase.GetTripUseCase
}

func NewGetTripController(usecase usecase.GetTripUseCase) *GetTripController {
	return &GetTripController{
		usecase: usecase,
	}
}

func (c *GetTripController) Handle(fbrCtx fiber.Ctx, req *GetTripRequest) (*GetTripResponse, error) {
	userID := fbrCtx.Get("X-User-ID")
	fmt.Println("userID:", userID)
	// if userID == "" {
	// 	return nil, fiber.ErrUnauthorized
	// }
	trip, err := c.usecase.Execute(fbrCtx.Context(), req.TripID)
	if err != nil {
		return nil, err
	}
	return &GetTripResponse{Trip: trip}, nil
}
