package domain

import (
	"context"

	"github.com/google/uuid"
)

type UpdateProfileParams struct {
	FirstName   *string
	LastName    *string
	Bio         *string
	Location    *string
	Website     *string
	Preferences []string
}
type UserRepository interface {
	CreateUser(ctx context.Context, id uuid.UUID, username, email string) error
	UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) error
	UpdateProfile(ctx context.Context, userID uuid.UUID, params UpdateProfileParams) error
	Close() error
}
