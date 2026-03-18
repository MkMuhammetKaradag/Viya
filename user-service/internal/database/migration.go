package database

import (
	"database/sql"
	"fmt"
	"log"
)

func runMigrations(db *sql.DB) error {
	if _, err := db.Exec(usersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	// if _, err := db.Exec(updatedTrigger); err != nil {
	// 	return fmt.Errorf("trigger: %w", err)
	// }

	log.Println("Database migrated")
	return nil
}
