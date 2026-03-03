package app

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/domain"
	"auth-service/internal/server"
	"fmt"
	"time"
)

type App struct {
	// Add your application fields here
	server *server.Server
	repo   domain.AuthRepository
}

func NewApp(cfg *config.Config) (*App, error) {
	c, err := buildContainer(cfg)
	if err != nil {
		return nil, err
	}
	return &App{server: c.server, repo: c.repo}, nil
}

type container struct {
	server *server.Server
	repo   domain.AuthRepository
}

func buildContainer(cfg *config.Config) (*container, error) {
	repo, err := initStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}
	srv := server.NewServer(
		getServerConfig(),
	)
	return &container{server: srv, repo: repo}, nil
}

func initStorage(cfg *config.Config) (domain.AuthRepository, error) {
	repo, err := database.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres init error: %w", err)
	}

	return repo, nil
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
