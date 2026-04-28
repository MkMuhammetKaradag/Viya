package database

import (
	"context"
	"database/sql"
	"fmt"
	"user-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) GetUser(ctx context.Context, currentUserID, userID uuid.UUID) (*domain.UserSummary, error) {
	query := `
        SELECT id, username, first_name, last_name, website, bio, location, avatar_url, banner_url, is_private,
        EXISTS (
            SELECT 1 FROM local_follows 
            WHERE follower_id = $1 AND following_id = users.id AND status = 'ACCEPTED'
        ) AS is_following
        FROM users
        WHERE id = $2 
          AND deleted_at IS NULL
          AND NOT EXISTS (
              SELECT 1 FROM local_blocks
              WHERE (blocker_id = $1 AND blocked_id = $2) 
                 OR (blocker_id = $2 AND blocked_id = $1)
          )
    `
	var user domain.UserSummary
	err := r.db.QueryRowContext(ctx, query, currentUserID, userID).Scan(
		&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Website,
		&user.Bio, &user.Location, &user.AvatarURL, &user.BannerURL,
		&user.IsPrivate, &user.IsFollowing,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}
