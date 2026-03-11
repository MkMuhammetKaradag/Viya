package usecase

import (
	"auth-service/internal/domain"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
)

type SignInUseCase interface {
	Execute(fiberCtx fiber.Ctx, identifier, password string) (string, error)
}

type signInUsecase struct {
	repo        domain.AuthRepository
	sessionRepo domain.SessionRepository
}

func NewSignInUseCase(repo domain.AuthRepository, sessionRepo domain.SessionRepository) SignInUseCase {
	return &signInUsecase{
		repo:        repo,
		sessionRepo: sessionRepo,
	}
}

func (uc *signInUsecase) Execute(fiberCtx fiber.Ctx, identifier, password string) (string, error) {

	fmt.Println(
		"Executing SignInUseCase with identifier:", identifier,
		"and password:", password,
	)
	user, err := uc.repo.SignIn(fiberCtx.Context(), identifier, password)
	if err != nil {
		return "", fmt.Errorf("signin error: %w", err)
	}
	device := fiberCtx.Get("User-Agent")
	ip := fiberCtx.IP()
	userData := &domain.SessionData{
		UserID:    user.ID.String(), // In a real implementation, this would be the user's unique ID from the database
		CreatedAt: time.Now(),
		Device:    device,
		Ip:        ip,
	}
	sessionID, err := uc.sessionRepo.CreateSession(fiberCtx.Context(), 24*time.Hour*7, userData)
	if err != nil {
		return "", err
	}
	// fmt.Println("sessionID", sessionID)
	return sessionID, nil
}
