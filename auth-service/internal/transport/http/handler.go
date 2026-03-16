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
	ResetPassword  *controller.ResetPasswordController
	SignOut        *controller.SignOutController
	AllSignOut     *controller.AllSignOutController
}

func NewHandlers(repo domain.AuthRepository, sessionRepo domain.SessionRepository, rabbitClient domain.RabbitMQClient) *Handler {
	return &Handler{
		Auth: &AuthHandlers{
			SignUp:         controller.NewSignUpController(usecase.NewSignupUseCase(repo, rabbitClient)),
			SignIn:         controller.NewSignInController(usecase.NewSignInUseCase(repo, sessionRepo)),
			ForgotPassword: controller.NewForgotPasswordController(usecase.NewForgotPasswordUseCase(repo, sessionRepo)),
			ResetPassword:  controller.NewResetPasswordController(usecase.NewResetPasswordUseCase(repo, sessionRepo)),
			SignOut:        controller.NewSignOutController(usecase.NewSignOutUseCase(sessionRepo)),
			AllSignOut:     controller.NewAllSignOutontroller(usecase.NewAllSignOutUseCase(sessionRepo)),
		},
	}
}
