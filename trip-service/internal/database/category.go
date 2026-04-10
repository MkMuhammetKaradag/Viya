package database

import (
	"context"
	"trip-service/internal/domain"
)

func (r *Repository) SearchCategories(ctx context.Context, query string) ([]domain.Category, error) {
	sqlQuery := `
        SELECT id, name
        FROM categories 
        WHERE name ILIKE $1 
        ORDER BY name ASC 
        LIMIT 10`

	rows, err := r.db.QueryContext(ctx, sqlQuery, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}
