package http

import (
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/controller"
	"trip-service/internal/transport/http/usecase"
)

type Handlers struct {
	Trip     *tripHandlers
	WayPoint *waypointHandlers
}

type tripHandlers struct {
	Create *controller.CreateTripController
}

type waypointHandlers struct {
	Add *controller.AddWayPointController
}

func NewHandlers(repo domain.TripRepository) *Handlers {
	return &Handlers{
		Trip: &tripHandlers{
			// UseCase ve Controller birleşimi
			Create: controller.NewCreateTripController(usecase.NewCreateTripUseCase(repo)),
		},
		WayPoint: &waypointHandlers{
			Add: controller.NewAddWaypointController(usecase.NewAddWayPointUseCase(repo)),
		},
	}
}
