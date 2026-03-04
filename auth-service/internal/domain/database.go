package domain

import "context"

type AuthRepository interface {
	SignUp(ctx context.Context, username, email, password string) error
	SignIn(ctx context.Context, identifier, password string) (*User, error)
	Close() error
}
