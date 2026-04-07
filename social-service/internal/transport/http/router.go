package http

import (
	"github.com/gofiber/fiber/v3"
)

type Router struct {
	handlers *Handlers
}

func NewRouter(handlers *Handlers) *Router {
	return &Router{handlers: handlers}
}

func (r *Router) Register(app *fiber.App) {

	api := app.Group("/api/v1")
	// social Endpoints
	api.Post("/social", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "active",
			"service": "trip-service",
		})
	})

}
