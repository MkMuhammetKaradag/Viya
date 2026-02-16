package database

import (
	"database/sql"
	"errors"
	"trip-service/internal/config"
	"trip-service/internal/domain"

	_ "github.com/lib/pq"
)

var (
	ErrDuplicateResource = errors.New("duplicate resource")
)

type Repository struct {
	db *sql.DB
}

func NewRepository(cfg config.Config) (domain.TripRepository, error) {
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
		return r.db.Close()
	}
	return nil
}
