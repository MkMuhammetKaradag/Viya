package http

import (
	"auth-service/internal/domain"
	"auth-service/internal/transport/http/controller"
	"auth-service/internal/transport/http/usecase"
)

type Handler struct {
	Auth *AuthHandlers
}

type AuthHandlers struct {
	SignUp *controller.SignUpController
}

func NewHandlers(repo domain.AuthRepository) *Handler {
	return &Handler{
		Auth: &AuthHandlers{
			SignUp: controller.NewSignUpController(usecase.NewSignupUseCase(repo)),
		},
	}
}
