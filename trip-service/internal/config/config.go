package config

type ServerConfig struct {
	Port string
}
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}
type Config struct {
	Server ServerConfig
	DB     DBConfig
	// Add your configuration fields here
}

func Read() (Config, error) {
	return Config{
		Server: ServerConfig{
			Port: ":8081",
		},
		DB: DBConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "user",
			Password: "password",
			DBName:   "tripdb",
		},
	}, nil
}
