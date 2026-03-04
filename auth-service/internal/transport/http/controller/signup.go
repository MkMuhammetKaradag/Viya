package controller

import (
	"auth-service/internal/transport/http/usecase"
	"context"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Message string `json:"message"`
}

type SignUpController struct {
	usecase usecase.SignupUseCase
}

func NewSignUpController(usecase usecase.SignupUseCase) *SignUpController {
	return &SignUpController{usecase: usecase}
}

func (c *SignUpController) Handle(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error) {
	err := c.usecase.Execute(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &SignUpResponse{Message: "User registered successfully"}, nil
}
