package database

import (
	"auth-service/internal/config"
	"auth-service/internal/domain"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	ErrDuplicateResource = errors.New("duplicate resource")
)

type Repository struct {
	db *sql.DB
}

func NewRepository(cfg *config.Config) (domain.AuthRepository, error) {
	db, err := newPostgresDB(cfg)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	repo := &Repository{db: db}

	return repo, nil
}

func (r *Repository) Close() error {
	if r.db != nil {
		err := r.db.Close()
		if err != nil {
			return fmt.Errorf("postgres connection close error: %w", err)
		}
		fmt.Println("📡 postgre connection closed successfully.")
	}
	return nil
}
