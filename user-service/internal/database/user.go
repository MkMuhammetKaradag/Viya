package database

import (
	"context"
	"fmt"
	"user-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) SearchUsers(ctx context.Context, query string, currentUserID uuid.UUID, page int, limit int) ([]domain.UserSummary, error) {
	// Sorgu parametresini %...% formatına getiriyoruz
	searchQuery := "%" + query + "%"

	offset := (page - 1) * limit
	fmt.Println("search user query", searchQuery, "ofset", offset, "limit", limit, "currentUserID", currentUserID)

	sql := `
		SELECT 
			u.id, u.username, u.first_name, u.last_name, u.avatar_url, u.is_private,
			EXISTS (
				SELECT 1 FROM local_follows f 
				WHERE f.follower_id = $2 AND f.following_id = u.id AND f.status = 'following'
			) as is_following
		FROM users u
		WHERE (
			u.username ILIKE $1 OR 
			u.first_name ILIKE $1 OR 
			u.last_name ILIKE $1 OR
			(u.first_name || ' ' || u.last_name) ILIKE $1
		) 
		AND u.id != $2 
		AND u.deleted_at IS NULL
		AND NOT EXISTS (
			SELECT 1 FROM local_blocks b 
			WHERE (b.blocker_id = $2 AND b.blocked_id = u.id)
			   OR (b.blocker_id = u.id AND b.blocked_id = $2)
		)
		ORDER BY u.username ASC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, sql, searchQuery, currentUserID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.UserSummary
	for rows.Next() {
		var u domain.UserSummary
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.FirstName,
			&u.LastName,
			&u.AvatarURL,
			&u.IsPrivate,
			&u.IsFollowing,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
