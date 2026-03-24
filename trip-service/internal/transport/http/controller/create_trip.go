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
	// UserID        uuid.UUID  `json:"user_id" validate:"required"`
	Title         string     `json:"title" validate:"required,min=3"`
	Desc          string     `json:"desc,omitempty"`
	CoverImageURL *string    `json:"cover_image_url,omitempty"`
	IsPublic      bool       `json:"is_public"`
	PublishedAt   *time.Time `json:"published_at,omitempty"`
	IsActive      bool       `json:"is_active"`
	// İşte can alıcı nokta: Waypoint listesi opsiyonel (hibrit)
	Waypoints []WaypointRequest `json:"waypoints,omitempty"`
}

type WaypointRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	OrderIndex  int      `json:"order_index"`
	Latitude    float64  `json:"latitude" validate:"required"`
	Longitude   float64  `json:"longitude" validate:"required"`
	Note        string   `json:"note"`
	Photos      []string `json:"photos,omitempty"` // Cloudinary URL'leri
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
	}

	// PublishedAt boşsa şu anı set et
	if req.PublishedAt != nil {
		tripModel.PublishedAt = *req.PublishedAt
	} else {
		tripModel.PublishedAt = time.Now()
	}

	// 2. Eğer Waypoints varsa onları da domain modeline ekleyelim
	for _, wr := range req.Waypoints {
		waypoint := domain.Waypoint{
			Title:       wr.Title,
			Description: wr.Description,
			OrderIndex:  wr.OrderIndex,
			Latitude:    wr.Latitude,
			Longitude:   wr.Longitude,
			Note:        wr.Note,
		}
		for _, photoURL := range wr.Photos {
			waypoint.Photos = append(waypoint.Photos, domain.Photo{
				URL: photoURL,
			})
		}
		tripModel.Waypoints = append(tripModel.Waypoints, waypoint)
	}

	// 3. UseCase'i çalıştır (Repository'deki Transaction'ı bu tetikleyecek)
	id, err := c.usecase.Execute(fbrctx.Context(), tripModel)
	if err != nil {
		return nil, err
	}

	return &CreateTripResponse{
		Message: "Trip successfully created",
		TripID:  id,
	}, nil
}
