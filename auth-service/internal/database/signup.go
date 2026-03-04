package database

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (r *Repository) SignUp(ctx context.Context, username, email, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("hashing error: %w", err)
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("transaction begin error: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
	`
	_, err = tx.ExecContext(ctx, query, username, email, hashedPassword)
	if err != nil {
		return fmt.Errorf("insert error: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit error: %w", err)
	}
	return nil

}
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}
