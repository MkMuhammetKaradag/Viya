package database

import (
	"context"
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
