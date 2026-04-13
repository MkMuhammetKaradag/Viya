package usecase

import (
	"context"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

type UpdatedUserUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, isPrivate *bool, avatarURL *string) error
}
type updatedUserUseCase struct {
	repository domain.TripRepository
}

func NewUserUpdateddUseCase(repository domain.TripRepository) UpdatedUserUseCase {
	return &updatedUserUseCase{
		repository: repository,
	}
}

func (uc *updatedUserUseCase) Execute(ctx context.Context, userID uuid.UUID, isPrivate *bool, avatarURL *string) error {
	return uc.repository.UpdateUserSocialInfo(ctx, userID, isPrivate, avatarURL)
}
