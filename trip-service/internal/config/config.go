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
	Server     ServerConfig     `mapstructure:"server"`
	DB         DBConfig         `mapstructure:"db"`
	Cloudinary CloudinaryConfig `mapstructure:"cloudinary"`
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

type CloudinaryConfig struct {
	CloudName string `mapstructure:"cloud_name"`
	APIKey    string `mapstructure:"api_key"`
	APISecret string `mapstructure:"api_secret"`
}

func Load() (*Config, error) {

	// 1️⃣ .env yükle (varsa)
	_ = godotenv.Load()

	v := viper.New()

	// 2️⃣ Config path
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("config path tespit edilemedi")
	}

	configDir := filepath.Dir(file)
	v.AddConfigPath(configDir)
	// 3️⃣ YAML dosyalarını merge et
	files := []string{"server.yaml", "database.yaml"}

	for _, file := range files {
		v.SetConfigFile(filepath.Join(configDir, file))
		if err := v.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("%s okunamadı: %w", file, err)
		}
	}

	// 4️⃣ ENV otomatik yükleme
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 5️⃣ ENV binding (override için)
	v.BindEnv("cloudinary.cloud_name", "CLOUD_NAME")
	v.BindEnv("cloudinary.api_key", "API_KEY")
	v.BindEnv("cloudinary.api_secret", "API_SECRET")

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config parse edilemedi: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	fmt.Println("cfg-load", cfg)

	return &cfg, nil
}

func (c *Config) Validate() error {

	if c.Server.Port == "" {
		return errors.New("server.port boş olamaz")
	}

	if c.DB.Host == "" {
		return errors.New("db.host boş olamaz")
	}

	if c.DB.Port == "" {
		return errors.New("db.port boş olamaz")
	}

	if c.DB.User == "" {
		return errors.New("db.user boş olamaz")
	}

	if c.DB.DBName == "" {
		return errors.New("db.dbname boş olamaz")
	}

	return nil
}

// func getCurrentConfigDir() string {
// 	_, file, _, ok := runtime.Caller(0)
// 	if !ok {
// 		panic("Config klasörü tespit edilemedi")
// 	}
// 	return filepath.Dir(file)
// }
