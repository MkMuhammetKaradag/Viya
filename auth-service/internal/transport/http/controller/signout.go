// internal/user-service/transport/http/controller/signout.go
package controller

import (
	"auth-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
)

type SignOutRequest struct{}
type SignOutResponse struct {
	Message string `json:"message"`
}
type SignOutController struct {
	usecase usecase.SignOutUseCase
}

func NewSignOutController(usecase usecase.SignOutUseCase) *SignOutController {
	return &SignOutController{
		usecase: usecase,
	}
}

func (h *SignOutController) Handle(fbrCtx fiber.Ctx, req *SignOutRequest) (*SignOutResponse, error) {
	err := h.usecase.Execute(fbrCtx)
	if err != nil {
		return nil, err
	}
	//fbrCtx.("X-Invalidate-Session", "true")
	return &SignOutResponse{Message: "logout successfully"}, nil
}
