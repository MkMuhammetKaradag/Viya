package database

const (
	usersTable = `CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	email VARCHAR(255) UNIQUE NOT NULL,
	username VARCHAR(50) NOT NULL UNIQUE,
	first_name VARCHAR(50),
	last_name VARCHAR(50),
	phone_number VARCHAR(15),
	avatar_url TEXT,
	banner_url TEXT,
	password VARCHAR(255) NOT NULL,
	failed_login_attempts INT DEFAULT 0,
	account_locked BOOLEAN DEFAULT false,

	lock_until TIMESTAMP WITH TIME ZONE,
	last_login TIMESTAMP WITH TIME ZONE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)`

	forgotPasswordTokensTable = `CREATE TABLE IF NOT EXISTS forgot_password_tokens (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL,
			token TEXT NOT NULL UNIQUE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`
)
