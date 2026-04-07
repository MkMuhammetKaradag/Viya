// trip-service/internal/app/app.go
package app

import (
	"fmt"
	"log"
	"time"

	"social-service/internal/config"
	"social-service/internal/database"
	"social-service/internal/domain"
	"social-service/internal/graceful"
	"social-service/internal/server"
	httptransport "social-service/internal/transport/http"
)

type App struct {
	config *config.Config

	server *server.Server
	repo   domain.SocialRepository
}

func NewApp(cfg *config.Config) (*App, error) {
	c, err := buildContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("bootstrap failed: %w", err)
	}
	return &App{config: cfg,
		server: c.server,

		repo: c.repo,
	}, nil
}

type container struct {
	server *server.Server
	repo   domain.SocialRepository
}

func buildContainer(cfg *config.Config) (*container, error) {
	repo, err := initStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}

	httpRouter := setupHttpRouter(cfg, repo)
	return &container{

		server: server.NewServer(getServerConfig(cfg), httpRouter),
		repo:   repo,
	}, nil
}
func getServerConfig(cfg *config.Config) server.Config {
	return server.Config{
		Port:         cfg.Server.Port,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
func initStorage(cfg *config.Config) (domain.SocialRepository, error) {
	repo, err := database.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres init error: %w", err)
	}

	return repo, nil
}

func (a *App) Start() error {

	// 2. HTTP Server'ı başlat (Arka planda)
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("HTTP Server %s portunda başlatılıyor...", a.server.Address())
		if err := a.server.Start(); err != nil {
			serverErr <- err
		}
	}()

	// 3. Graceful Shutdown (BLOKLAYICI OLMALI)
	// Bu fonksiyon içeride os.Interrupt (SIGINT, SIGTERM) bekler.
	// Ctrl+C gelene kadar BURADA BEKLER.
	graceful.Shutdown(a.server.FiberApp(), 5*time.Second, a.repo)

	// Ctrl+C gelince kod buraya düşer.
	// Eğer server başlatılırken bir hata oluştuysa onu kontrol et, yoksa temizce çık.
	select {
	case err := <-serverErr:
		return err
	default:
		log.Println("Uygulama başarıyla kapatıldı.")
		return nil
	}
}
func setupHttpRouter(cfg *config.Config, r domain.SocialRepository) server.RouteRegistrar {

	httpHandlers := httptransport.NewHandlers(r)
	return httptransport.NewRouter(httpHandlers)
}
