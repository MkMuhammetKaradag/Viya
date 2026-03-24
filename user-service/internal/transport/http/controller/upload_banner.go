package controller

import (
	"fmt"
	"user-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UploadBannerRequest struct {
}

type UploadBannerController struct {
	usecase usecase.UploadBannerUseCase
}

type UploadBannerResponse struct {
	Message string `json:"message"`
}

func NewUploadBannerController(usecase usecase.UploadBannerUseCase) *UploadBannerController {
	return &UploadBannerController{
		usecase: usecase,
	}
}

func (h *UploadBannerController) Handle(fbrctx fiber.Ctx, req *UploadBannerRequest) (*UploadBannerResponse, error) {
	fmt.Println("sdsd")
	userIDStr := fbrctx.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid or missing user id")
	}

	fileHeader, err := fbrctx.FormFile("banner")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	file, err := fileHeader.Open()
	if err != nil {
		fmt.Println("open", err)

		return nil, err
	}
	defer file.Close()

	if err := h.usecase.Execute(fbrctx.Context(), userID, file); err != nil {
		return nil, err
	}
	return &UploadBannerResponse{Message: "Avatar Banner successfully"}, nil
}
