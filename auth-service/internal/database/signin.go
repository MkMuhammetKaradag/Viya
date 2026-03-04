package database

import (
	"auth-service/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username, email or password")
)

func (r *Repository) SignIn(ctx context.Context, identifier, password string) (*domain.User, error) {

	user := &domain.User{}
	query := `
		SELECT id,username,email,password
		FROM users
		WHERE (username = $1 OR email = $1)
	`
	err := r.db.QueryRowContext(ctx, query, identifier).Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("query error: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}
