package usecase

import (
	"context"
	"social-service/internal/domain"

	"github.com/google/uuid"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, userName, email string) error
}
type createUserUseCase struct {
	repository domain.SocialRepository
}

func NewUserCreatedUseCase(repository domain.SocialRepository) CreateUserUseCase {
	return &createUserUseCase{
		repository: repository,
	}
}

func (uc *createUserUseCase) Execute(ctx context.Context, userID uuid.UUID, userName, email string) error {
	return uc.repository.SaveUser(ctx, userID, userName, email)
}
