package domain

import (
	"time"

	"github.com/google/uuid"
)

type PendingRequest struct {
	FollowerID uuid.UUID `json:"follower_id"`
	Username   string    `json:"username"`
	AvatarURL  *string   `json:"avatar_url"`
	CreatedAt  time.Time `json:"created_at"`
}
