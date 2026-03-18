package usecase

import (
	"context"
	"fmt"
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

func (u *createUserUseCase) Execute(ctx context.Context, userID uuid.UUID, userName, email string) error {
	fmt.Println("save user email", email, " userid:", userID, "-user name:", userName)

	return nil
}
