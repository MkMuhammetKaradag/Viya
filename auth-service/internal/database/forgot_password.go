package database

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) CreateForgotPasswordToken(ctx context.Context, userID uuid.UUID, token string) error {
	query := `INSERT INTO forgot_password_tokens (user_id, token) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, userID, token)
	return err
}
