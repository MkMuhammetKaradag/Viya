package app

import (
	"fmt"
	"time"
	"user-service/internal/config"
	"user-service/internal/database"
	"user-service/internal/domain"
	"user-service/internal/graceful"
	"user-service/internal/server"
	httptransport "user-service/internal/transport/http"
)

type App struct {
	server   *server.Server
	config   *config.Config
	userRepo domain.UserRepository
}

func NewApp(cfg *config.Config) (*App, error) {
	c, err := buildContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("app failed:%w", err)
	}
	return &App{
		config:   cfg,
		userRepo: c.userRepo,
		server:   c.server,
	}, nil
}

type container struct {
	server   *server.Server
	userRepo domain.UserRepository
}

func buildContainer(cfg *config.Config) (*container, error) {
	repo, err := initStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("init progres repository:%w", err)
	}
	httpRouter := setupHttpRouter(cfg, repo)
	return &container{
		userRepo: repo,
		server: server.NewServer(
			getServerConfig(cfg),
			httpRouter,
		),
	}, nil
}

func initStorage(cfg *config.Config) (domain.UserRepository, error) {

	return database.NewRepository(cfg)
}

func setupHttpRouter(cfg *config.Config, userRepo domain.UserRepository) server.RouterRegister {
	handler := httptransport.NewHandlers(userRepo)
	return httptransport.NewRouter(handler)
}
func getServerConfig(cfg *config.Config) server.ServerConfig {
	return server.ServerConfig{
		Port:         cfg.Server.Port,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
}

func (a *App) Start() error {
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- a.server.Start()
	}()
	graceful.Shutdown(a.server.FiberApp(), 5*time.Second, a.userRepo)
	select {
	case err := <-serverErr:
		return err
	default:
		return nil
	}
}
