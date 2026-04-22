package controller

import (
	"user-service/internal/domain"
	"user-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type SearchUsersRequest struct {
	Query string `query:"query" validate:"omitempty,min=3"`
	Page  int    `query:"page" validate:"omitempty,min=1"`
	Limit int    `query:"limit" validate:"omitempty,min=1,max=50"`
}

type SearchUsersController struct {
	usecase usecase.SearchUsersUseCase
}

type SearchUsersResponse struct {
	Users []domain.UserSummary `json:"users"`
}

func NewSearchUsersController(usecase usecase.SearchUsersUseCase) *SearchUsersController {
	return &SearchUsersController{
		usecase: usecase,
	}
}

func (c *SearchUsersController) Handle(fbrCtx fiber.Ctx, req *SearchUsersRequest) (*SearchUsersResponse, error) {
	rawUserID := fbrCtx.Get("X-User-ID")

	var currentUserID uuid.UUID
	if rawUserID != "" {
		parsedID, err := uuid.Parse(rawUserID)
		if err == nil {
			currentUserID = parsedID
		}
	}

	users, err := c.usecase.Execute(fbrCtx.Context(), currentUserID, req.Limit, req.Page, req.Query)
	if err != nil {
		return nil, err
	}
	return &SearchUsersResponse{Users: users}, nil
}
