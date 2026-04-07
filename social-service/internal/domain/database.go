package domain

import (
	"context"

	"github.com/google/uuid"
)

type SocialRepository interface {
	SaveUser(ctx context.Context, id uuid.UUID, username, email string) error
	Close() error
}
