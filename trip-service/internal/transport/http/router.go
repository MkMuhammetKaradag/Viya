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
	}

	// WAYPOINTS
	waypoints := api.Group("/waypoints")
	{
		waypoints.Post("/", handler.HandleWithFiber[controller.AddWayPointRequest, controller.AddWayPointResponse](h.WayPoint.Add))
	}

}
