package database

const (
	createTripColon = `
        ALTER TABLE trips 
		ADD COLUMN IF NOT EXISTS content_vector vector(768), 
		
		ADD COLUMN IF NOT EXISTS location_name VARCHAR(255),
		ADD COLUMN IF NOT EXISTS total_likes INTEGER DEFAULT 0,
		ADD COLUMN IF NOT EXISTS total_comments INTEGER DEFAULT 0;
		ALTER TABLE waypoints
		ADD COLUMN IF NOT EXISTS category_id UUID REFERENCES categories(id) ON DELETE SET NULL;
		ALTER TABLE users
		ADD COLUMN IF NOT EXISTS interest_vector vector(768);
    `
	createExtension = `
        CREATE EXTENSION IF NOT EXISTS vector;
		CREATE EXTENSION IF NOT EXISTS pg_trgm;
		CREATE EXTENSION IF NOT EXISTS unaccent;
    `
	creatCategoryIndex = `
    CREATE INDEX IF NOT EXISTS idx_categories_name_trgm ON categories 
    USING gin (public.immutable_unaccent(name) gin_trgm_ops);
`
	tripsTable = `CREATE TABLE IF NOT EXISTS trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL, -- JWT'den gelecek
    
    -- İçerik
    title VARCHAR(255) NOT NULL,
    description TEXT,
    cover_image_url TEXT, -- Rota için Cloudinary kapak fotoğrafı
    
    -- İstatistik ve Durum
    is_active BOOLEAN DEFAULT true,
    is_public BOOLEAN DEFAULT true, 
    view_count INTEGER DEFAULT 0, -- Popülerlik takibi için
    
    -- Zamanlama
    start_date TIMESTAMP WITH TIME ZONE, -- Seyahat ne zaman başladı?
    end_date TIMESTAMP WITH TIME ZONE,

	published_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Otomatik Alanlar
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)`

	categoriesTable = `CREATE TABLE IF NOT EXISTS categories (
    -- UUID Standartı
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Üst kategori (Self-referencing)
    parent_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    
    -- İçerik
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    icon_url TEXT,
    description TEXT,
    
    -- AI/Vektör Altyapısı
    category_vector vector(768), 
    
    -- Zamanlama
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)`
	tripCategoriesTable = `CREATE TABLE IF NOT EXISTS trip_categories (
    trip_id UUID REFERENCES trips(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (trip_id, category_id)
)`
	waypointsTable = `CREATE TABLE IF NOT EXISTS waypoints (
		id UUID PRIMARY  KEY DEFAULT gen_random_uuid(),
		title VARCHAR(255) NOT NULL,
		description TEXT,
		order_index INT NOT NULL,
		trip_id UUID REFERENCES trips(id) ON DELETE CASCADE,
		latitude DOUBLE PRECISION NOT NULL,
		longitude DOUBLE PRECISION NOT NULL,
		note TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`

	photosTable = `CREATE TABLE IF NOT EXISTS photos (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		waypoint_id UUID REFERENCES waypoints(id) ON DELETE CASCADE,
		url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	photoTagsTable = `CREATE TABLE IF NOT EXISTS photo_tags (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		photo_id UUID REFERENCES photos(id) ON DELETE CASCADE,
		label TEXT NOT NULL,
		x_pos DOUBLE PRECISION NOT NULL,
		y_pos DOUBLE PRECISION NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	tripWiewsTable = `
		CREATE TABLE  IF NOT EXISTS  trip_views (
			trip_id UUID REFERENCES trips(id),
			user_id UUID REFERENCES users(id),
			viewed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (trip_id, user_id) -- Bu ikili UNIQUE olmalı ki aynı kişi 2. kez eklenemesin
		)
	`
	userasTable = `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY, -- Auth servisindeki ID ile aynı olmalı
		username VARCHAR(50) NOT NULL,
		email VARCHAR(100),
		avatar_url TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`

	tripWiewsConflict = `
		INSERT INTO trip_views (trip_id, user_id) 
		VALUES ($1, $2) 
		ON CONFLICT (trip_id, user_id) DO NOTHING;
	`
)
