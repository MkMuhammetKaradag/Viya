package app

import (
	"fmt"
	"log"
	"time"
	"user-service/internal/config"
	"user-service/internal/database"
	"user-service/internal/domain"
	"user-service/internal/graceful"
	"user-service/internal/server"
	httptransport "user-service/internal/transport/http"
	"user-service/internal/transport/rabbitmq"
	"viya/pkg/messaging"
)

type App struct {
	server   *server.Server
	config   *config.Config
	userRepo domain.UserRepository
	rabbit   domain.RabbitMQClient
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
		rabbit:   c.rabbit,
	}, nil
}

type container struct {
	server   *server.Server
	userRepo domain.UserRepository
	rabbit   domain.RabbitMQClient
}

func buildContainer(cfg *config.Config) (*container, error) {
	repo, _, err := initStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("init progres repository:%w", err)
	}

	rabbitRouter := setupRabbitRouter(cfg, repo)
	messageRouter := func(msg messaging.Message) error {
		handler, ok := rabbitRouter[msg.Type]
		if !ok {
			return nil
		}
		return handler.Handle(msg)
	}
	rabbitClient, err := initMessaging(messageRouter)
	if err != nil {
		return nil, fmt.Errorf("init rabbit :%w", err)
	}

	httpRouter := setupHttpRouter(cfg, repo)

	return &container{
		userRepo: repo,
		server: server.NewServer(
			getServerConfig(cfg),
			httpRouter,
		),
		rabbit: rabbitClient,
	}, nil
}

func initStorage(cfg *config.Config) (domain.UserRepository, domain.RabbitMQClient, error) {
	repo, err := database.NewRepository(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("postgres init error: %w", err)
	}

	// rabbitClient, err := messaging.NewRabbitClient(cfg.RabbitMQ.URL)
	// if err != nil {
	// 	return nil, nil, fmt.Errorf("rabbitmq init error: %w", err)
	// }
	return repo, nil, nil
}
func initMessaging(handler func(messaging.Message) error) (domain.RabbitMQClient, error) {
	config := messaging.NewDefaultConfig("")
	config.RetryTypes = []messaging.MessageType{}
	rabbitMQ, err := messaging.NewRabbitClient(config, messaging.UserService)
	if err != nil {
		log.Fatalf("RabbitMQ bağlantısı kurulamadı: %v", err)
		return nil, err
	}

	go func() {

		err = rabbitMQ.ConsumeMessages(func(msg messaging.Message) error {

			return handler(msg)

		})
		if err != nil {
			log.Fatal("Mesaj dinleyici başlatılamadı:", err)
		}

	}()
	return rabbitMQ, nil

}
func setupHttpRouter(cfg *config.Config, userRepo domain.UserRepository) server.RouterRegister {
	handler := httptransport.NewHandlers(userRepo)
	return httptransport.NewRouter(handler)
}

func setupRabbitRouter(cfg *config.Config, userRepo domain.UserRepository) map[messaging.MessageType]domain.MessageHandler {
	return rabbitmq.SetupMessageHandlers(userRepo)
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
