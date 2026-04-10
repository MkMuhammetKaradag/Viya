package database

import (
	"context"
	"trip-service/internal/domain"
)

func (r *Repository) SearchCategories(ctx context.Context, query string) ([]domain.Category, error) {

	sqlQuery := `
        SELECT id, name, description
        FROM categories 
        WHERE 
            unaccent(name) % unaccent($1) -- Benzerlik var mı?
            OR unaccent(description) % unaccent($1) -- Açıklamada benzerlik var mı?
            OR unaccent(name) ILIKE $2 OR unaccent(description) ILIKE $2 -- İçinde geçiyor mu?
        ORDER BY 
            unaccent(name) <-> unaccent($1) ASC, -- En benzer olan en üste (Mesafe en az olan)
            strict_word_similarity(unaccent($1), unaccent(name)) DESC
        LIMIT 10`

	searchTerm := "%" + query + "%"

	rows, err := r.db.QueryContext(ctx, sqlQuery, query, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}
