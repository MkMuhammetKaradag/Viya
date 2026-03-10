package usecase

import (
	"auth-service/internal/domain"

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
	_ = ctx.Cookies("session_id")

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

	return nil
}
