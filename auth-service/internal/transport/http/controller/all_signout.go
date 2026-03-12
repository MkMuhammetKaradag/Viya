package controller

import (
	"auth-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
)

type AllSignOutRequest struct {
}

type AllSignOutResponse struct {
	Message string `json:"string"`
}
type AllSignOutController struct {
	usecase usecase.AllSignOutUseCase
}

func NewAllSignOutontroller(usecase usecase.AllSignOutUseCase) *AllSignOutController {
	return &AllSignOutController{
		usecase: usecase,
	}
}

func (c *AllSignOutController) Handle(fbrCtx fiber.Ctx, req *AllSignOutRequest) (*AllSignOutResponse, error) {
	err := c.usecase.Execute(fbrCtx)
	if err != nil {
		return nil, err
	}

	return &AllSignOutResponse{
		Message: "All SignOut Successful",
	}, nil
}
