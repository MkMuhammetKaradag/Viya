package database

const (
	usersTable = `CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        username VARCHAR(50) UNIQUE NOT NULL,
        first_name VARCHAR(50),
        last_name VARCHAR(50),
        bio TEXT,
        phone_number VARCHAR(15),
        avatar_url TEXT,
        banner_url TEXT,
		location VARCHAR(100),
		website TEXT,
		preferences JSONB DEFAULT '[]',
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
		is_private BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ -- Soft delete desteği için
    );`

	updatedTrigger = `
			CREATE OR REPLACE FUNCTION update_updated_at_column()
			RETURNS TRIGGER AS $$
			BEGIN
				NEW.updated_at = NOW();
				RETURN NEW;
			END;
			$$ language 'plpgsql';

			CREATE TRIGGER update_user_modtime
				BEFORE UPDATE ON users
				FOR EACH ROW
				EXECUTE PROCEDURE update_updated_at_column();
			`
)
