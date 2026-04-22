package app

import (
	"fmt"
	"log"
	"time"
	"user-service/internal/config"
	"user-service/internal/database"
	"user-service/internal/domain"
	"user-service/internal/graceful"
	"user-service/internal/infrastructure/img"
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
	cldSvc, err := img.NewCloudinary(cfg.Cloudinary.CloudName, cfg.Cloudinary.APIKey, cfg.Cloudinary.APISecret)

	if err != nil {
		return nil, err
	}
	router := rabbitmq.NewRabbitRouter(repo)

	rabbitClient, err := initMessaging()
	if err != nil {
		return nil, fmt.Errorf("init rabbit :%w", err)
	}

	httpRouter := setupHttpRouter(cfg, repo, cldSvc, rabbitClient)

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
func setupHttpRouter(cfg *config.Config, userRepo domain.UserRepository, cloudinaryService domain.CloudinaryService, rabbitClient domain.RabbitMQClient) server.RouterRegister {
	handler := httptransport.NewHandlers(userRepo, cloudinaryService, rabbitClient)
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

	go func() {
		log.Println("RabbitMQ Consumer başlatılıyor...")
		if err := a.rabbit.ConsumeMessages(a.rabbitRouter.Route); err != nil {
			log.Printf("Consumer hatası: %v", err)
		}
	}()
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("HTTP Server %s portunda başlatılıyor...", a.server.Address())
		if err := a.server.Start(); err != nil {
			serverErr <- err
		}
	}()
	graceful.Shutdown(a.server.FiberApp(), 5*time.Second, a.userRepo, a.rabbit)
	select {
	case err := <-serverErr:
		return err
	default:
		log.Println("Uygulama başarıyla kapatıldı.")
		return nil
	}
}
