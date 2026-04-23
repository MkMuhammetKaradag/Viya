package domain

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID              uuid.UUID  `json:"id"`
	TripID          uuid.UUID  `json:"trip_id"`
	UserID          uuid.UUID  `json:"user_id"`
	ParentID        *uuid.UUID `json:"parent_id"` // Boş olabilir (ana yorumsa)
	Username        string     `json:"username"`
	AvatarURL       string     `json:"avatar_url"`
	Content         string     `json:"content"`
	ReplyToUsername string     `json:"reply_to_username,omitempty"` // Yanıtlanan kişinin adı
	CreatedAt       time.Time  `json:"created_at"`
}
