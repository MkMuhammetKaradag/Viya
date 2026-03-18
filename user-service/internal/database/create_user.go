package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) CreateUser(ctx context.Context, id uuid.UUID, username, email string) error {

	query := `
        INSERT INTO users (id, username, email, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, true, NOW(), NOW())
        ON CONFLICT (id) DO NOTHING; 
    `

	_, err := r.db.ExecContext(ctx, query, id, username, email)
	if err != nil {
		return fmt.Errorf("failed to insert user to database: %w", err)
	}

	return nil
}
