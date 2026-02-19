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

	httptransport "trip-service/internal/transport/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/hibiken/asynq"
)

type App struct {
	config    *config.Config
	registrar RouteRegistrar
	// Add your application fields here
}

func NewApp(cfg *config.Config) (*App, error) {
	c, err := buildContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("bootstrap failed: %w", err)
	}
	return &App{config: cfg, registrar: c.httpRouter}, nil
}

type container struct {
	tripRepo   domain.TripRepository
	httpRouter RouteRegistrar
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
	redisOpt := asynq.RedisClientOpt{Addr: "localhost:6379", DB: 2}
	asynqClient := asynq.NewClient(redisOpt)
	wrk := worker.NewWorker(asynqClient)

	processor := worker.NewTaskProcessor(redisOpt, repo, imgSvc)
	go func() {
		if err := processor.Start(); err != nil {
			log.Printf("Task Processor error: %v", err)
		}
	}()

	httpRouter := setupHttpRouter(cfg, repo, imgSvc, wrk)
	return &container{
		tripRepo:   repo,
		httpRouter: httpRouter,
	}, nil
}
func getServerConfig(cfg *config.Config) AppConfig {
	return AppConfig{
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

type AppConfig struct {
	Port         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
type RouteRegistrar interface {
	Register(app *fiber.App)
}

func (a *App) Run() error {
	cfg := getServerConfig(a.config)
	app := fiber.New(fiber.Config{
		AppName:      "Viya Trip Service v1.0",
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Concurrency:  256 * 1024,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"},
		AllowCredentials: true,
	}))

	if a.registrar != nil {
		a.registrar.Register(app)
	}
	app.Get("/health", func(c fiber.Ctx) error {
		// return c.SendString("OK")
		return c.JSON(fiber.Map{
			"status":  "active",
			"service": "trip-service",
		})
	})

	return app.Listen(cfg.Port)

}

func setupHttpRouter(cfg *config.Config, r domain.TripRepository, i domain.ImageService, w domain.Worker) RouteRegistrar {

	httpHandlers := httptransport.NewHandlers(r, i, w)
	return httptransport.NewRouter(httpHandlers)
}
