package usecase

import (
	"auth-service/internal/domain"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

type SignOutUseCase interface {
	Execute(ctx fiber.Ctx) error
}
type signOutUseCase struct {
	sessionRepository domain.SessionRepository
}

func NewSignOutUseCase(repository domain.SessionRepository) SignOutUseCase {
	return &signOutUseCase{
		sessionRepository: repository,
	}
}

func (u *signOutUseCase) Execute(ctx fiber.Ctx) error {
	sessionID := ctx.Cookies("session_id")
	fmt.Println(sessionID)
	err := u.sessionRepository.DeleteSession(ctx.Context(), sessionID)
	// Hata olsa bile tarayıcıdaki çerezi temizle
	ctx.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})
	if err != nil {
		return err
	}
	return nil
}
