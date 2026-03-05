package domain

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID    `json:"id"`
	Username            string       `json:"username"`
	Email               string       `json:"email"`
	Password            string       `json:"-"`
	FailedLoginAttempts int          `json:"failed_login_attempts"`
	AccountLocked       bool         `json:"account_locked"`
	LockUntil           sql.NullTime `json:"lock_until,omitempty"`
}
