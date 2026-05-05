package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (r *Repository) CreateFollow(ctx context.Context, followerID, followingID uuid.UUID) (string, error) {
	var isPrivate bool

	err := r.db.QueryRowContext(ctx, "SELECT is_private FROM users WHERE id = $1", followingID).Scan(&isPrivate)
	if err != nil {
		return "", err
	}

	status := "ACCEPTED"
	if isPrivate {
		status = "PENDING"
	}

	query := `
        INSERT INTO follows (follower_id, following_id, status)
        VALUES ($1, $2, $3)
    `
	_, err = r.db.ExecContext(ctx, query, followerID, followingID, status)

	if err != nil {

		if pgErr, ok := err.(*pq.Error); ok {

			if pgErr.Code == "23505" {
				return "", errors.New("already_exists")
			}
		}
		return "", err
	}

	return status, nil
}
