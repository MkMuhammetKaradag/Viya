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
		trips.Get("/user/:user_id", handler.HandleWithFiber[controller.GetUserTripsRequest, controller.GetUserTripsResponse](h.Trip.GetUserTrips))
		trips.Get("/me", handler.HandleWithFiber[controller.GetMeTripsRequest, controller.GetMeTripsResponse](h.Trip.GetMeTrips))

		trips.Get("/explore", handler.HandleWithFiber[controller.GetExploreTripsRequest, controller.GetExploreTripsResponse](h.Trip.GetExploreTrips))
		trips.Get("/home-feed", handler.HandleWithFiber[controller.GetHomeFeedTripsRequest, controller.GetHomeFeedTripsResponse](h.Trip.GetHomeFeedTrips))
		trips.Get("/liked", handler.HandleWithFiber[controller.GetLikedTripsRequest, controller.GetLikedTripsResponse](h.Trip.GetLikedTrips))
		trips.Patch("/:trip_id/like", handler.HandleWithFiber[controller.ToggleTripLikeRequest, controller.ToggleTripLikeResponse](h.Trip.ToggleTripLike))
		trips.Post("/:trip_id/fork", handler.HandleWithFiber[controller.ForkTripRequest, controller.ForkTripResponse](h.Trip.ForkedTrip))
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
	comments := api.Group("/comments")
	{
		comments.Post("/", handler.HandleWithFiber[controller.CreateCommentRequest, controller.CreateCommentResponse](h.Comment.Create))
		comments.Get("/trip/:trip_id", handler.HandleWithFiber[controller.GetTripCommentsRequest, controller.GetTripCommentsResponse](h.Comment.GetTripComments))
		comments.Get("/replies/:comment_id", handler.HandleWithFiber[controller.GetCommentRepliseRequest, controller.GetCommentRepliseResponse](h.Comment.GetCommentReplies))
	}

}
