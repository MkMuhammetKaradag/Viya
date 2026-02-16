package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"trip-service/internal/config"
)

func newPostgresDB(cfg config.Config) (*sql.DB, error) {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DBName,
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	log.Println("Connected to PostgreSQL")
	return db, nil
}
