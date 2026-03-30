package controller

import (
	"fmt"
	"user-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UploadAvatarRequest struct {
}

type UploadAvatarController struct {
	usecase usecase.UploadAvatarUseCase
}

type UploadAvatarResponse struct {
	Message   string `json:"message"`
	AvatarUrl string `json:"avatar_url"`
}

func NewUploadAvatarController(usecase usecase.UploadAvatarUseCase) *UploadAvatarController {
	return &UploadAvatarController{
		usecase: usecase,
	}
}

func (h *UploadAvatarController) Handle(fbrctx fiber.Ctx, req *UploadAvatarRequest) (*UploadAvatarResponse, error) {
	fmt.Println("sdsd")
	userIDStr := fbrctx.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid or missing user id")
	}

	fileHeader, err := fbrctx.FormFile("avatar")
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

	url, err := h.usecase.Execute(fbrctx.Context(), userID, file)
	if err != nil {
		return nil, err
	}
	return &UploadAvatarResponse{Message: "Avatar uploaded successfully", AvatarUrl: url}, nil
}
