package app

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/domain"
	"auth-service/internal/graceful"
	"auth-service/internal/server"
	"auth-service/internal/session"
	grpctransport "auth-service/internal/transport/grpc"
	httptransport "auth-service/internal/transport/http"
	"fmt"
	"log"
	"time"
	"viya/pkg/messaging"
)

type App struct {
	// Add your application fields here
	server  *server.Server
	repo    domain.AuthRepository
	session domain.SessionRepository
	rabbit  domain.RabbitMQClient
}

func NewApp(cfg *config.Config) (*App, error) {
	c, err := buildContainer(cfg)
	if err != nil {
		return nil, err
	}
	return &App{server: c.server, repo: c.repo, session: c.session, rabbit: c.rabbit}, nil
}

type container struct {
	server  *server.Server
	repo    domain.AuthRepository
	session domain.SessionRepository
	rabbit  domain.RabbitMQClient
}

func buildContainer(cfg *config.Config) (*container, error) {
	repo, sessionRepo, _, err := initStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}

	rabbitClient, err := initMessaging()
	if err != nil {
		return nil, fmt.Errorf("init rabbit :%w", err)
	}

	httpRouter := setupHttpRouter(cfg, repo, sessionRepo, rabbitClient)
	grpcHandler := grpctransport.NewAuthGrpcHandler(sessionRepo)
	return &container{server: server.NewServer(
		getServerConfig(cfg),
		httpRouter,
		grpcHandler,
	), repo: repo, session: sessionRepo, rabbit: rabbitClient}, nil
}

func initStorage(cfg *config.Config) (domain.AuthRepository, domain.SessionRepository, domain.RabbitMQClient, error) {
	repo, err := database.NewRepository(cfg)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("postgres init error: %w", err)
	}
	sessionRepo, err := session.NewSessionRepository(cfg)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("redis init error: %w", err)
	}

	// rabbitClient, err := messaging.NewRabbitClient(cfg.RabbitMQ.URL)
	// if err != nil {
	// 	return nil, nil, nil, fmt.Errorf("rabbitmq init error: %w", err)
	// }

	return repo, sessionRepo, nil, nil
}

func initMessaging() (domain.RabbitMQClient, error) {
	config := messaging.NewDefaultConfig("")
	config.RetryTypes = []messaging.MessageType{}
	rabbitMQ, err := messaging.NewRabbitClient(config, messaging.AuthService)
	if err != nil {
		log.Fatalf("RabbitMQ bağlantısı kurulamadı: %v", err)
		return nil, err
	}
	return rabbitMQ, nil

}

func getServerConfig(cfg *config.Config) server.ServerConfig {
	return server.ServerConfig{
		GrpcPort:     cfg.Srver.GrpcPort,
		Port:         cfg.Srver.Port,
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

	// 2. Shutdown sinyalini burada (ana akışta) bekle
	// Bu fonksiyon bloklayıcı (blocking) olmalı
	graceful.Shutdown(a.server.FiberApp(), 5*time.Second, a.repo, a.session, a.rabbit)

	// Eğer server Start sırasında bir hata aldıysa (örn: port meşgul) onu dön
	select {
	case err := <-serverErr:
		return err
	default:
		return nil
	}
}

func setupHttpRouter(cfg *config.Config, repo domain.AuthRepository, sessionRepo domain.SessionRepository, rabbitClient domain.RabbitMQClient) server.RouteRegistrar {

	handler := httptransport.NewHandlers(repo, sessionRepo, rabbitClient)

	return httptransport.NewRouter(handler)
}
