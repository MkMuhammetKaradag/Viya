package controller

import (
	"auth-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
)

type ForgotPasswordRequest struct {
	Identifier string `json:"identifier" validate:"required"`
}

type ForgotPasswordResponse struct {
	Message   string `json:"message"`
	SessionID string `json:"session_id"`
}

type ForgotPasswordController struct {
	usecase usecase.ForgotPasswordUseCase
}

func NewForgotPasswordController(usecase usecase.ForgotPasswordUseCase) *ForgotPasswordController {
	return &ForgotPasswordController{usecase: usecase}
}

func (c *ForgotPasswordController) Handle(ctx fiber.Ctx, req *ForgotPasswordRequest) (*ForgotPasswordResponse, error) {
	clientOS := ctx.Get("X-Platform")

	// Karar mekanizması
	platformType := "web"
	if clientOS == "ios" || clientOS == "android" {
		platformType = "mobile"
	}
	sessionID, err := c.usecase.Execute(ctx.Context(), req.Identifier, platformType)
	if err != nil {
		return nil, err
	}
	return &ForgotPasswordResponse{Message: "If an account with that identifier exists, a password reset link has been sent.", SessionID: sessionID}, nil
}
