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
	IsPrivate   *bool
	Preferences []string
}
type UserRepository interface {
	CreateUser(ctx context.Context, id uuid.UUID, username, email string) error
	UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) error
	UpdateBanner(ctx context.Context, userID uuid.UUID, banerURL string) error
	UpdateProfile(ctx context.Context, userID uuid.UUID, params UpdateProfileParams) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)

	UpsertLocalBlock(ctx context.Context, blockerID, blockedID uuid.UUID) error
	UpsertLocalFollow(ctx context.Context, followerID, followingID uuid.UUID, status string) error

	Close() error
}
