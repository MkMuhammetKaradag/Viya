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
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)`
)
