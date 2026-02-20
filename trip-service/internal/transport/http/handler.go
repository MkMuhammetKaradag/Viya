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
	Get    *controller.GetTripController
}

type waypointHandlers struct {
	Add       *controller.AddWayPointController
	AddPhotos *controller.AddWayPointPhotosController
	Delete    *controller.DeleteWaypointController
}

func NewHandlers(repo domain.TripRepository, imgSvc domain.ImageService, worker domain.Worker) *Handlers {
	return &Handlers{
		Trip: &tripHandlers{
			// UseCase ve Controller birleşimi
			Create: controller.NewCreateTripController(usecase.NewCreateTripUseCase(repo)),
			Get:    controller.NewGetTripController(usecase.NewGetTripUseCase(repo)),
		},
		WayPoint: &waypointHandlers{
			Add:       controller.NewAddWaypointController(usecase.NewAddWayPointUseCase(repo)),
			AddPhotos: controller.NewAddWayPointPhotosController(usecase.NewAddWayPointPhotosUseCase(repo, imgSvc, worker)),
			Delete:    controller.NewDeleteWaypointController(usecase.NewDeleteWaypointUseCase(repo)),
		},
	}
}
