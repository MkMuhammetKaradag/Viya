package database

import (
	"auth-service/internal/domain"
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const maxLoginAttempts = 5

var (
	ErrInvalidCredentials = errors.New("invalid username, email or password")
	ErrAccountLocked      = errors.New("account is locked, please try again later")
)

func (r *Repository) SignIn(ctx context.Context, identifier, password string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, username, email, password, failed_login_attempts, account_locked, lock_until FROM users WHERE username = $1 OR email = $1`

	err := r.db.QueryRowContext(ctx, query, identifier).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.FailedLoginAttempts, &user.AccountLocked, &user.LockUntil,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// 1. Kilit Kontrolü
	if user.AccountLocked && user.LockUntil.Valid && user.LockUntil.Time.After(time.Now()) {
		return nil, ErrAccountLocked
	}

	// 2. Şifre Doğrulama
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// SQL tarafında artık failed_login_attempts ve account_locked alanlarını güncelle
		updateQuery := `
            UPDATE users 
            SET failed_login_attempts = failed_login_attempts + 1,
                account_locked = CASE WHEN failed_login_attempts + 1 >= $1 THEN true ELSE false END,
                lock_until = CASE WHEN failed_login_attempts + 1 >= $1 THEN NOW() + INTERVAL '1 minute' ELSE NULL END
            WHERE id = $2`

		_, _ = r.db.ExecContext(ctx, updateQuery, maxLoginAttempts, user.ID)
		return nil, ErrInvalidCredentials
	}

	// 3. Başarılı Giriş (Sayaçları sıfırla)
	updateQuery := `UPDATE users SET failed_login_attempts = 0, account_locked = false, last_login = NOW(), lock_until = NULL WHERE id = $1`
	_, _ = r.db.ExecContext(ctx, updateQuery, user.ID)

	return user, nil
}
