package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, id uuid.UUID, username, email string) error
	Close() error
}
