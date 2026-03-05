package usecase

import (
	"auth-service/internal/domain"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type ForgotPasswordUseCase interface {
	Execute(ctx context.Context, identifier string) error
}

type forgotPasswordUseCase struct {
	repo         domain.AuthRepository
	redisManager domain.SessionRepository
}

func NewForgotPasswordUseCase(repo domain.AuthRepository, redisManager domain.SessionRepository) ForgotPasswordUseCase {
	return &forgotPasswordUseCase{repo: repo, redisManager: redisManager}
}

func (u *forgotPasswordUseCase) Execute(ctx context.Context, identifier string) error {

	limitKey := fmt.Sprintf("pwd_limit:%s", identifier)

	// Redis'ten bu anahtarı kontrol et (Bu metodun sessionManager'da olduğunu varsayıyorum)
	isLocked, err := u.redisManager.IsActionLocked(ctx, limitKey)
	if err != nil {
		return err
	}
	if isLocked {
		// Güvenlik ipucu: Hata mesajını çok açık vermeyip "Too many requests" diyebilirsin
		return fmt.Errorf("too many requests, please wait a few minutes")
	}

	user, err := u.repo.GetUserByIdentifier(ctx, identifier)
	if err != nil {

		return fmt.Errorf("user not found: %w", err)
	}

	token, err := generateSecureToken()
	if err != nil {
		return err
	}
	fmt.Println("Generated token:", token)
	fmt.Println("User ID:", user.ID)
	err = u.repo.CreateForgotPasswordToken(ctx, user.ID, token)
	if err != nil {
		return err
	}
	_ = u.redisManager.SetActionLock(ctx, limitKey, 2*time.Minute)
	return nil
}
func generateSecureToken() (string, error) {
	b := make([]byte, 32) // 32 byte = 64 karakterli hex string
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
