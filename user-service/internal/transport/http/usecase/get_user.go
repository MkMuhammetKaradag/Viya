package usecase

import (
	"context"
	"fmt"
	"user-service/internal/domain"

	"github.com/google/uuid"
)

type GetUserUseCase interface {
	Execute(ctx context.Context, currentUserId, userID uuid.UUID) (*domain.UserSummary, error)
}
type getUserUseCase struct {
	userRepository domain.UserRepository
}

func NewGetUserUseCase(userRepository domain.UserRepository) GetUserUseCase {
	return &getUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *getUserUseCase) Execute(ctx context.Context, currentUserID, userID uuid.UUID) (*domain.UserSummary, error) {
	fmt.Println("userid:", userID)
	user, err := uc.userRepository.GetUser(ctx, currentUserID, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
