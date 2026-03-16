package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	DB     DBConfig     `mapstructure:"db"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabitmq"`
}
type RabbitMQConfig struct {
	URL string `mapstructure:"url"`
}
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstucture:"dbname"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	v := viper.New()

	_, file, _, ok := runtime.Caller(0)

	if !ok {
		return nil, fmt.Errorf("config path  not found")
	}

	configDir := filepath.Dir(file)
	v.AddConfigPath(configDir)

	files := []string{"server.yaml", "database.yaml"}

	for _, file := range files {
		v.SetConfigFile(filepath.Join(configDir, file))
		if err := v.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("%s unreadable %w ", file, err)
		}
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "-"))

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config  could not parsed : %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return errors.New("server.Port cannot  be empty")
	}
	if c.DB.Host == "" {
		return errors.New("db.host cannot be empty")
	}

	if c.DB.Password == "" {
		return errors.New("db.port can not be empty")
	}
	if c.DB.DBName == "" {
		return errors.New("db.dbname  cannot be empty")
	}

	return nil
}
