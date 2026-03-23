package http

import (
	"user-service/internal/handler"
	"user-service/internal/transport/http/controller"

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

	users := api.Group("/users")
	{
		users.Post("/upload-avatar", handler.HandleWithFiber[controller.UploadAvatarRequest, controller.UploadAvatarResponse](h.User.UploadAvatar))
		users.Put("/update-profile", handler.HandleWithFiber[controller.UpdateProfileRequest, controller.UpdateProfileResponse](h.User.UpdateProfile))
	}

}
