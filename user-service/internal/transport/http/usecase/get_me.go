package usecase

import (
	"context"
	"user-service/internal/domain"

	"github.com/google/uuid"
)

type GetMeUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}
type getMeUseCase struct {
	userRepository domain.UserRepository
}

func NewGetMeUseCase(userRepository domain.UserRepository) GetMeUseCase {
	return &getMeUseCase{
		userRepository: userRepository,
	}
}

func (uc *getMeUseCase) Execute(ctx context.Context, userID uuid.UUID) (*domain.User, error) {

	user, err := uc.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
