package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (r *Repository) SignUp(ctx context.Context, username, email, password string) (uuid.UUID, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("hashing error: %w", err)
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, fmt.Errorf("transaction begin error: %w", err)
	}
	defer tx.Rollback()
	var userID uuid.UUID
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id

	`
	err = tx.QueryRowContext(ctx, query, username, email, hashedPassword).Scan(&userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert and scan error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, fmt.Errorf("transaction commit error: %w", err)
	}

	return userID, nil

}
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}
