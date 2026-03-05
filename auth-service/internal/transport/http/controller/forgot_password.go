package controller

import (
	"auth-service/internal/transport/http/usecase"
	"context"
)

type ForgotPasswordRequest struct {
	Identifier string `json:"identifier" validate:"required"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

type ForgotPasswordController struct {
	usecase usecase.ForgotPasswordUseCase
}

func NewForgotPasswordController(usecase usecase.ForgotPasswordUseCase) *ForgotPasswordController {
	return &ForgotPasswordController{usecase: usecase}
}

func (c *ForgotPasswordController) Handle(ctx context.Context, req *ForgotPasswordRequest) (*ForgotPasswordResponse, error) {

	err := c.usecase.Execute(ctx, req.Identifier)
	if err != nil {
		return nil, err
	}
	return &ForgotPasswordResponse{Message: "If an account with that identifier exists, a password reset link has been sent."}, nil
}
