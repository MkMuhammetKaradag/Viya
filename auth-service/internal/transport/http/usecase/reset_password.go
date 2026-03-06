package usecase

import (
	"auth-service/internal/domain"
	"context"
)

type ResetPasswordUseCase interface {
	Execute(ctx context.Context, newPassword, token string) error
}

type resetPasswordUseCase struct {
	repo domain.AuthRepository
}

func NewResetPasswordUseCase(repo domain.AuthRepository) ResetPasswordUseCase {
	return &resetPasswordUseCase{
		repo: repo,
	}
}

func (uc *resetPasswordUseCase) Execute(ctx context.Context, newPassword, token string) error {

	err := uc.repo.ResetPassword(ctx, token, newPassword)
	if err != nil {
		return err

	}

	return nil
}
