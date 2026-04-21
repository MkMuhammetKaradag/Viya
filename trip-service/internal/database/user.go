package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"trip-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) UpdateUserInterest(ctx context.Context, userID uuid.UUID, tripID uuid.UUID, weight float32) error {
	var oldInterestStr, tripVecStr sql.NullString

	// 1. Verileri string (metin) olarak çekiyoruz (Hata almamak için)
	err := r.db.QueryRowContext(ctx, "SELECT interest_vector::text FROM users WHERE id = $1", userID).Scan(&oldInterestStr)
	err = r.db.QueryRowContext(ctx, "SELECT content_vector::text FROM trips WHERE id = $1", tripID).Scan(&tripVecStr)

	if err != nil || !tripVecStr.Valid {
		return fmt.Errorf("trip vector not found or null: %w", err)
	}

	// 2. Bu stringleri (ör: "[0.1,0.2]") Go'nun anlayacağı []float32'ye çeviriyoruz
	oldInterestVec := parseVectorString(oldInterestStr.String)
	tripVec := parseVectorString(tripVecStr.String)

	// 3. Go tarafında ağırlıklı ortalamayı hesapla
	newInterestVec := calculateWeightedAverage(oldInterestVec, tripVec, weight)

	// 4. Tekrar stringe çevir ve DB'ye yaz
	strForDB := formatVector(newInterestVec)
	_, err = r.db.ExecContext(ctx, "UPDATE users SET interest_vector = $1 WHERE id = $2", strForDB, userID)

	return err
}
func parseVectorString(str string) []float32 {
	if str == "" || str == "NULL" {
		return nil
	}
	// Köşeli parantezleri temizle: [0.1,0.2] -> 0.1,0.2
	str = strings.Trim(str, "[]")
	if str == "" {
		return nil
	}

	parts := strings.Split(str, ",")
	vec := make([]float32, len(parts))
	for i, p := range parts {
		val, _ := strconv.ParseFloat(strings.TrimSpace(p), 32)
		vec[i] = float32(val)
	}
	return vec
}
func calculateWeightedAverage(oldVec, newVec []float32, weight float32) []float32 {
	if len(oldVec) == 0 {
		return newVec
	}

	result := make([]float32, len(oldVec))
	for i := range oldVec {
		// Formül: (Eski * 0.95) + (Yeni * 0.05)
		result[i] = (oldVec[i] * (1.0 - weight)) + (newVec[i] * weight)
	}
	return result
}
func formatVector(v []float32) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, val := range v {
		// Sayıları metne çeviriyoruz (f formatı: ondalıklı sayı)
		sb.WriteString(fmt.Sprintf("%f", val))
		if i < len(v)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func (r *Repository) GetLikedTrips(ctx context.Context, userID uuid.UUID, limit, page int) ([]domain.TripSummary, error) {

	offset := (page - 1) * limit
	query := `
        SELECT 
            t.id, 
            t.title, 
            COALESCE(
                t.cover_image_url, 
                (SELECT p.url FROM photos p 
                 JOIN waypoints w ON p.waypoint_id = w.id 
                 WHERE w.trip_id = t.id 
                 ORDER BY w.order_index ASC, p.created_at ASC 
                 LIMIT 1)
            ) as effective_cover_url,
            t.is_public, 
            t.view_count, 
			t.total_likes,
            t.created_at,
            (SELECT COUNT(*) FROM waypoints WHERE trip_id = t.id) as waypoint_count
        FROM trips t
        INNER JOIN trip_likes tl ON t.id = tl.trip_id
        WHERE tl.user_id = $1
        ORDER BY tl.created_at DESC
        LIMIT $2 OFFSET $3 
    `

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	count := 0
	var summaries []domain.TripSummary
	for rows.Next() {
		count++
		var s domain.TripSummary
		err := rows.Scan(
			&s.ID, &s.Title, &s.CoverImageURL, &s.IsPublic, &s.ViewCount, &s.LikeCount, &s.CreatedAt, &s.WaypointCount,
		)
		if err != nil {
			fmt.Printf("SCAN HATASI: %v\n", err)
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, nil
}
