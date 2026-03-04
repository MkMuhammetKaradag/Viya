package usecase

import (
	"auth-service/internal/domain"
	"context"
	"fmt"
	"time"
)

type SignInUseCase interface {
	Execute(ctx context.Context, identifier, password string) error
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

func (uc *signInUsecase) Execute(ctx context.Context, identifier, password string) error {

	fmt.Println(
		"Executing SignInUseCase with identifier:", identifier,
		"and password:", password,
	)
	userData := &domain.SessionData{
		UserID: identifier, // In a real implementation, this would be the user's unique ID from the database
	}
	sessionID, err := uc.sessionRepo.CreateSession(ctx, 24*time.Hour, userData)
	if err != nil {
		return err
	}
	fmt.Println("sessionID", sessionID)
	return nil
}
