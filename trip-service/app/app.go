package app

import (
	"time"
	"trip-service/config"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

type App struct {
	config config.Config
	// Add your application fields here
}

func NewApp(cfg config.Config) (*App, error) {

	return &App{config: cfg}, nil
}
func getServerConfig(cfg config.Config) AppConfig {
	return AppConfig{
		Port:         cfg.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

type AppConfig struct {
	Port         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (a *App) Run() error {
	cfg := getServerConfig(a.config)
	app := fiber.New(fiber.Config{
		AppName:      "Viya Trip Service v1.0",
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Concurrency:  256 * 1024,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"},
		AllowCredentials: true,
	}))
	app.Get("/health", func(c fiber.Ctx) error {
		// return c.SendString("OK")
		return c.JSON(fiber.Map{
			"status":  "active",
			"service": "trip-service",
		})
	})

	return app.Listen(cfg.Port)

}
