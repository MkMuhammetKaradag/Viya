package usecase

import (
	"auth-service/internal/domain"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

type AllSignOutUseCase interface {
	Execute(fbrCtx fiber.Ctx) error
}

type allSignOutUseCase struct {
	sessionRepo domain.SessionRepository
}

func NewAllSignOutUseCase(sessionRepo domain.SessionRepository) AllSignOutUseCase {
	return &allSignOutUseCase{
		sessionRepo: sessionRepo,
	}
}

func (uc *allSignOutUseCase) Execute(fbrCtx fiber.Ctx) error {
	userID := fbrCtx.Get("X-User-ID")
	fmt.Println("user:", userID)
	if err := uc.sessionRepo.DeleteAllSession(fbrCtx.Context(), userID); err != nil {
		return err

	}

	fbrCtx.Cookie(&fiber.Cookie{
		Name:     "Session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})
	return nil
}
