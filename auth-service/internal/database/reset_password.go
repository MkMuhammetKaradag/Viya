package database

import (
	"context"
	"database/sql"
	"fmt"
)

func (r *Repository) ResetPassword(ctx context.Context, token, newPassword string) error {
	// 1. Yeni şifreyi hashle
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("hashing error: %w", err)
	}

	// 2. Transaction başlat
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("transaction begin error: %w", err)
	}
	defer tx.Rollback()

	// 3. Token'ın kime ait olduğunu ve süresini kontrol et
	// (Tablonda expires_at olduğunu varsayıyorum)
	var userID string
	tokenQuery := `
        DELETE FROM forgot_password_tokens 
        WHERE token = $1 AND expires_at > NOW() 
        RETURNING user_id`

	err = tx.QueryRowContext(ctx, tokenQuery, token).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid or expired token")
		}
		return fmt.Errorf("token validation error: %w", err)
	}

	// 4. Kullanıcının şifresini güncelle
	// Ayrıca başarılı sıfırlamadan sonra failed_login_attempts'i de sıfırlamak iyi bir pratiktir.
	updateQuery := `
        UPDATE users 
        SET password = $1, 
            failed_login_attempts = 0, 
            account_locked = false, 
            lock_until = NULL 
        WHERE id = $2`

	_, err = tx.ExecContext(ctx, updateQuery, hashedPassword, userID)
	if err != nil {
		return fmt.Errorf("password update error: %w", err)
	}

	// 5. Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit error: %w", err)
	}

	return nil
}
