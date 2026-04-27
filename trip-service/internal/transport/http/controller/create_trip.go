package controller

import (
	"time"
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// Hibrit yapıyı destekleyen yeni Request DTO'su
type CreateTripRequest struct {
	Title         string     `json:"title" validate:"required,min=3"`
	Desc          string     `json:"desc,omitempty"`
	CoverImageURL *string    `json:"cover_image_url,omitempty"`
	IsPublic      bool       `json:"is_public"`
	PublishedAt   *time.Time `json:"published_at,omitempty"`
	IsActive      bool       `json:"is_active"`
	// İşte can alıcı nokta: Waypoint listesi opsiyonel (hibrit)
	Waypoints   []WaypointRequest `json:"waypoints,omitempty"`
	CategoryIDs []uuid.UUID       `json:"category_ids,omitempty"`
}

type WaypointRequest struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	OrderIndex  int            `json:"order_index"`
	Latitude    float64        `json:"latitude" validate:"required"`
	Longitude   float64        `json:"longitude" validate:"required"`
	Note        string         `json:"note"`
	Photos      []PhotoRequest `json:"photos,omitempty"`
	CategoryID  *uuid.UUID     `json:"category_id"`
}
type PhotoRequest struct {
	URL  string       `json:"url"`
	Tags []TagRequest `json:"tags,omitempty"`
}

type TagRequest struct {
	Label string  `json:"label"`
	XPos  float64 `json:"x_pos"`
	YPos  float64 `json:"y_pos"`
}
type CreateTripResponse struct {
	Message string    `json:"message"`
	TripID  uuid.UUID `json:"trip_id"`
}

type CreateTripController struct {
	usecase usecase.CreateTripUseCase
}

func NewCreateTripController(usecase usecase.CreateTripUseCase) *CreateTripController {
	return &CreateTripController{
		usecase: usecase,
	}
}

func (c *CreateTripController) Handle(fbrctx fiber.Ctx, req *CreateTripRequest) (*CreateTripResponse, error) {

	userIDStr := fbrctx.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid or missing user id")
	}

	// 1. Request verisini Domain modeline dönüştürelim
	tripModel := &domain.Trip{
		UserID:        userID,
		Title:         req.Title,
		Description:   req.Desc,
		CoverImageURL: req.CoverImageURL,
		IsPublic:      req.IsPublic,
		IsActive:      req.IsActive,

		CategoryIDs: req.CategoryIDs,
	}

	if req.PublishedAt != nil {
		tripModel.PublishedAt = *req.PublishedAt
	} else {
		tripModel.PublishedAt = time.Now()
	}

	// 2. Waypoints Döngüsü
	for _, wr := range req.Waypoints {
		waypoint := domain.Waypoint{
			Title:       wr.Title,
			Description: wr.Description,

			CategoryID: wr.CategoryID,
			OrderIndex: wr.OrderIndex,
			Latitude:   wr.Latitude,
			Longitude:  wr.Longitude,
			Note:       wr.Note,
		}

		// Fotoğraflar ve Etiketler
		for _, pr := range wr.Photos {
			photo := domain.Photo{URL: pr.URL}
			for _, tr := range pr.Tags {
				photo.Tags = append(photo.Tags, domain.Tag{
					Label: tr.Label,
					XPos:  tr.XPos,
					YPos:  tr.YPos,
				})
			}
			waypoint.Photos = append(waypoint.Photos, photo)
		}
		tripModel.Waypoints = append(tripModel.Waypoints, waypoint)
	}

	// 3. UseCase'i çalıştır
	id, err := c.usecase.Execute(fbrctx.Context(), tripModel)
	if err != nil {
		return nil, err
	}

	return &CreateTripResponse{
		Message: "Trip successfully created",
		TripID:  id,
	}, nil
}
