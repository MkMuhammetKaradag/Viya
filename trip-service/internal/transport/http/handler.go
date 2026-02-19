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
	Add       *controller.AddWayPointController
	AddPhotos *controller.AddWayPointPhotosController
}

func NewHandlers(repo domain.TripRepository, imgSvc domain.ImageService) *Handlers {
	return &Handlers{
		Trip: &tripHandlers{
			// UseCase ve Controller birleşimi
			Create: controller.NewCreateTripController(usecase.NewCreateTripUseCase(repo)),
		},
		WayPoint: &waypointHandlers{
			Add:       controller.NewAddWaypointController(usecase.NewAddWayPointUseCase(repo)),
			AddPhotos: controller.NewAddWayPointPhotosController(usecase.NewAddWayPointPhotosUseCase(repo, imgSvc)),
		},
	}
}
