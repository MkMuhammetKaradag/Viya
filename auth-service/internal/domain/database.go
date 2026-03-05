package domain

import (
	"context"

	"github.com/google/uuid"
)

type AuthRepository interface {
	SignUp(ctx context.Context, username, email, password string) error
	SignIn(ctx context.Context, identifier, password string) (*User, error)
	GetUserByIdentifier(ctx context.Context, identifier string) (*User, error)
	CreateForgotPasswordToken(ctx context.Context, userID uuid.UUID, token string) error
	Close() error
}
