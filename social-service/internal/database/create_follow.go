package database

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) CreateFollow(ctx context.Context, followerID, followingID uuid.UUID) (string, error) {
	var isPrivate bool

	// hedef kullanıcının gizli olup olmadığını kontrol et
	err := r.db.QueryRowContext(ctx, "SELECT is_private FROM users WHERE id = $1", followingID).Scan(&isPrivate)
	if err != nil {
		return "", err
	}

	status := "ACCEPTED"
	if isPrivate {
		status = "PENDING"
	}

	// 2. Takip ilişkisini ekle (Eğer engelli değilse - Block kontrolü servis katmanında yapılacak)
	query := `
        INSERT INTO follows (follower_id, following_id, status)
        VALUES ($1, $2, $3)
        ON CONFLICT (follower_id, following_id) DO NOTHING
    `
	_, err = r.db.ExecContext(ctx, query, followerID, followingID, status)
	return status, err
}
