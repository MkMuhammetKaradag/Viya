package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Email       string    `json:"email" db:"email"`
	Username    string    `json:"username" db:"username"`
	FirstName   *string   `json:"first_name" db:"first_name"`
	LastName    *string   `json:"last_name" db:"last_name"`
	Bio         *string   `json:"bio" db:"bio"`
	Location    *string   `json:"location" db:"location"`
	Website     *string   `json:"website" db:"website"`
	AvatarURL   *string   `json:"avatar_url" db:"avatar_url"`
	BannerURL   *string   `json:"banner_url" db:"banner_url"`
	IsPrivate   *bool     `json:"is_private" db:"is_private"`
	Preferences []string  `json:"preferences" db:"preferences"` // JSONB için slice
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type UserSummary struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	FirstName   *string   `json:"first_name"`
	LastName    *string   `json:"last_name"`
	AvatarURL   *string   `json:"avatar_url"`
	BannerURL   *string   `json:"banner_url"`
	Website     *string   `json:"website"`
	Bio         *string   `json:"bio"`
	Location    *string   `json:"location"`
	IsPrivate   bool      `json:"is_private"`
	IsFollowing bool      `json:"is_following"`
	IsRequested bool      `json:"is_requested"`
}
