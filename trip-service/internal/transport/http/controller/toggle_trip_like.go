package controller

import (
	"fmt"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ToggleTripLikeRequest struct {
	TripID uuid.UUID `uri:"trip_id" validate:"required"`
}

type ToggleTripLikeResponse struct {
	Message string `json:"messages"`
	IsLiked bool   `json:"is_liked"`
}
type ToggleTripLikeController struct {
	usecase usecase.ToggleTripLikeUseCase
}

func NewToggleTripLikeController(usecase usecase.ToggleTripLikeUseCase) *ToggleTripLikeController {
	return &ToggleTripLikeController{
		usecase: usecase,
	}
}

func (c *ToggleTripLikeController) Handle(fbrCtx fiber.Ctx, req *ToggleTripLikeRequest) (*ToggleTripLikeResponse, error) {
	rawUserID := fbrCtx.Get("X-User-ID")

	var currentUserID uuid.UUID
	if rawUserID != "" {
		parsedID, err := uuid.Parse(rawUserID)
		if err == nil {
			currentUserID = parsedID
		}
	}
	fmt.Println("geldi", currentUserID, req.TripID)
	isLiked, err := c.usecase.Execute(fbrCtx.Context(), req.TripID, currentUserID)
	if err != nil {
		return nil, err
	}
	message := "Trip liked successfully"
	if !isLiked {
		message = "Trip unliked successfully"
	}
	return &ToggleTripLikeResponse{Message: message, IsLiked: isLiked}, nil
}
