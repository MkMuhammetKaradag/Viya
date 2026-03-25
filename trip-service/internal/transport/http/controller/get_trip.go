package controller

import (
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
	rawUserID := fbrCtx.Get("X-User-ID")

	var currentUserID uuid.UUID
	if rawUserID != "" {
		// Eğer giriş yapılmışsa UUID'ye çevirelim
		parsedID, err := uuid.Parse(rawUserID)
		if err == nil {
			currentUserID = parsedID
		}
	}
	trip, err := c.usecase.Execute(fbrCtx.Context(), req.TripID, currentUserID)
	if err != nil {
		return nil, err
	}
	return &GetTripResponse{Trip: trip}, nil
}
