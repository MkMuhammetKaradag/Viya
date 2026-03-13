package usecase

import (
	"auth-service/internal/domain"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type ResetPasswordUseCase interface {
	Execute(ctx context.Context, newPassword, token, code, sessionID, platformType string) error
}

type resetPasswordUseCase struct {
	repo         domain.AuthRepository
	redisManager domain.SessionRepository
}

func NewResetPasswordUseCase(repo domain.AuthRepository, redisManager domain.SessionRepository) ResetPasswordUseCase {
	return &resetPasswordUseCase{
		repo:         repo,
		redisManager: redisManager,
	}
}

func (uc *resetPasswordUseCase) Execute(ctx context.Context, newPassword, token, code, sessionID, platformType string) error {

	if platformType == "mobile" {
		// --- MOBİL AKIŞI: SessionID (UUID) ve OTP (Code) ---
		if sessionID == "" || code == "" {
			return fmt.Errorf("session_id and code are required for mobile reset")
		}

		otpKey := fmt.Sprintf("otp_sess:%s", sessionID)

		// 1. Redis'ten değeri çek ("OTP:UserID")
		val, err := uc.redisManager.Get(ctx, otpKey)
		if err != nil {
			return fmt.Errorf("invalid or expired session")
		}

		// 2. Parçala ve Kontrol Et
		parts := strings.Split(val, ":")
		if len(parts) != 2 {
			return fmt.Errorf("internal session error")
		}

		savedOTP := parts[0]
		userIDStr := parts[1]

		if savedOTP != code {
			return fmt.Errorf("invalid verification code")
		}

		// 3. UserID üzerinden şifreyi güncelle
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return fmt.Errorf("invalid user identification")
		}

		err = uc.repo.UpdatePasswordByUserID(ctx, userID, newPassword)
		if err != nil {
			return err
		}

		// 4. Başarılıysa Redis'teki oturumu temizle
		_ = uc.redisManager.Delete(ctx, otpKey)

	} else {
		// --- WEB AKIŞI: Klasik Token ---
		if token == "" {
			return fmt.Errorf("token is required for web reset")
		}

		err := uc.repo.ResetPassword(ctx, token, newPassword)
		if err != nil {
			return err
		}
	}

	return nil
}
