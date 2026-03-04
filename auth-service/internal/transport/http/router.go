package http

import (
	"auth-service/internal/handler"
	"auth-service/internal/transport/http/controller"

	"github.com/gofiber/fiber/v3"
)

type Router struct {
	handler *Handler
}

func NewRouter(handler *Handler) *Router {
	return &Router{handler: handler}
}

func (r *Router) Register(app *fiber.App) {
	h := r.handler
	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.Post("/signup", handler.HandleBasic[controller.SignUpRequest, controller.SignUpResponse](h.Auth.SignUp))
		auth.Post("/signin", handler.HandleWithFiber[controller.SignInRequest, controller.SignInResponse](h.Auth.SignIn))
	}

}
