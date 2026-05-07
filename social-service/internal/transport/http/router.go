package http

import (
	"social-service/internal/handler"
	"social-service/internal/transport/http/controller"

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

	// social Endpoints
	social := api.Group("/social")
	{
		social.Post("/follow/:target_user_id", handler.HandleWithFiber[controller.FollowUserRequest, controller.FollowUserResponse](h.Social.FollowUser))
		social.Post("/unfollow/:target_user_id", handler.HandleWithFiber[controller.UnFollowUserRequest, controller.UnFollowUserResponse](h.Social.UnFollowUser))
		social.Delete("/remove-follower/:follower_id", handler.HandleWithFiber[controller.RemoveFollowerRequest, controller.RemoveFollowerResponse](h.Social.RemoveFollower))
		social.Post("/follow-request/:follower_id", handler.HandleWithFiber[controller.FollowRequestRequest, controller.FollowRequestResponse](h.Social.FollowRequest))
		social.Get("/pending-requests", handler.HandleWithFiber[controller.PendingRequestsRequest, controller.PendingRequestsResponse](h.Social.PendingRequests))
		social.Get("/sent-follow-requests", handler.HandleWithFiber[controller.GetSentFollowRequestsRequest, controller.GetSentFollowRequestsResponse](h.Social.GetSentFollowRequests))
		social.Post("/block/:target_user_id", handler.HandleWithFiber[controller.BlockUserRequest, controller.BlockUserResponse](h.Social.BlockUser))
		social.Delete("/unblock/:target_user_id", handler.HandleWithFiber[controller.UnblockUserRequest, controller.UnblockUserResponse](h.Social.UnblockUser))
	}

}
