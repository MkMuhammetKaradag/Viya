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
		users.Post("/upload-banner", handler.HandleWithFiber[controller.UploadBannerRequest, controller.UploadBannerResponse](h.User.UploadBanner))
		users.Put("/update-profile", handler.HandleWithFiber[controller.UpdateProfileRequest, controller.UpdateProfileResponse](h.User.UpdateProfile))
		users.Get("/me", handler.HandleWithFiber[controller.GetMeRequest, controller.GetMeResponse](h.User.GetMe))
		users.Get("/search", handler.HandleWithFiber[controller.SearchUsersRequest, controller.SearchUsersResponse](h.User.SearchUsers))
		users.Get("/profile/:user_id", handler.HandleWithFiber[controller.GetUserRequest, controller.GetUserResponse](h.User.GetUser))
	}

}
