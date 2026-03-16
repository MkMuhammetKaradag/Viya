package usecase

import (
	"auth-service/internal/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type SignupUseCase interface {
	Execute(ctx context.Context, username, email, password string) error
}
type signupUseCase struct {
	repo         domain.AuthRepository
	rabbitClient domain.RabbitMQClient
}

func NewSignupUseCase(repo domain.AuthRepository, rabbitClient domain.RabbitMQClient) SignupUseCase {
	return &signupUseCase{repo: repo, rabbitClient: rabbitClient}
}

type UserCreatedEvent struct {
	ID    uuid.UUID
	Email string
}

func (uc *signupUseCase) Execute(ctx context.Context, username, email, password string) error {
	fmt.Println(
		"Executing SignupUseCase with username:", username,
		"email:", email,
		"password:", password,
	)

	if err := uc.repo.SignUp(ctx, username, email, password); err != nil {
		return fmt.Errorf("signup error: %w", err)
	}
	event := UserCreatedEvent{ID: uuid.New(), Email: email}
	err := uc.rabbitClient.Publish("user_created", event)
	if err != nil {
		fmt.Println("rabbit err ", err)
	}
	return nil
}
