// trip-service/internal/app/app.go
package app

import (
	"fmt"
	"log"
	"time"
	"trip-service/infrastructure/img"
	"trip-service/infrastructure/worker"
	"trip-service/internal/config"
	"trip-service/internal/database"
	"trip-service/internal/domain"
	"trip-service/internal/graceful"
	"trip-service/internal/server"
	"viya/pkg/messaging"

	httptransport "trip-service/internal/transport/http"
	"trip-service/internal/transport/rabbitmq"

	"github.com/hibiken/asynq"
)

type App struct {
	config       *config.Config
	processor    *worker.TaskProcessor
	server       *server.Server
	repo         domain.TripRepository
	rabbit       domain.RabbitMQClient
	rabbitRouter domain.RabbitRouter
	// Add your application fields here
}

func NewApp(cfg *config.Config) (*App, error) {
	c, err := buildContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("bootstrap failed: %w", err)
	}
	return &App{config: cfg,
		processor:    c.processor,
		server:       c.server,
		rabbit:       c.rabbit,
		rabbitRouter: c.rabbitRouter,
	}, nil
}

type container struct {
	tripRepo     domain.TripRepository
	processor    *worker.TaskProcessor
	server       *server.Server
	repo         domain.TripRepository
	rabbit       domain.RabbitMQClient
	rabbitRouter domain.RabbitRouter
}

func buildContainer(cfg *config.Config) (*container, error) {
	repo, err := initStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}
	imgSvc, err := img.NewCloudinaryService(cfg.Cloudinary.CloudName, cfg.Cloudinary.APIKey, cfg.Cloudinary.APISecret)
	if err != nil {
		return nil, err
	}
	router := rabbitmq.NewRabbitRouter(repo)
	rabbitClient, err := initMessaging()
	if err != nil {
		return nil, fmt.Errorf("init rabbit :%w", err)
	}

	redisOpt := asynq.RedisClientOpt{Addr: "localhost:6379", DB: 2, Password: "password"}

	asynqClient := asynq.NewClient(redisOpt)
	wrk := worker.NewWorker(asynqClient)

	processor := worker.NewTaskProcessor(redisOpt, repo, imgSvc)
	go func() {
		log.Println("Starting Task Processor on Redis DB 2...")
		if err := processor.Start(); err != nil {
			log.Fatalf("Task Processor fatal error: %v", err)
		}
	}()

	httpRouter := setupHttpRouter(cfg, repo, imgSvc, wrk)
	return &container{
		tripRepo:     repo,
		processor:    processor,
		server:       server.NewServer(getServerConfig(cfg), httpRouter),
		repo:         repo,
		rabbit:       rabbitClient,
		rabbitRouter: router,
	}, nil
}
func getServerConfig(cfg *config.Config) server.Config {
	return server.Config{
		Port:         cfg.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
func initStorage(cfg *config.Config) (domain.TripRepository, error) {
	repo, err := database.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres init error: %w", err)
	}

	return repo, nil
}

func initMessaging() (domain.RabbitMQClient, error) {
	rabbitCfg := messaging.NewDefaultConfig("")
	// Servis adını ve diğer ayarları ver
	return messaging.NewRabbitClient(rabbitCfg, messaging.TripService)

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
	graceful.Shutdown(a.server.FiberApp(), 5*time.Second, a.processor, a.repo)
	select {
	case err := <-serverErr:
		return err
	default:
		return nil
	}
}
func setupHttpRouter(cfg *config.Config, r domain.TripRepository, i domain.ImageService, w domain.Worker) server.RouteRegistrar {

	httpHandlers := httptransport.NewHandlers(r, i, w)
	return httptransport.NewRouter(httpHandlers)
}
