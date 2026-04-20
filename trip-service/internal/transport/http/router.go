package http

import (
	"trip-service/internal/handler" // HandleWithFiber'ın olduğu yer
	"trip-service/internal/transport/http/controller"

	"github.com/gofiber/fiber/v3"
)

type Router struct {
	handlers *Handlers
}

func NewRouter(handlers *Handlers) *Router {
	return &Router{handlers: handlers}
}

func (r *Router) Register(app *fiber.App) {
	h := r.handlers

	api := app.Group("/api/v1")

	// TRIPS
	trips := api.Group("/trips")
	{
		trips.Post("/", handler.HandleWithFiber[controller.CreateTripRequest, controller.CreateTripResponse](h.Trip.Create))
		trips.Get("/me", handler.HandleWithFiber[controller.GetUserTripsRequest, controller.GetUserTripsResponse](h.Trip.GetUserTrips))
		trips.Get("/explore", handler.HandleWithFiber[controller.GetExploreTripsRequest, controller.GetExploreTripsResponse](h.Trip.GetExploreTrips))
		trips.Patch("/:trip_id/like", handler.HandleWithFiber[controller.ToggleTripLikeRequest, controller.ToggleTripLikeResponse](h.Trip.ToggleTripLike))
		trips.Get("/:trip_id", handler.HandleWithFiber[controller.GetTripRequest, controller.GetTripResponse](h.Trip.Get))

	}

	// WAYPOINTS
	waypoints := api.Group("/waypoints")
	{
		waypoints.Post("/", handler.HandleWithFiber[controller.AddWayPointRequest, controller.AddWayPointResponse](h.WayPoint.Add))
		//waypoints.Post("/:waypoint_id/photos", handler.HandleWithFiber[controller.AddWayPointPhotosRequest, controller.AddWayPointPhotosResponse](h.WayPoint.AddPhotos))
		waypoints.Delete("/:waypoint_id", handler.HandleWithFiber[controller.DeleteWaypointRequest, controller.DeleteWaypointResponse](h.WayPoint.Delete))
		waypoints.Patch("/:waypoint_id/reorder", handler.HandleWithFiber[controller.ReorderRequest, controller.ReorderResponse](h.WayPoint.Reorder))
		waypoints.Put("/:waypoint_id", handler.HandleWithFiber[controller.UpdateWaypointRequest, controller.UpdateWaypointResponse](h.WayPoint.Update))
	}
	categories := api.Group("/categories")
	{
		categories.Get("/search", handler.HandleWithFiber[controller.SearchCategoriesRequest, controller.SearchCategoriesResponse](h.Category.Search))
	}

}
