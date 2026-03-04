package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Srver ServerConfig `mapstructure:"server"`
	Redis RedisConfig  `mapstructure:"redis"`
	DB    DBConfig     `mapstructure:"db"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func Load() (*Config, error) {

	_ = godotenv.Load()
	v := viper.New()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("config path tespit edilemedi")
	}

	configDir := filepath.Dir(file)

	v.AddConfigPath(configDir)

	files := []string{"server.yaml", "database.yaml"}

	for _, f := range files {
		v.SetConfigFile(filepath.Join(configDir, f))
		if err := v.MergeInConfig(); err != nil {
			return nil, err
		}
	}
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if c.Srver.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if c.DB.Host == "" || c.DB.Port == "" || c.DB.User == "" || c.DB.Password == "" || c.DB.DBName == "" {
		return fmt.Errorf("database configuration is incomplete")
	}
	return nil
}
