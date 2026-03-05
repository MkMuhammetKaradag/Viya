package database

import (
	"auth-service/internal/domain"
	"context"
)

func (r *Repository) GetUserByIdentifier(ctx context.Context, identifier string) (*domain.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE username = $1 OR email = $1`
	row := r.db.QueryRowContext(ctx, query, identifier)
	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
