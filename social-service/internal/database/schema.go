package database

const (
	usersTable = `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY, -- Auth servisindeki ID ile aynı olmalı
		username VARCHAR(50) NOT NULL,
		email VARCHAR(100),
		avatar_url TEXT,
		is_private BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`

	createTypeSQL = `
    DO $$ 
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'follow_status') THEN
            CREATE TYPE follow_status AS ENUM ('PENDING', 'ACCEPTED');
        END IF;
    END $$;`

	followsTable = `
	CREATE TABLE IF NOT EXISTS follows (
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status follow_status DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (follower_id, following_id),
    CONSTRAINT cannot_follow_self CHECK (follower_id <> following_id)
)
	`

	blocksTable = `
	CREATE TABLE IF NOT EXISTS blocks (
    blocker_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (blocker_id, blocked_id),
    CONSTRAINT cannot_block_self CHECK (blocker_id <> blocked_id)
)
	`
)
