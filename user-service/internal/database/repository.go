package database

import (
	"database/sql"
	"user-service/internal/config"
	"user-service/internal/domain"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(cfg *config.Config) (domain.UserRepository, error) {
	db, err := newPostgresDB(cfg)

	if err != nil {
		return nil, err
	}
	repo := &Repository{
		db: db,
	}
	return repo, nil
}

func (r *Repository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
