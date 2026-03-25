package usecase

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, userName, email string) error
}
type createUserUseCase struct {
	repository domain.TripRepository
}

func NewUserCreatedUseCase(repository domain.TripRepository) CreateUserUseCase {
	return &createUserUseCase{
		repository: repository,
	}
}

func (uc *createUserUseCase) Execute(ctx context.Context, userID uuid.UUID, userName, email string) error {
	return uc.repository.CreateUser(ctx, userID, userName, email)
}
