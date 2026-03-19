package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, id uuid.UUID, username, email string) error
	UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) error
	Close() error
}
