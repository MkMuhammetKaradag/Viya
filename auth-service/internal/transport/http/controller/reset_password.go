package controller

import (
	"auth-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
)

type ResetPasswordRequest struct {
	Token       string `json:"token" v`
	SessionID   string `json:"session_id" `
	Code        string `json:"code" `
	NewPassword string `json:"new_password" validate:"required,min=8,max=32"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

type ResetPasswordController struct {
	usecase usecase.ResetPasswordUseCase
}

func NewResetPasswordController(usecase usecase.ResetPasswordUseCase) *ResetPasswordController {
	return &ResetPasswordController{
		usecase: usecase,
	}
}

func (c *ResetPasswordController) Handle(fbrCtx fiber.Ctx, req *ResetPasswordRequest) (*ResetPasswordResponse, error) {
	clientOS := fbrCtx.Get("X-Platform")
	platformType := "web"
	if clientOS == "ios" || clientOS == "android" {
		platformType = "mobile"
	}
	err := c.usecase.Execute(fbrCtx.Context(), req.NewPassword, req.Token, req.Code, req.SessionID, platformType)
	if err != nil {
		return nil, err
	}
	return &ResetPasswordResponse{
		Message: "password reset success ",
	}, nil
}
