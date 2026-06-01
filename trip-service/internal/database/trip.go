package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"trip-service/internal/domain"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (r *Repository) GetTripByIDForAI(ctx context.Context, id uuid.UUID) (*domain.Trip, error) {
	var trip domain.Trip

	// 1. Ana Trip bilgilerini ve Kategori İsimlerini çekiyoruz
	// trip_categories üzerinden categories tablosuna join yapıyoruz
	query := `
        SELECT t.id, t.title, t.description, t.location_name,
               ARRAY_AGG(c.name) FILTER (WHERE c.name IS NOT NULL) as category_names
        FROM trips t
        LEFT JOIN trip_categories tc ON t.id = tc.trip_id
        LEFT JOIN categories c ON tc.category_id = c.id
        WHERE t.id = $1
        GROUP BY t.id`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&trip.ID, &trip.Title, &trip.Description, &trip.LocationName,
		pq.Array(&trip.CategoryNames), // pq kütüphanesi ARRAY_AGG için lazım
	)
	if err != nil {
		fmt.Println("GetTripByIDForAI err:", err)
		return nil, err
	}

	// 2. Waypoint'leri de çekelim ki prompt zenginleşsin
	wpQuery := `SELECT title, description FROM waypoints WHERE trip_id = $1 ORDER BY order_index`
	rows, err := r.db.QueryContext(ctx, wpQuery, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var wp domain.Waypoint
			if err := rows.Scan(&wp.Title, &wp.Description); err == nil {
				trip.Waypoints = append(trip.Waypoints, wp)
			}
		}
	}

	return &trip, nil
}
func (r *Repository) UpdateTripEmbedding(ctx context.Context, id uuid.UUID, vector []float32) error {

	fmt.Println("veri tabanına geldi")
	// pgvector formatına çeviriyoruz: [0.1, 0.2, ...]
	strVector := "["
	for i, v := range vector {
		strVector += fmt.Sprintf("%f", v)
		if i < len(vector)-1 {
			strVector += ","
		}
	}
	strVector += "]"

	query := `UPDATE trips SET content_vector = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, strVector, id)
	return err
}
func (r *Repository) GetExploreTrips(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.TripExploreDTO, error) {
	fmt.Println(userID)
	query := `
	WITH user_context AS (
		SELECT interest_vector FROM users WHERE id = $1
	)
	SELECT 
		t.id, t.user_id, t.title, t.cover_image_url, t.total_likes, t.total_comments, t.view_count, t.published_at,
		u.username as owner_username,
		u.avatar_url as owner_avatar,
		-- 📍 Toplam waypoint sayısını alt sorgu ile alıyoruz
		(SELECT COUNT(*) FROM waypoints WHERE trip_id = t.id) as waypoint_count,
		-- 📸 Eğer kapak fotoğrafı yoksa, ilk waypoint'in ilk fotoğrafını getir
		COALESCE(t.cover_image_url, (
			SELECT p.url FROM photos p 
			JOIN waypoints w ON p.waypoint_id = w.id 
			WHERE w.trip_id = t.id 
			ORDER BY w.order_index ASC, p.created_at ASC 
			LIMIT 1
		)) as display_image,
		-- 🧠 Keşfet Puanı
		COALESCE(
			(
				((1 - (t.content_vector <=> (SELECT interest_vector FROM user_context))) * 5.0) +
				(LOG(t.total_likes + 1) * 0.8 + LOG(t.total_comments + 1) * 1.2 + LOG(t.view_count + 1) * 0.3) +
				(exp(-0.01 * EXTRACT(DAY FROM (now() - t.published_at))))
			), 
			0.0 -- Eğer yukarıdaki hesaplama NULL dönerse skoru 0 yap
		) as explore_score
	FROM trips t
	JOIN users u ON t.user_id = u.id
	WHERE t.is_active = true AND t.is_public = true
	  -- 🛡️ Blok ve Gizlilik Duvarı (Aynı kalıyor)
	  AND NOT EXISTS (
		  SELECT 1 FROM local_blocks 
		  WHERE (blocker_id = $1 AND blocked_id = t.user_id) OR (blocker_id = t.user_id AND blocked_id = $1)
	  )
	  AND (u.is_private = false OR EXISTS (
		  SELECT 1 FROM local_follows WHERE follower_id = $1 AND following_id = t.user_id AND status = 'ACCEPTED'
	  ))
	ORDER BY explore_score DESC
	LIMIT $2 OFFSET $3;
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("lightweight explore query failed: %w", err)
	}
	defer rows.Close()

	var trips []domain.TripExploreDTO
	for rows.Next() {
		var t domain.TripExploreDTO
		var score float64 // Puanı okuyoruz ama dışarı vermiyoruz

		err := rows.Scan(
			&t.ID, &t.UserID, &t.Title, &t.DisplayImage, &t.TotalLikes, &t.TotalComments, &t.ViewCount, &t.PublishedAt,
			&t.OwnerUsername, &t.OwnerAvatar, &t.WaypointCount, &t.DisplayImage, &score,
		)
		if err != nil {
			return nil, fmt.Errorf("scan explore trip failed: %w", err)
		}
		trips = append(trips, t)
	}

	return trips, nil
}
func (r *Repository) GetHomeFeedTrips(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.TripExploreDTO, error) {
	query := `
    SELECT 
        t.id, t.user_id, t.title, t.cover_image_url, t.total_likes, t.total_comments, t.view_count, t.published_at,
        u.username as owner_username,
        u.avatar_url as owner_avatar,
        (SELECT COUNT(*) FROM waypoints WHERE trip_id = t.id) as waypoint_count,
        COALESCE(t.cover_image_url, (
            SELECT p.url FROM photos p 
            JOIN waypoints w ON p.waypoint_id = w.id 
            WHERE w.trip_id = t.id 
            ORDER BY w.order_index ASC, p.created_at ASC 
            LIMIT 1
        )) as display_image
    FROM trips t
    JOIN users u ON t.user_id = u.id
    WHERE t.is_active = true 
      AND t.is_public = true
      -- 🚀  Sadece ben veya takip ettiğim kişiler
      AND (
          t.user_id = $1  -- Kendi gezilerim
          OR t.user_id IN (
              SELECT following_id FROM local_follows 
              WHERE follower_id = $1 AND status = 'ACCEPTED'
          )
      )
      -- Blok kontrolü (Güvenlik için her zaman olmalı)
      AND NOT EXISTS (
          SELECT 1 FROM local_blocks 
          WHERE (blocker_id = $1 AND blocked_id = t.user_id) OR (blocker_id = t.user_id AND blocked_id = $1)
      )
    ORDER BY t.published_at DESC -- Ana sayfada genellikle en yeni olan en üsttedir
    LIMIT $2 OFFSET $3;
    `

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("home feed query failed: %w", err)
	}
	defer rows.Close()

	var trips []domain.TripExploreDTO
	for rows.Next() {
		var t domain.TripExploreDTO
		err := rows.Scan(
			&t.ID, &t.UserID, &t.Title, &t.DisplayImage, &t.TotalLikes, &t.TotalComments, &t.ViewCount, &t.PublishedAt,
			&t.OwnerUsername, &t.OwnerAvatar, &t.WaypointCount, &t.DisplayImage,
		)
		if err != nil {
			return nil, fmt.Errorf("scan home feed trip failed: %w", err)
		}
		trips = append(trips, t)
	}

	return trips, nil
}
func (r *Repository) ToggleTripLike(ctx context.Context, tripID, userID uuid.UUID) (bool, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	var isLiked bool
	query := `
        WITH deleted AS (
            DELETE FROM trip_likes 
            WHERE trip_id = $1 AND user_id = $2 
            RETURNING *
        ),
        inserted AS (
            INSERT INTO trip_likes (trip_id, user_id)
            SELECT $1, $2 
            WHERE NOT EXISTS (SELECT 1 FROM deleted)
            RETURNING true
        )
        SELECT COALESCE(
            (SELECT true FROM inserted), 
            (SELECT false FROM deleted)
        )`

	err = tx.QueryRowContext(ctx, query, tripID, userID).Scan(&isLiked)
	if err != nil {
		return false, fmt.Errorf("failed to toggle like: %w", err)
	}

	//  Sayacı güncelle
	var updateQuery string
	if isLiked {
		updateQuery = `UPDATE trips SET total_likes = total_likes + 1 WHERE id = $1`
	} else {
		updateQuery = `UPDATE trips SET total_likes = total_likes - 1 WHERE id = $1`
	}

	_, err = tx.ExecContext(ctx, updateQuery, tripID)
	if err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return isLiked, nil
}
func (r *Repository) GetTripStatus(ctx context.Context, tripID uuid.UUID) (*domain.TripStatusDTO, error) {
	query := `
        SELECT 
            t.user_id, 
            t.is_public, 
            u.is_private as owner_is_private
        FROM trips t
        JOIN users u ON t.user_id = u.id
        WHERE t.id = $1`

	var status domain.TripStatusDTO
	err := r.db.QueryRowContext(ctx, query, tripID).Scan(
		&status.UserID,
		&status.IsPublic,
		&status.OwnerIsPrivate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("gezi bulunamadı")
		}
		return nil, fmt.Errorf("gezi durumu sorgulanırken hata: %w", err)
	}

	return &status, nil
}

func (r *Repository) ForkTrip(ctx context.Context, originalTripID uuid.UUID, forkUserID uuid.UUID) (uuid.UUID, error) {
	// 1. TRANSACTION BAŞLATALIM
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	// Hata durumunda rollback yapması için güvenceye alalım
	defer tx.Rollback()

	// 2. GÜVENLİK VE GİZLİLİK KONTROLÜ
	// (Gezi aktif mi, public mi, forkable mı, engel var mı, gizli hesap mı kontrolü)
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 
			FROM trips t
			JOIN users u ON t.user_id = u.id
			WHERE t.id = $1 AND t.is_active = true AND t.is_public = true AND t.is_forkable = true
			AND NOT EXISTS (
				SELECT 1 FROM local_blocks 
				WHERE (blocker_id = $2 AND blocked_id = t.user_id) OR (blocker_id = t.user_id AND blocked_id = $2)
			)
			AND (
				u.is_private = false OR t.user_id = $2
				OR EXISTS (
					SELECT 1 FROM local_follows 
					WHERE follower_id = $2 AND following_id = t.user_id AND status = 'ACCEPTED'
				)
			)
		);`

	var isAllowed bool
	err = tx.QueryRowContext(ctx, checkQuery, originalTripID, forkUserID).Scan(&isAllowed)
	if err != nil {
		return uuid.Nil, fmt.Errorf("privacy check failed: %w", err)
	}
	if !isAllowed {
		return uuid.Nil, fmt.Errorf("action not allowed: trip is not forkable or privacy constraints violated")
	}

	// 3. ORİJİNAL GEZİ BİLGİLERİNİ ALALIM
	getTripQuery := `
		SELECT title, description, cover_image_url, location_name
		FROM trips WHERE id = $1`

	var title, description, locationName sql.NullString
	var coverImageURL sql.NullString
	err = tx.QueryRowContext(ctx, getTripQuery, originalTripID).Scan(&title, &description, &coverImageURL, &locationName)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to fetch original trip: %w", err)
	}

	// Forklanan gezinin başlığına bir işaret koyalım (Örn: "Ege Turu (Forked)")
	forkedTitle := fmt.Sprintf("%s (Forked)", title.String)

	// 4. YENİ GEZİ KAYDINI OLUŞTURMA (trips tablosuna insert)
	newTripID := uuid.New()
	insertTripQuery := `
		INSERT INTO trips (id, user_id, parent_id, title, description, cover_image_url, location_name, is_forkable, is_public, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, true, false, true);` // Çatallanan geziyi varsayılan olarak gizli (is_public=false) başlatalım, kullanıcı isterse yayınlar

	_, err = tx.ExecContext(ctx, insertTripQuery,
		newTripID,
		forkUserID,
		originalTripID, // parent_id bağlantısı
		forkedTitle,
		description,
		coverImageURL,
		locationName,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert forked trip: %w", err)
	}

	// 5. WAYPOINT'LERİ KOPYALAMA (Fotoğraflar hariç!)
	// Önce orijinal durakları çekiyoruz
	getWaypointsQuery := `
		SELECT title, description, order_index, latitude, longitude, note, category_id
		FROM waypoints WHERE trip_id = $1 ORDER BY order_index ASC`

	rows, err := tx.QueryContext(ctx, getWaypointsQuery, originalTripID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to fetch original waypoints: %w", err)
	}

	// Geçici bir struct tanımlayarak verileri hafızaya alalım
	type tempWaypoint struct {
		title, description, note sql.NullString
		orderIndex               int
		lat, lng                 float64
		categoryID               uuid.NullUUID
	}
	var waypointList []tempWaypoint

	// Okuma döngüsü
	for rows.Next() {
		var w tempWaypoint
		err := rows.Scan(&w.title, &w.description, &w.orderIndex, &w.lat, &w.lng, &w.note, &w.categoryID)
		if err != nil {
			rows.Close() // Hata durumunda bağlantıyı hemen kapat
			return uuid.Nil, fmt.Errorf("failed to scan waypoint: %w", err)
		}
		waypointList = append(waypointList, w)
	}
	rows.Close() // 🚀 KRİTİK NOKTA: Okuma bitti, bağlantıyı INSERT'ler için serbest bırakıyoruz!

	// Şimdi temiz bağlantı üzerinden yeni durakları güvenle insert edebiliriz
	insertWaypointQuery := `
		INSERT INTO waypoints (id, title, description, order_index, trip_id, latitude, longitude, note, category_id)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8)`

	for _, w := range waypointList {
		_, err = tx.ExecContext(ctx, insertWaypointQuery,
			w.title, w.description, w.orderIndex, newTripID, w.lat, w.lng, w.note, w.categoryID,
		)
		if err != nil {
			return uuid.Nil, fmt.Errorf("failed to clone waypoint: %w", err)
		}
	}

	// 6. HER ŞEY YOLUNDAYSA COMMIT
	if err = tx.Commit(); err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	return newTripID, nil
}
func (r *Repository) GetTripByID(ctx context.Context, id uuid.UUID) (*domain.Trip, error) {
	// Soft delete (deleted_at IS NULL) kontrolünü sorguya ekledik usta!
	query := `
		SELECT 
			id, user_id, title, description, cover_image_url, 
			parent_id, is_forkable, is_active, is_public, 
			view_count, start_date, end_date, published_at, 
			created_at, updated_at
		FROM trips 
		WHERE id = $1 AND deleted_at IS NULL
	`

	trip := &domain.Trip{}

	// Nullable olabilecek alanlar için sql.NullX yapıları veya domain modelinde pointer kullanmalısın.
	// Domain modelinde alanların pointer (*string, *time.Time) olduğunu varsayarak scan ediyoruz:
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&trip.ID,
		&trip.UserID,
		&trip.Title,
		&trip.Description,   // domain.Trip içinde string veya pointer string
		&trip.CoverImageURL, // *string
		&trip.ParentID,      // *uuid.UUID
		&trip.IsForkable,
		&trip.IsActive,
		&trip.IsPublic,
		&trip.ViewCount,
		&trip.StartDate, // *time.Time
		&trip.EndDate,   // *time.Time
		&trip.PublishedAt,
		&trip.CreatedAt,
		&trip.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("trip not found")
		}
		return nil, err
	}

	return trip, nil
}

func (r *Repository) UpdateTrip(ctx context.Context, trip *domain.Trip) error {
	// updated_at = NOW() ekleyerek tablonun güncelliğini koruyoruz usta.
	query := `
		UPDATE trips 
		SET 
			title = $1, 
			description = $2, 
			cover_image_url = $3, 
			is_forkable = $4, 
			is_active = $5, 
			is_public = $6, 
			start_date = $7, 
			end_date = $8, 
			published_at = $9,
			updated_at = NOW()
		WHERE id = $10 AND user_id = $11 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		trip.Title,
		trip.Description,
		trip.CoverImageURL,
		trip.IsForkable,
		trip.IsActive,
		trip.IsPublic,
		trip.StartDate,
		trip.EndDate,
		trip.PublishedAt,
		trip.ID,
		trip.UserID, // 🛡️ Güvenlik için sorguda da user_id kontrolü yapıyoruz usta
	)
	if err != nil {
		return err
	}

	// Etkilenen satır sayısını kontrol edelim (Gezinin silinmediğinden emin olmak için)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no trip was updated (either not found or unauthorized)")
	}

	return nil
}
