package http

import (
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/controller"
	"trip-service/internal/transport/http/usecase"
)

type Handlers struct {
	Trip     *tripHandlers
	WayPoint *waypointHandlers
	Category *categoryHandlers
}

type tripHandlers struct {
	Create          *controller.CreateTripController
	Get             *controller.GetTripController
	GetUserTrips    *controller.GetUserTripsController
	GetExploreTrips *controller.GetExploreTripsController
}
type categoryHandlers struct {
	Search *controller.SearchCategoriesController
}

type waypointHandlers struct {
	Add *controller.AddWayPointController

	//AddPhotos *controller.AddWayPointPhotosController
	Delete  *controller.DeleteWaypointController
	Reorder *controller.ReorderController
	Update  *controller.UpdateWaypointController
}

func NewHandlers(repo domain.TripRepository, imgSvc domain.ImageService, worker domain.Worker) *Handlers {
	return &Handlers{
		Trip: &tripHandlers{
			// UseCase ve Controller birleşimi
			Create:          controller.NewCreateTripController(usecase.NewCreateTripUseCase(repo, worker)),
			Get:             controller.NewGetTripController(usecase.NewGetTripUseCase(repo, worker)),
			GetUserTrips:    controller.NewGetUserTripsController(usecase.NewGetUserTripsUseCase(repo)),
			GetExploreTrips: controller.NewGetExploreTripsController(usecase.NewGetExploreTripsUseCase(repo, worker)),
		},
		WayPoint: &waypointHandlers{
			Add: controller.NewAddWaypointController(usecase.NewAddWayPointUseCase(repo, worker)),
			//AddPhotos: controller.NewAddWayPointPhotosController(usecase.NewAddWayPointPhotosUseCase(repo, imgSvc, worker)),
			Delete:  controller.NewDeleteWaypointController(usecase.NewDeleteWaypointUseCase(repo)),
			Reorder: controller.NewReorderController(usecase.NewReorderUseCase(repo)),
			Update:  controller.NewUpdateWaypointController(usecase.NewUpdateWaypointUseCase(repo)),
		},
		Category: &categoryHandlers{
			Search: controller.NewSearchCategoriesController(usecase.NewSearchCategoriesUseCase(repo)),
		},
	}
}
