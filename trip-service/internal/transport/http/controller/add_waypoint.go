package controller

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type AddWayPointRequest struct {
	// TripID     uuid.UUID `json:"trip_id"`
	// Lat        float64   `json:"lat" validate:"required"`
	// Lon        float64   `json:"lon" validate:"required"`
	// Desc       string    `json:"desc,omitempty"`
	// Title      string    `json:"title,omitempty"`
	// OrderIndex int       `json:"order_index,omitempty"`
}

type AddWayPointResponse struct {
	Message    string    `json:"message"`
	WayPointID uuid.UUID `json:"waypoint_id"`
}

type AddWayPointController struct {
	usecase usecase.AddWayPointUseCase
}

func NewAddWaypointController(usecase usecase.AddWayPointUseCase) *AddWayPointController {
	return &AddWayPointController{
		usecase: usecase,
	}
}

func (c *AddWayPointController) Handle(fiberCtx fiber.Ctx, req *AddWayPointRequest) (*AddWayPointResponse, error) {
	// 1. Form verilerini (text alanlarını) okuyalım
	tripID, err := uuid.Parse(fiberCtx.FormValue("trip_id"))
	if err != nil {
		return nil, fmt.Errorf("invalid trip_id: %w", err)
	}

	lat, _ := strconv.ParseFloat(fiberCtx.FormValue("lat"), 64)
	lon, _ := strconv.ParseFloat(fiberCtx.FormValue("lon"), 64)
	orderIndex, _ := strconv.Atoi(fiberCtx.FormValue("order_index"))

	wayPointModel := &domain.Waypoint{
		TripID:      tripID,
		Latitude:    lat,
		Longitude:   lon,
		Description: fiberCtx.FormValue("desc"),
		Title:       fiberCtx.FormValue("title"),
		OrderIndex:  orderIndex,
		Note:        fiberCtx.FormValue("note"),
	}

	// 2. Fotoğrafları (dosyaları) alalım
	form, err := fiberCtx.MultipartForm()
	var files []*multipart.FileHeader
	if err == nil { // Eğer dosya gönderilmişse form'dan alalım
		files = form.File["images"]
	}

	// 3. UseCase'e hem modeli hem de dosyaları gönderiyoruz
	// UseCase önce waypoint'i kaydedecek, sonra fotoğrafları worker'a atacak.
	wpID, err := c.usecase.Execute(fiberCtx.Context(), wayPointModel, files)
	if err != nil {
		return nil, err
	}

	return &AddWayPointResponse{
		Message:    "Waypoint created and photos are being processed",
		WayPointID: wpID,
	}, nil
}
