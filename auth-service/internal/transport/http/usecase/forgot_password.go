package usecase

import (
	"auth-service/internal/domain"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ForgotPasswordUseCase interface {
	Execute(ctx context.Context, identifier, platformType string) (string, error)
}

type forgotPasswordUseCase struct {
	repo         domain.AuthRepository
	redisManager domain.SessionRepository
}

func NewForgotPasswordUseCase(repo domain.AuthRepository, redisManager domain.SessionRepository) ForgotPasswordUseCase {
	return &forgotPasswordUseCase{repo: repo, redisManager: redisManager}
}

func (u *forgotPasswordUseCase) Execute(ctx context.Context, identifier, platformType string) (string, error) {

	limitKey := fmt.Sprintf("pwd_limit:%s", identifier)

	// Redis'ten bu anahtarı kontrol et (Bu metodun sessionManager'da olduğunu varsayıyorum)
	isLocked, err := u.redisManager.IsActionLocked(ctx, limitKey)
	if err != nil {
		return "", err
	}
	if isLocked {
		// Güvenlik ipucu: Hata mesajını çok açık vermeyip "Too many requests" diyebilirsin
		return "", fmt.Errorf("too many requests, please wait a few minutes")
	}

	user, err := u.repo.GetUserByIdentifier(ctx, identifier)
	if err != nil {

		return "", fmt.Errorf("user not found: %w", err)
	}

	if platformType == "mobile" {
		// --- MOBİL AKIŞI: 6 Haneli OTP ---
		otpCode := generateOTP(6)
		sessionID := uuid.New().String() // Reset session ID

		// Redis'e kaydedilecek değer: "OTP:UserID"
		// Örn: "123456:550e8400-e29b-41d4-a716-446655440000"
		val := fmt.Sprintf("%s:%s", otpCode, user.ID.String())
		otpKey := fmt.Sprintf("otp_sess:%s", sessionID)

		// 10 dakika süreyle kaydet
		err = u.redisManager.Set(ctx, otpKey, val, 10*time.Minute)
		if err != nil {
			return "", err
		}

		fmt.Printf("MOBİL İÇİN OTP: %s (Session: %s)\n", otpCode, sessionID)
		return sessionID, nil // Mobilde state'de tutulacak olan UUID

	} else {
		// --- WEB AKIŞI: Uzun Token (Link) ---
		token, err := generateSecureToken()
		if err != nil {
			return "", err
		}

		// DB'ye kaydet (Senin mevcut metodun)
		err = u.repo.CreateForgotPasswordToken(ctx, user.ID, token)
		if err != nil {
			return "", err
		}

		// Web için linkli mail simülasyonu
		fmt.Printf("WEB LİNKİ GÖNDERİLDİ: https://viya.com/reset-password?token=%s\n", token)
	}
	_ = u.redisManager.SetActionLock(ctx, limitKey, 2*time.Minute)
	return "", nil
}
func generateSecureToken() (string, error) {
	b := make([]byte, 32) // 32 byte = 64 karakterli hex string
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
func generateOTP(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)
	rand.Read(result)
	for i := 0; i < length; i++ {
		result[i] = digits[result[i]%10]
	}
	return string(result)
}
