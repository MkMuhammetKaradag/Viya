package database

import (
	"context"
	"social-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) GetPendingFollowRequests(ctx context.Context, userID uuid.UUID) ([]domain.PendingRequest, error) {
	// query := `
	//     SELECT f.follower_id, u.username, u.avatar_url, f.created_at
	//     FROM follows f
	//     JOIN users u ON f.follower_id = u.id
	//     WHERE f.following_id = $1 AND f.status = 'PENDING'
	//     ORDER BY f.created_at DESC
	// `

	query := `
        SELECT 
			f.follower_id, 
			u.username, 
			-- MANTIK: Eğer kullanıcı gizli DEĞİLSE veya ben onu TAKİP EDİYORSAM (ACCEPTED) avatarı gör, yoksa boş/varsayılan dön.
			CASE 
				WHEN u.is_private = false THEN u.avatar_url
				WHEN EXISTS (
					SELECT 1 FROM follows f2 
					WHERE f2.follower_id = $1 -- Ben
					AND f2.following_id = u.id -- O
					AND f2.status = 'ACCEPTED'
				) THEN u.avatar_url
				ELSE 'default_private_avatar.png' -- Gizli profil fotoğrafı
			END as avatar_url  ,
			f.created_at
		FROM follows f
		JOIN users u ON f.follower_id = u.id
		WHERE f.following_id = $1 AND f.status = 'PENDING'
		ORDER BY f.created_at DESC
    `
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []domain.PendingRequest
	for rows.Next() {
		var req domain.PendingRequest
		if err := rows.Scan(&req.UserID, &req.Username, &req.AvatarURL, &req.CreatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}
func (r *Repository) GetSentFollowRequests(ctx context.Context, userID uuid.UUID) ([]domain.PendingRequest, error) {
	query := `
        SELECT 
            f.following_id, 
            u.username, 
            -- MANTIK: İstek attığım kişi gizli değilse avatarı gör, gizliyse varsayılanı dön.
            -- (Zaten istek PENDING aşamasındaysa henüz kabul edilmemiş demektir)
            CASE 
                WHEN u.is_private = false THEN u.avatar_url
                ELSE 'default_private_avatar.png'
            END as avatar_url,
            f.created_at
        FROM follows f
        JOIN users u ON f.following_id = u.id
        WHERE f.follower_id = $1 AND f.status = 'PENDING'
        ORDER BY f.created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []domain.PendingRequest
	for rows.Next() {
		var req domain.PendingRequest
		if err := rows.Scan(&req.UserID, &req.Username, &req.AvatarURL, &req.CreatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *Repository) GetPendingFollowRequestsCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
        SELECT COUNT(*) 
        FROM follows 
        WHERE following_id = $1 AND status = 'PENDING'
    `
	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
