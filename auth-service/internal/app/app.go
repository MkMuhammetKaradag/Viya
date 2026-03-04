package app

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/domain"
	"auth-service/internal/server"
	"auth-service/internal/session"
	httptransport "auth-service/internal/transport/http"

	"fmt"
	"time"
)

type App struct {
	// Add your application fields here
	server  *server.Server
	repo    domain.AuthRepository
	session domain.SessionRepository
}

func NewApp(cfg *config.Config) (*App, error) {
	c, err := buildContainer(cfg)
	if err != nil {
		return nil, err
	}
	return &App{server: c.server, repo: c.repo, session: c.session}, nil
}

type container struct {
	server  *server.Server
	repo    domain.AuthRepository
	session domain.SessionRepository
}

func buildContainer(cfg *config.Config) (*container, error) {
	repo, sessionRepo, err := initStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}

	httpRouter := setupHttpRouter(cfg, repo, sessionRepo)

	return &container{server: server.NewServer(
		getServerConfig(),
		httpRouter,
	), repo: repo, session: sessionRepo}, nil
}

func initStorage(cfg *config.Config) (domain.AuthRepository, domain.SessionRepository, error) {
	repo, err := database.NewRepository(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("postgres init error: %w", err)
	}
	sessionRepo, err := session.NewSessionRepository(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("redis init error: %w", err)
	}

	return repo, sessionRepo, nil
}

func getServerConfig() server.ServerConfig {
	return server.ServerConfig{
		Port:         "8082",
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
}

func (a *App) Start() error {
	return a.server.Start()
}

func setupHttpRouter(cfg *config.Config, repo domain.AuthRepository, sessionRepo domain.SessionRepository) server.RouteRegistrar {

	handler := httptransport.NewHandlers(repo, sessionRepo)

	return httptransport.NewRouter(handler)
}
