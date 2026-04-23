package database

import (
	"context"
	"fmt"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) CreateComment(ctx context.Context, comment *domain.Comment) (uuid.UUID, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()

	// 1. ŞART: Trip aktif mi ve sahibi kim?
	// Ayrıca gizli profili takip edip etmediğimizi kontrol etmek için joinliyoruz
	var tripOwnerID uuid.UUID
	var isActive bool
	var isPrivate bool
	var isFollowing bool

	checkQuery := `
		SELECT t.user_id, t.is_active, u.is_private,
		       EXISTS(SELECT 1 FROM local_follows f WHERE f.follower_id = $1 AND f.following_id = t.user_id AND f.status = 'ACCEPTED') as is_following
		FROM trips t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = $2 AND t.deleted_at IS NULL`

	err = tx.QueryRowContext(ctx, checkQuery, comment.UserID, comment.TripID).Scan(&tripOwnerID, &isActive, &isPrivate, &isFollowing)
	if err != nil {
		return uuid.Nil, fmt.Errorf("trip not found: %w", err)
	}

	// 2. KONTROLLER
	if !isActive {
		return uuid.Nil, fmt.Errorf("cannot comment on inactive trip")
	}

	// Eğer profil gizliyse ve yorum yazan kişi trip sahibi DEĞİLSE takibe bakılır
	if isPrivate && tripOwnerID != comment.UserID && !isFollowing {
		return uuid.Nil, fmt.Errorf("follow user to comment on their private trip")
	}

	// 3. EĞER CEVAPSA (ParentID varsa): Ana yorum aynı trip'e mi ait?
	if comment.ParentID != nil {
		var parentTripID uuid.UUID
		err = tx.QueryRowContext(ctx, "SELECT trip_id FROM comments WHERE id = $1", comment.ParentID).Scan(&parentTripID)
		if err != nil || parentTripID != comment.TripID {
			return uuid.Nil, fmt.Errorf("invalid parent comment")
		}
	}

	// 4. YORUMU EKLE
	insertQuery := `
		INSERT INTO comments (trip_id, user_id, parent_id, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	var commentID uuid.UUID
	err = tx.QueryRowContext(ctx, insertQuery, comment.TripID, comment.UserID, comment.ParentID, comment.Content).Scan(&commentID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert comment: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return commentID, nil
}
