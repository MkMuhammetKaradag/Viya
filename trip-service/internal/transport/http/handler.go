package http

import (
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/controller"
)

type Handlers struct {
	Trip *tripHandlers
}

type tripHandlers struct {
	Create *controller.CreateTripController
}

func NewHandlers(repo domain.TripRepository) *Handlers {
	return &Handlers{
		Trip: &tripHandlers{
			// UseCase ve Controller birleşimi
			Create: controller.NewCreateTripController(),
		},
	}
}
