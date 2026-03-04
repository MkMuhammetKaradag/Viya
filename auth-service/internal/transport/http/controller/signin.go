package controller

import (
	"auth-service/internal/transport/http/usecase"
	"time"

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
	sessionID, err := c.usecase.Execute(fbrCtx.Context(), req.Identifier, req.Password)
	if err != nil {
		return nil, err
	}
	fbrCtx.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,  // JavaScript tarafından okunamaz (XSS koruması)
		Secure:   true,  // Sadece HTTPS üzerinden gönderilir
		SameSite: "Lax", // CSRF saldırılarına karşı koruma sağlar
		Path:     "/",
	})
	return &SignInResponse{
		Message: "Sign in successful",
	}, nil
}
