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
	server       *server.Server
	config       *config.Config
	userRepo     domain.UserRepository
	rabbit       domain.RabbitMQClient
	rabbitRouter domain.RabbitRouter
}

func NewApp(cfg *config.Config) (*App, error) {
	c, err := buildContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("app failed:%w", err)
	}
	return &App{
		config:       cfg,
		userRepo:     c.userRepo,
		server:       c.server,
		rabbit:       c.rabbit,
		rabbitRouter: c.rabbitRouter,
	}, nil
}

type container struct {
	server       *server.Server
	userRepo     domain.UserRepository
	rabbit       domain.RabbitMQClient
	rabbitRouter domain.RabbitRouter
}

func buildContainer(cfg *config.Config) (*container, error) {
	repo, _, err := initStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("init progres repository:%w", err)
	}

	router := rabbitmq.NewRabbitRouter(repo)

	rabbitClient, err := initMessaging()
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
		rabbit:       rabbitClient,
		rabbitRouter: router,
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
func initMessaging() (domain.RabbitMQClient, error) {
	rabbitCfg := messaging.NewDefaultConfig("")
	// Servis adını ve diğer ayarları ver
	return messaging.NewRabbitClient(rabbitCfg, messaging.UserService)

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
		log.Println("RabbitMQ Consumer başlatılıyor...")
		if err := a.rabbit.ConsumeMessages(a.rabbitRouter.Route); err != nil {
			log.Printf("Consumer hatası: %v", err)
		}
	}()
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
