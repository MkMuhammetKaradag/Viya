package database

const (
	tripsTable = `CREATE TABLE IF NOT EXISTS trips (
		id UUID PRIMARY  KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
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
		url TEXT NOT NULL
	)`
)
