package database

import (
	"context"
	"encoding/json"
	"fmt"
	"user-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, username, first_name, last_name, bio, location, website, preferences, avatar_url,banner_url,is_private,created_at 
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL
	`
	var user domain.User
	var prefsJSON []byte // PostgreSQL'deki JSONB verisini tutmak için

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName,
		&user.Bio, &user.Location, &user.Website, &prefsJSON, &user.AvatarURL, &user.BannerURL, &user.IsPrivate, &user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Eğer preferences boş değilse, JSON'dan Go slice'ına çeviriyoruz
	if len(prefsJSON) > 0 {
		if err := json.Unmarshal(prefsJSON, &user.Preferences); err != nil {
			return nil, fmt.Errorf("failed to unmarshal preferences: %w", err)
		}
	}

	return &user, nil
}
