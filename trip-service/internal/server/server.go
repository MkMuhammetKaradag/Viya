// trip-service/internal/server/server.go
package server

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
)

type RouteRegistrar interface {
	Register(app *fiber.App)
}
type Config struct {
	Port         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
type Server struct {
	app *fiber.App
	cfg Config
}

func NewServer(cfg Config, registrar RouteRegistrar) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Concurrency:  256 * 1024,
	})

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "active",
			"service": "trip-service",
		})
	})

	if registrar != nil {
		registrar.Register(app)
	}
	return &Server{
		app: app,
		cfg: cfg,
	}
}

func (s *Server) Start() error {
	return s.app.Listen(s.Address())
}
func (s *Server) Address() string {
	return fmt.Sprintf("0.0.0.0:%s", s.cfg.Port)
}

func (s *Server) Shutdown(timeout time.Duration) error {
	return s.app.ShutdownWithTimeout(timeout)
}
func (s *Server) FiberApp() *fiber.App {
	return s.app
}
