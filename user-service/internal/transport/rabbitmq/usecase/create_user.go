package usecase

import (
	"context"
	"user-service/internal/domain"

	"github.com/google/uuid"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, userName, email string) error
}
type createUserUseCase struct {
	repository domain.UserRepository
}

func NewUserCreatedUseCase(repository domain.UserRepository) CreateUserUseCase {
	return &createUserUseCase{
		repository: repository,
	}
}

func (uc *createUserUseCase) Execute(ctx context.Context, userID uuid.UUID, userName, email string) error {
	return uc.repository.CreateUser(ctx, userID, userName, email)
}
