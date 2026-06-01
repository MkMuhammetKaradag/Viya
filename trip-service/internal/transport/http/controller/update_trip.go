package controller

import (
	"fmt"
	"time"
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UpdateTripRequest struct {
	TripID uuid.UUID `uri:"trip_id" validate:"required"`

	Title         *string     `json:"title,omitempty" validate:"omitempty,min=3"`
	Desc          *string     `json:"desc,omitempty"`
	CoverImageURL *string     `json:"cover_image_url,omitempty"`
	IsPublic      *bool       `json:"is_public,omitempty"`
	PublishedAt   *time.Time  `json:"published_at,omitempty"`
	IsActive      *bool       `json:"is_active,omitempty"`
	CategoryIDs   []uuid.UUID `json:"category_ids,omitempty"`
}

type UpdateTripResponse struct {
	Message string    `json:"message"`
	TripID  uuid.UUID `json:"trip_id"`
}

type UpdateTripController struct {
	usecase usecase.UpdateTripUseCase
}

func NewUpdateTripController(usecase usecase.UpdateTripUseCase) *UpdateTripController {
	return &UpdateTripController{
		usecase: usecase,
	}
}

func (c *UpdateTripController) Handle(fbrctx fiber.Ctx, req *UpdateTripRequest) (*UpdateTripResponse, error) {
	fmt.Println("UpdateTripController.Handle called with TripID:", req.TripID)
	userIDStr := fbrctx.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid or missing user id")
	}

	// 2. Sadece gelen alanları domain katmanına taşıyoruz
	tripModel := &domain.Trip{
		ID:     req.TripID,
		UserID: userID,
	}

	// Pointer kontrolleri: Sadece gelen verileri modele set ediyoruz
	if req.Title != nil {
		tripModel.Title = *req.Title
	}
	if req.Desc != nil {
		tripModel.Description = req.Desc
	}
	if req.CoverImageURL != nil {
		tripModel.CoverImageURL = req.CoverImageURL
	}
	if req.IsPublic != nil {
		tripModel.IsPublic = *req.IsPublic
	}
	if req.IsActive != nil {
		tripModel.IsActive = *req.IsActive
	}
	if req.CategoryIDs != nil {
		tripModel.CategoryIDs = req.CategoryIDs
	}
	if req.PublishedAt != nil {
		tripModel.PublishedAt = *req.PublishedAt
	}
	// NOT: `else { tripModel.PublishedAt = time.Now() }` kısmını sildik.
	// Eğer PublishedAt gelmediyse veritabanındaki eski değer korunmalı (bunu Usecase/Repository katmanında handle etmelisin usta).

	id, err := c.usecase.Execute(fbrctx.Context(), tripModel)
	if err != nil {
		return nil, err
	}

	return &UpdateTripResponse{
		Message: "Trip successfully updated",
		TripID:  id,
	}, nil
}
