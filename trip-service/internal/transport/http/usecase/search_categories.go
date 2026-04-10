package usecase

import (
	"context"
	"trip-service/internal/domain"
)

type SearchCategoriesUseCase interface {
	Execute(ctx context.Context, searchQuery string) ([]domain.Category, error)
}

type searchCategoriesUseCase struct {
	tripRepo domain.TripRepository
}

func NewSearchCategoriesUseCase(tripRepo domain.TripRepository) SearchCategoriesUseCase {
	return &searchCategoriesUseCase{tripRepo: tripRepo}
}

func (uc *searchCategoriesUseCase) Execute(ctx context.Context, searchQuery string) ([]domain.Category, error) {
	id, err := uc.tripRepo.SearchCategories(ctx, searchQuery)
	if err != nil {
		return nil, err
	}

	return id, nil
}
