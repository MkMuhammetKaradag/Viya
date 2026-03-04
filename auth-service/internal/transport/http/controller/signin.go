package controller

import (
	"auth-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
)

type SignInRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required,min=8"`
}

type SignInResponse struct {
	Message string `json:"message"`
}
type SignInController struct {
	usecase usecase.SignInUseCase
}

func NewSignInController(usecase usecase.SignInUseCase) *SignInController {
	return &SignInController{
		usecase: usecase,
	}
}

func (c *SignInController) Handle(fbrCtx fiber.Ctx, req *SignInRequest) (*SignInResponse, error) {
	err := c.usecase.Execute(fbrCtx.Context(), req.Identifier, req.Password)
	if err != nil {
		return nil, err
	}
	return &SignInResponse{
		Message: "Sign in successful",
	}, nil
}
