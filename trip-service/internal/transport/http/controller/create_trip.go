package controller

import (
	"fmt"
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CreateTripRequest struct {
	UserID uuid.UUID `json:"user_id"`
	// ID       uuid.UUID `params:"id,omitempty" validate:"required" `
	// Search   string    `query:"search,omitempty"`
	Title    string `json:"title" validate:"required,min=3"`
	Desc     string `json:"desc,omitempty"`
	IsActive bool   `json:"is_active"`
}

type CreateTripResponse struct {
	Message string `json:"message"`
}

type CreateTripController struct {
	usecase usecase.CreateTripUseCase
}

func NewCreateTripController(usecase usecase.CreateTripUseCase) *CreateTripController {
	return &CreateTripController{
		usecase: usecase,
	}
}

func (c *CreateTripController) Handle(fiberCtx fiber.Ctx, req *CreateTripRequest) (*CreateTripResponse, error) {
	// Burada UseCase'i çağırıp trip oluşturma işlemini yapacağız
	// Şimdilik sadece dummy response dönüyoruz

	tripModel := &domain.Trip{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Desc,
		IsActive:    req.IsActive,
	}

	id, err := c.usecase.Execute(fiberCtx.Context(), tripModel)
	if err != nil {
		return nil, err
	}
	return &CreateTripResponse{Message: fmt.Sprintf("trip id %s", id)}, nil
}
