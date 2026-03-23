package controller

import (
	"user-service/internal/domain"
	"user-service/internal/transport/http/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UpdateProfileRequest struct {
	FirstName   *string  `json:"first_name"`
	LastName    *string  `json:"last_name"`
	Bio         *string  `json:"bio"`
	Location    *string  `json:"location"`
	Website     *string  `json:"website"`
	Preferences []string `json:"preferences"`
}

type UpdateProfileController struct {
	usecase usecase.UpdateProfileUseCase
}

type UpdateProfileResponse struct {
	Message string `json:"message"`
}

func NewUpdateProfileController(usecase usecase.UpdateProfileUseCase) *UpdateProfileController {
	return &UpdateProfileController{
		usecase: usecase,
	}
}

func (h *UpdateProfileController) Handle(fbrctx fiber.Ctx, req *UpdateProfileRequest) (*UpdateProfileResponse, error) {
	userIDStr := fbrctx.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid or missing user id")
	}
	parans := domain.UpdateProfileParams{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Bio:         req.Bio,
		Location:    req.Location,
		Website:     req.Website,
		Preferences: req.Preferences,
	}
	if err := h.usecase.Execute(fbrctx.Context(), userID, parans); err != nil {
		return nil, err
	}
	return &UpdateProfileResponse{Message: "user updated successfully"}, nil
}
