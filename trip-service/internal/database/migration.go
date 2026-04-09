package database

import (
	"database/sql"
	"fmt"
	"log"
)

func runMigrations(db *sql.DB) error {

	if _, err := db.Exec(createExtension); err != nil {
		return fmt.Errorf("failed to create extension: %w", err)
	}
	if _, err := db.Exec(tripsTable); err != nil {
		return fmt.Errorf("failed to create trips table: %w", err)
	}
	if _, err := db.Exec(categoriesTable); err != nil {
		return fmt.Errorf("failed to create catogories table: %w", err)

	}
	if _, err := db.Exec(tripCategoriesTable); err != nil {
		return fmt.Errorf("failed to create tripCategory table : %w", err)
	}
	if _, err := db.Exec(createTripColon); err != nil {
		return fmt.Errorf("failed to create trip table colon : %w", err)
	}
	if _, err := db.Exec(waypointsTable); err != nil {
		return fmt.Errorf("failed to create waypoints table: %w", err)
	}
	if _, err := db.Exec(photosTable); err != nil {
		return fmt.Errorf("failed to create photos table: %w", err)
	}
	if _, err := db.Exec(userasTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	if _, err := db.Exec(tripWiewsTable); err != nil {
		return fmt.Errorf("failed to create trip_views table: %w", err)
	}
	if _, err := db.Exec(photoTagsTable); err != nil {
		return fmt.Errorf("failed to create photo_tags table: %w", err)
	}

	log.Println("Database migrated")
	return nil
}
