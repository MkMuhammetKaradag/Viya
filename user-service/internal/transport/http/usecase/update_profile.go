package usecase

import (
	"context"
	"fmt"
	"user-service/internal/domain"

	"github.com/google/uuid"
)

type UpdateProfileUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, params domain.UpdateProfileParams) error
}
type updateProfileUseCase struct {
	userRepository domain.UserRepository
}

func NewUpdateProfileUseCase(userRepository domain.UserRepository) UpdateProfileUseCase {
	return &updateProfileUseCase{
		userRepository: userRepository,
	}
}

func (uc *updateProfileUseCase) Execute(ctx context.Context, userID uuid.UUID, params domain.UpdateProfileParams) error {

	err := uc.userRepository.UpdateProfile(ctx, userID, params)
	if err != nil {
		return err
	}
	fmt.Println(params)
	return nil
}
