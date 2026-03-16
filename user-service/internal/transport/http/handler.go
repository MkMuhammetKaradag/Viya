package http

import "user-service/internal/domain"

type Handler struct {
	User *UserHandlers
}

type UserHandlers struct {
}

func NewHandlers(userRepo domain.UserRepository) *Handler {
	return &Handler{
		User: &UserHandlers{},
	}
}
