package app

import (
	"auth-service/internal/server"
	"time"
)

type App struct {
	// Add your application fields here
	server *server.Server
}

func NewApp() (*App, error) {
	c, err := buildContainer()
	if err != nil {
		return nil, err
	}
	return &App{server: c.server}, nil
}

type container struct {
	server *server.Server
}

func buildContainer() (*container, error) {

	srv := server.NewServer(
		getServerConfig(),
	)
	return &container{server: srv}, nil
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
