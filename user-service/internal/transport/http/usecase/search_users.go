package usecase

import (
	"context"
	"user-service/internal/domain"

	"github.com/google/uuid"
)

type SearchUsersUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, limit, page int, query string) ([]domain.UserSummary, error)
}
type searchUsersUseCase struct {
	userRepository domain.UserRepository
}

func NewSearchUsersUseCase(userRepository domain.UserRepository) SearchUsersUseCase {
	return &searchUsersUseCase{
		userRepository: userRepository,
	}
}

func (uc *searchUsersUseCase) Execute(ctx context.Context, userID uuid.UUID, limit, page int, query string) ([]domain.UserSummary, error) {

	users, err := uc.userRepository.SearchUsers(ctx, query, userID, page, limit)
	if err != nil {
		return nil, err
	}

	return users, nil
}
