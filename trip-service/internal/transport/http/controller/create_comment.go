package controller

import (
	"time"
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CreateCommentRequest struct {
	TripID   uuid.UUID  `json:"trip_id" validate:"required"`
	Content  string     `json:"content" validate:"required,min=2"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
}

type CreateCommentResponse struct {
	Message string `json:"message"`
}

type CreateCommentController struct {
	usecase usecase.CreateCommentUseCase
}

func NewCreateCommentController(usecase usecase.CreateCommentUseCase) *CreateCommentController {
	return &CreateCommentController{
		usecase: usecase,
	}
}

func (c *CreateCommentController) Handle(fbrctx fiber.Ctx, req *CreateCommentRequest) (*CreateCommentResponse, error) {

	userIDStr := fbrctx.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid or missing user id")
	}

	comment := &domain.Comment{
		UserID:    userID,
		TripID:    req.TripID,
		Content:   req.Content,
		ParentID:  req.ParentID,
		CreatedAt: time.Now(),
	}

	_, err = c.usecase.Execute(fbrctx.Context(), comment)
	if err != nil {
		return nil, err
	}

	return &CreateCommentResponse{
		Message: "comment successfully created",
	}, nil
}
