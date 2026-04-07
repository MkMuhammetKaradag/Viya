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
	if _, err := db.Exec(createTypeSQL); err != nil {
		return fmt.Errorf("failed to create follow_status type: %w", err)
	}
	if _, err := db.Exec(followsTable); err != nil {
		return fmt.Errorf("failed to create follows table: %w", err)
	}
	if _, err := db.Exec(blocksTable); err != nil {
		return fmt.Errorf("failed to create blocks table: %w", err)
	}

	log.Println("Database migrated")
	return nil
}
