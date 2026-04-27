package controller

import (
	"fmt"
	"trip-service/internal/domain"
	"trip-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetCommentRepliseRequest struct {
	Page      int       `query:"page" validate:"omitempty,min=1"`
	Limit     int       `query:"limit" validate:"omitempty,min=1,max=50"`
	CommentID uuid.UUID `uri:"comment_id" validate:"required"`
}

type GetCommentRepliseResponse struct {
	Comments []domain.Comment `json:"comments"`
}
type GetCommentRepliseController struct {
	usecase usecase.GetCommentRepliseUseCase
}

func NewGetCommentRepliseController(usecase usecase.GetCommentRepliseUseCase) *GetCommentRepliseController {
	return &GetCommentRepliseController{
		usecase: usecase,
	}
}

func (c *GetCommentRepliseController) Handle(fbrCtx fiber.Ctx, req *GetCommentRepliseRequest) (*GetCommentRepliseResponse, error) {
	rawUserID := fbrCtx.Get("X-User-ID")

	var currentUserID uuid.UUID
	if rawUserID != "" {
		parsedID, err := uuid.Parse(rawUserID)
		if err == nil {
			currentUserID = parsedID
		}
	}

	fmt.Println("limit:", req.Limit, "commentid:", req.CommentID, "page:", req.Page)
	comments, err := c.usecase.Execute(fbrCtx.Context(), currentUserID, req.CommentID, req.Limit, req.Page)
	if err != nil {
		return nil, err
	}

	return &GetCommentRepliseResponse{Comments: comments}, nil
}
