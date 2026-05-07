package domain

import (
	"time"

	"github.com/google/uuid"
)

type PendingRequest struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}
