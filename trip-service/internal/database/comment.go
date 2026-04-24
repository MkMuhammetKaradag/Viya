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

	updateQuery := `UPDATE trips SET total_comments = total_comments + 1 WHERE id = $1`

	_, err = tx.ExecContext(ctx, updateQuery, comment.TripID)
	if err != nil {
		return uuid.Nil, err
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return commentID, nil
}

func (r *Repository) GetTripComments(ctx context.Context, viewerID uuid.UUID, tripID uuid.UUID, page, limit int) ([]domain.Comment, error) {
	offset := (page - 1) * limit
	fmt.Println("offset:", offset, " limit:", limit)

	query := `
    SELECT 
        c.id, c.trip_id, c.user_id, u.username, u.avatar_url, c.content, c.created_at,
        (SELECT COUNT(*) FROM comments rc WHERE rc.parent_id = c.id AND rc.deleted_at IS NULL) as reply_count
    FROM comments c
    JOIN users u ON c.user_id = u.id
    LEFT JOIN local_blocks b ON (
        (b.blocker_id = $1 AND b.blocked_id = c.user_id) OR 
        (b.blocker_id = c.user_id AND b.blocked_id = $1)
    )
    WHERE c.trip_id = $2 
      AND c.parent_id IS NULL 
      AND c.deleted_at IS NULL
      AND b.blocker_id IS NULL
    ORDER BY c.created_at DESC
    LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, viewerID, tripID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var c domain.Comment
		err := rows.Scan(
			&c.ID, &c.TripID, &c.UserID, &c.Username, &c.AvatarURL,
			&c.Content, &c.CreatedAt, &c.ReplyCount,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}
func (r *Repository) GetCommentReplies(ctx context.Context, parentID uuid.UUID, page, limit int) ([]domain.Comment, error) {
	offset := (page - 1) * limit

	query := `
        SELECT c.id, c.trip_id, c.user_id, c.parent_id, u.username, u.avatar_url, c.content, c.created_at,
		(SELECT COUNT(*) FROM comments rc WHERE rc.parent_id = c.id AND rc.deleted_at IS NULL) as reply_count
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.parent_id = $1 AND c.deleted_at IS NULL
        ORDER BY c.created_at ASC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, parentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []domain.Comment
	for rows.Next() {
		var rc domain.Comment
		err := rows.Scan(&rc.ID, &rc.TripID, &rc.UserID, &rc.ParentID, &rc.Username, &rc.AvatarURL, &rc.Content, &rc.CreatedAt, &rc.ReplyCount)
		if err != nil {
			return nil, err
		}
		replies = append(replies, rc)
	}
	return replies, nil
}
