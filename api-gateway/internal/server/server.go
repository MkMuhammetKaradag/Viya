package server

import (
	"api-gateway/internal/config"
	"api-gateway/internal/grpc_client"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"api-gateway/internal/service"
	"api-gateway/internal/session"
	"context"
	"fmt"
	"log"
	"net/http"
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
	f.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "active",
			"service": "api-service",
		})
	})
	f.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173", "*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		// Tarayıcı bazen ön kontrol yanıtlarının önbelleğe alınmasını ister
		MaxAge: 3600,
	}))
	protectedPrefixes := []string{
		"/api/v1/trips/me",
		"/api/v1/trips/user",
		"/api/v1/trips/liked",
		"/api/v1/trips/explore",
		"/api/v1/trips",
		"/api/v1/trips/:trip_id/like",
		"/api/v1/waypoints",
		"/api/v1/comments",
		"/api/v1/comments/trip/:trip_id",

		"/api/v1/waypoints",
		"/api/v1/auth/signout",
		"/api/v1/auth/all-signout",
		"/api/v1/users/upload-avatar",
		"/api/v1/users/upload-banner",
		"/api/v1/users/update-profile",
		"/api/v1/users/me",
		"/api/v1/users/search",
		"/api/v1/users/profile/:user_id",
		"/api/v1/users/profile",

		"/api/v1/social/follow",
		"/api/v1/social/unfollow",
		"/api/v1/social/remove-follower",
		"/api/v1/social/block",
		"/api/v1/social/unblock",
		"/api/v1/social/follow-request",
		"/api/v1/social/pending-requests",
		"/api/v1/social/sent-follow-requests",
		"/api/v1/social/pending-requests/count",
	}
	f.Use(middleware.AuthMiddleware(protectedPrefixes, sessionManager))

	proxyHandler := handlers.NewProxyHandler(registry, sessionManager)
	f.All("/*", proxyHandler.Handle)
	ctx := context.Background()
	registry.StartHealthChecks(ctx, 15*time.Second)
	return &Server{
		cfg:      cfg,
		app:      f,
		Registry: registry,
	}
}
func (s *Server) Start() error {
	go func() {
		if err := s.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("gRPC sunucusu hatası: %v", err)
		}
	}()
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
func (s *Server) Run() error {
	grpcAddress := "localhost:3001"

	if err := grpc_client.InitAuthClient(grpcAddress); err != nil {
		log.Fatalf("gRPC istemcisi başlatılamadı: %v", err)
		return err
	}
	return nil
}
