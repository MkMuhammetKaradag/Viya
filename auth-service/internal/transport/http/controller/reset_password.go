package controller

import (
	"auth-service/internal/transport/http/usecase"
	"context"
)

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
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

func (c *ResetPasswordController) Handle(ctx context.Context, req *ResetPasswordRequest) (*ResetPasswordResponse, error) {

	err := c.usecase.Execute(ctx, req.NewPassword, req.Token)
	if err != nil {
		return nil, err
	}
	return &ResetPasswordResponse{
		Message: "password reset success ",
	}, nil
}
