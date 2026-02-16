package config

type ServerConfig struct {
	Port string
}

type Config struct {
	Server ServerConfig
	// Add your configuration fields here
}

func Read() (Config, error) {
	return Config{
		Server: ServerConfig{
			Port: ":8081",
		},
	}, nil
}
