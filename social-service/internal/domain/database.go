package domain

import (
	"context"

	"github.com/google/uuid"
)

type SocialRepository interface {
	SaveUser(ctx context.Context, id uuid.UUID, username, email string) error
	UpdateUserSocialInfo(ctx context.Context, id uuid.UUID, isPrivate *bool, avatarURL *string) error

	CreatwFollow(ctx context.Context, followerID, followingID uuid.UUID) (string, error)

	Close() error
}
