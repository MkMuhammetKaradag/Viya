package database

import (
	"database/sql"
	"fmt"
)

func runMigrations(db *sql.DB) error {
	if _, err := db.Exec(usersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	// if _, err := db.Exec(updatedTrigger); err != nil {
	// 	return fmt.Errorf("trigger: %w", err)
	// }
	if _, err := db.Exec(localBlocksTable); err != nil {
		return fmt.Errorf("failed to create local_blocks table: %w", err)
	}
	if _, err := db.Exec(localFollowaTable); err != nil {
		return fmt.Errorf("failed to create local_follows table: %w", err)
	}

	return nil
}
