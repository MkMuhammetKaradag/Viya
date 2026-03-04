package usecase

import (
	"auth-service/internal/domain"
	"context"
	"fmt"
)

type SignupUseCase interface {
	Execute(ctx context.Context, username, email, password string) error
}
type signupUseCase struct {
	repo domain.AuthRepository
}

func NewSignupUseCase(repo domain.AuthRepository) SignupUseCase {
	return &signupUseCase{repo: repo}
}

func (uc *signupUseCase) Execute(ctx context.Context, username, email, password string) error {
	fmt.Println(
		"Executing SignupUseCase with username:", username,
		"email:", email,
		"password:", password,
	)
	if err := uc.repo.SignUp(ctx, username, email, password); err != nil {
		return fmt.Errorf("signup error: %w", err)
	}

	return nil
}
