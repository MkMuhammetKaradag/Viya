package domain

import "context"

type AuthRepository interface {
	Signup(ctx context.Context, username, email, password string) error
	Close() error
}
