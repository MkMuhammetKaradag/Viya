package controller

import (
	"fmt"
	"trip-service/internal/transport/http/usecase"

	"trip-service/internal/domain"

	"github.com/gofiber/fiber/v3"
)

type SearchCategoriesRequest struct {
	SearchQuery string `query:"search_query" validate:"required"`
}

type SearchCategoriesResponse struct {
	Categories []domain.Category `json:"categories,omitempty"`
}

type SearchCategoriesController struct {
	usecase usecase.SearchCategoriesUseCase
}

func NewSearchCategoriesController(usecase usecase.SearchCategoriesUseCase) *SearchCategoriesController {
	return &SearchCategoriesController{
		usecase: usecase,
	}
}

func (c *SearchCategoriesController) Handle(fiberCtx fiber.Ctx, req *SearchCategoriesRequest) (*SearchCategoriesResponse, error) {
	fmt.Println("Search query:", req.SearchQuery)
	categories, err := c.usecase.Execute(fiberCtx, req.SearchQuery)
	if err != nil {
		return nil, err
	}

	return &SearchCategoriesResponse{Categories: categories}, nil
}
