package controller

import (
	"fmt"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ForkTripRequest struct {
	TripID uuid.UUID `uri:"trip_id" validate:"required"`
}

type ForkTripResponse struct {
	TripId  uuid.UUID `json:"trip"`
	Message string    `json:"message,omitempty"`
}
type ForkTripController struct {
	usecase usecase.ForkTripUseCase
}

func NewForkTripController(usecase usecase.ForkTripUseCase) *ForkTripController {
	return &ForkTripController{
		usecase: usecase,
	}
}

func (c *ForkTripController) Handle(fbrCtx fiber.Ctx, req *ForkTripRequest) (*ForkTripResponse, error) {
	rawUserID := fbrCtx.Get("X-User-ID")

	var currentUserID uuid.UUID
	if rawUserID != "" {
		parsedID, err := uuid.Parse(rawUserID)
		if err == nil {
			currentUserID = parsedID
		}
	}
	trip, err := c.usecase.Execute(fbrCtx.Context(), req.TripID, currentUserID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &ForkTripResponse{TripId: trip, Message: "forked ✅"}, nil
}
