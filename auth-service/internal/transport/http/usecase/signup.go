package usecase

import (
	"auth-service/internal/domain"
	"context"
	"fmt"
	"log"
	"viya/pkg/messaging"

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

	userID, err := uc.repo.SignUp(ctx, username, email, password)
	if err != nil {
		return fmt.Errorf("signup error: %w", err)
	}
	userCreatedMessage := messaging.Message{
		Type:       messaging.AuthTypes.CreatedUser,
		ToServices: []messaging.ServiceType{messaging.UserService},
		Data: map[string]interface{}{
			"id":       userID,
			"email":    email,
			"username": username,
		},
		Critical: true,
	}

	err = uc.rabbitClient.PublishMessage(ctx, userCreatedMessage)
	if err != nil {
		log.Printf("User creation message could not be sent: %v", err)
		//return err
	}
	return nil
}
