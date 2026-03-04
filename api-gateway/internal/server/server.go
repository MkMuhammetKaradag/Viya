package server

import (
	"api-gateway/internal/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"api-gateway/internal/service"
	"api-gateway/internal/session"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type Server struct {
	cfg      *config.Config
	app      *fiber.App
	Registry *service.ServiceRegistry
}

func New(cfg *config.Config, sessionManager *session.SessionManager) *Server {
	registry := service.NewServiceRegistry()

	f := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	// Global Middleware
	f.Use(recover.New())
	f.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))
	f.Use(cors.New())
	protectedPrefixes := []string{
		"/api/v1/trips",
	}
	f.Use(middleware.AuthMiddleware(protectedPrefixes, sessionManager))

	proxyHandler := handlers.NewProxyHandler(registry, sessionManager)
	f.All("/*", proxyHandler.Handle)
	return &Server{
		cfg:      cfg,
		app:      f,
		Registry: registry,
	}
}
func (s *Server) Start() error {

	log.Printf("🌐 HTTP sunucusu %s adresinde dinliyor...", s.cfg.Server.Port)
	return s.app.Listen(s.Address())
}

func (s *Server) Shutdown(timeout time.Duration) error {

	return s.app.ShutdownWithTimeout(timeout)
}

func (s *Server) FiberApp() *fiber.App {
	return s.app
}

func (s *Server) Address() string {
	return fmt.Sprintf("0.0.0.0:%s", s.cfg.Server.Port)
}
