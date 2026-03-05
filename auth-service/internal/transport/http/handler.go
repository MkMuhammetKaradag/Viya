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
	SignUp         *controller.SignUpController
	SignIn         *controller.SignInController
	ForgotPassword *controller.ForgotPasswordController
}

func NewHandlers(repo domain.AuthRepository, sessionRepo domain.SessionRepository) *Handler {
	return &Handler{
		Auth: &AuthHandlers{
			SignUp:         controller.NewSignUpController(usecase.NewSignupUseCase(repo)),
			SignIn:         controller.NewSignInController(usecase.NewSignInUseCase(repo, sessionRepo)),
			ForgotPassword: controller.NewForgotPasswordController(usecase.NewForgotPasswordUseCase(repo, sessionRepo)),
		},
	}
}
