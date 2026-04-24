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
	Comment  *commentHandlers
}

type tripHandlers struct {
	Create          *controller.CreateTripController
	Get             *controller.GetTripController
	GetUserTrips    *controller.GetUserTripsController
	GetLikedTrips   *controller.GetLikedTripsController
	GetExploreTrips *controller.GetExploreTripsController
	ToggleTripLike  *controller.ToggleTripLikeController
}
type commentHandlers struct {
	Create          *controller.CreateCommentController
	GetTripComments *controller.GetTripCommentsController
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

func NewHandlers(repo domain.TripRepository, imgSvc domain.ImageService, worker domain.Worker, moderationService domain.ModerationService) *Handlers {
	return &Handlers{
		Trip: &tripHandlers{
			// UseCase ve Controller birleşimi
			Create:          controller.NewCreateTripController(usecase.NewCreateTripUseCase(repo, worker)),
			Get:             controller.NewGetTripController(usecase.NewGetTripUseCase(repo, worker)),
			GetUserTrips:    controller.NewGetUserTripsController(usecase.NewGetUserTripsUseCase(repo)),
			GetExploreTrips: controller.NewGetExploreTripsController(usecase.NewGetExploreTripsUseCase(repo, worker)),
			ToggleTripLike:  controller.NewToggleTripLikeController(usecase.NewToggleTripLikeUseCase(repo, worker)),
			GetLikedTrips:   controller.NewGetLikedTripsController(usecase.NewGetLikedTripsUseCase(repo)),
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
		Comment: &commentHandlers{
			Create:          controller.NewCreateCommentController(usecase.NewCreateCommentUseCase(repo, moderationService)),
			GetTripComments: controller.NewGetTripCommentsController(usecase.NewGetTripCommentsUseCase(repo)),
		},
	}
}
