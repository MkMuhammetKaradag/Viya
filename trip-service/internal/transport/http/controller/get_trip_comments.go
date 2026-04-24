package controller

import (
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetTripCommentsRequest struct {
	Page   int       `query:"page" validate:"omitempty,min=1"`
	Limit  int       `query:"limit" validate:"omitempty,min=1,max=50"`
	TripID uuid.UUID `uri:"trip_id" validate:"required"`
}

type GetTripCommentsResponse struct {
	Comments []domain.Comment `json:"comments"`
}
type GetTripCommentsController struct {
	usecase usecase.GetTripCommentsUseCase
}

func NewGetTripCommentsController(usecase usecase.GetTripCommentsUseCase) *GetTripCommentsController {
	return &GetTripCommentsController{
		usecase: usecase,
	}
}

func (c *GetTripCommentsController) Handle(fbrCtx fiber.Ctx, req *GetTripCommentsRequest) (*GetTripCommentsResponse, error) {
	rawUserID := fbrCtx.Get("X-User-ID")

	var currentUserID uuid.UUID
	if rawUserID != "" {
		parsedID, err := uuid.Parse(rawUserID)
		if err == nil {
			currentUserID = parsedID
		}
	}

	comments, err := c.usecase.Execute(fbrCtx.Context(), currentUserID, req.TripID, req.Limit, req.Page)
	if err != nil {
		return nil, err
	}

	return &GetTripCommentsResponse{Comments: comments}, nil
}
