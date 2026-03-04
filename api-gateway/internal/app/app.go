package app

import (
	"api-gateway/internal/config"
	"api-gateway/internal/graceful"
	"api-gateway/internal/server"
	"api-gateway/internal/session"
	"context"
	"fmt"
	"log"
	"time"
)

type App struct {
	cfg config.Config

	server         *server.Server
	sessionManager *session.SessionManager
}

func New(cfg *config.Config) *App {

	sessionManager, err := session.NewSessionManager(
		cfg.Redis.Addr,
		cfg.Redis.Password,
		cfg.Redis.DB,
		time.Hour,
	)
	if err != nil {
		log.Fatalf("❌ Redis Cache başlatılamadı: %v", err)
	}
	log.Println("✅ Redis Cache başarıyla bağlandı")

	server := server.New(cfg, sessionManager)
	return &App{

		server:         server,
		sessionManager: sessionManager,
	}
}
func (a *App) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Kafka consumer

	go graceful.WaitForShutdown(a.server.FiberApp(), 5*time.Second, ctx)

	log.Printf("starting user-service on %s", a.server.Address())
	if err := a.server.Start(); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	log.Println("server stopped, closing repository")
	return a.sessionManager.Close()
}
func (a *App) RegisterService(name string, baseURLs []string, prefix string) error {
	return a.server.Registry.Register(name, baseURLs, prefix)
}
