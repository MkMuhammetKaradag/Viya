package database

const (
	userasTable = `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY, -- Auth servisindeki ID ile aynı olmalı
		username VARCHAR(50) NOT NULL,
		email VARCHAR(100),
		avatar_url TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`
)
