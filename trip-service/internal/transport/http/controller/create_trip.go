package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CreateTripRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	ID       uuid.UUID `params:"id,omitempty" validate:"required" `
	Search   string    `query:"search,omitempty"`
	Title    string    `json:"title" validate:"required,min=3"`
	Desc     string    `json:"desc,omitempty"`
	IsActive bool      `json:"is_active"`
}

type CreateTripResponse struct {
	Message string `json:"message"`
}

type CreateTripController struct {
}

func NewCreateTripController() *CreateTripController {
	return &CreateTripController{}
}

func (c *CreateTripController) Handle(fiberCtx fiber.Ctx, req *CreateTripRequest) (*CreateTripResponse, error) {
	// Burada UseCase'i çağırıp trip oluşturma işlemini yapacağız
	// Şimdilik sadece dummy response dönüyoruz
	cookie := fiberCtx.Cookies("session_id")
	fmt.Println("cookie", cookie)
	fmt.Println("data", req)
	return &CreateTripResponse{Message: "Trip created successfully!"}, nil
}
