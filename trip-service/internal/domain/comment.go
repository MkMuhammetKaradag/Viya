package domain

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID         uuid.UUID  `json:"id"`
	TripID     uuid.UUID  `json:"trip_id"`
	UserID     uuid.UUID  `json:"user_id"`
	ParentID   *uuid.UUID `json:"parent_id"`
	Username   string     `json:"username"`
	AvatarURL  *string    `json:"avatar_url"`
	Content    string     `json:"content"`
	ReplyCount int        `json:"reply_count"`
	CreatedAt  time.Time  `json:"created_at"`
}
