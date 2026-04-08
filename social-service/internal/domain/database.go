package domain

import (
	"context"

	"github.com/google/uuid"
)

type SocialRepository interface {
	SaveUser(ctx context.Context, id uuid.UUID, username, email string) error
	UpdateUserSocialInfo(ctx context.Context, id uuid.UUID, isPrivate *bool, avatarURL *string) error

	CreateFollow(ctx context.Context, followerID, followingID uuid.UUID) (string, error)
	UpdateFollowStatus(ctx context.Context, followerID, followingID uuid.UUID, newStatus string) error
	GetPendingFollowRequests(ctx context.Context, userID uuid.UUID) ([]PendingRequest, error)
	UnfollowUser(ctx context.Context, followerID, targetUserID uuid.UUID) error
	RemoveFollower(ctx context.Context, myID, followerID uuid.UUID) error

	BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	IsBlocked(ctx context.Context, blockerID, blockedID uuid.UUID) ([]uuid.UUID, error)

	Close() error
}
