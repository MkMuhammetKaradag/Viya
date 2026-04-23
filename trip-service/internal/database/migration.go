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

	// 	createFuncSQL := `
	// CREATE OR REPLACE FUNCTION public.immutable_unaccent(text)
	//   RETURNS text AS
	// $func$
	//   SELECT public.unaccent('public.unaccent', $1)
	// $func$  LANGUAGE sql IMMUTABLE;`

	// 	_, err := db.Exec(createFuncSQL)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to create immutable_unaccent function: %w", err)
	// 	}
	// 	createCategoryIndex := `
	//     CREATE INDEX IF NOT EXISTS idx_categories_name_trgm ON categories
	//     USING gin (public.immutable_unaccent(name) gin_trgm_ops);
	// `
	// 	_, err = db.Exec(createCategoryIndex)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to create trip_category index: %w", err)
	// 	}

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
	if _, err := db.Exec(localFollowaTable); err != nil {
		return fmt.Errorf("failed to create  local_follows  table: %w", err)
	}
	if _, err := db.Exec(localBlocksTable); err != nil {
		return fmt.Errorf("failed to create local_blocks table: %w", err)
	}
	if _, err := db.Exec(tripLikesTable); err != nil {
		return fmt.Errorf("failed to create trip_likes table: %w", err)
	}
	if _, err := db.Exec(commentsTable); err != nil {
		return fmt.Errorf("failed to create comments table: %w", err)
	}

	log.Println("Database migrated")
	return nil
}
