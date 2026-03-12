package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
	"google.golang.org/grpc"
)

type RouteRegistrar interface {
	Register(app *fiber.App)
}
type GrpcServerRegistrar interface {
	Register(server *grpc.Server)
}
type ServerConfig struct {
	GrpcPort     string
	Port         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Server struct {
	app        *fiber.App
	cfg        ServerConfig
	grpcServer *grpc.Server
}

func NewServer(cfg ServerConfig, registrar RouteRegistrar, grpcRegistrar GrpcServerRegistrar) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Concurrency:  256 * 1024,
	})
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "active",
			"service": "auth-service",
		})
	})
	if registrar != nil {
		registrar.Register(app)
	}

	grpcSrv := grpc.NewServer()

	// GrpcRegistrar varsa, implementasyonları kaydet
	if grpcRegistrar != nil {
		grpcRegistrar.Register(grpcSrv)
	}

	return &Server{
		app:        app,
		cfg:        cfg,
		grpcServer: grpcSrv,
	}
}

func (s *Server) Start() error {
	return s.app.Listen(s.Address())
}

func (s *Server) Address() string {
	go func() {
		if err := s.startGrpc(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("gRPC sunucusu hatası: %v", err)
		}
	}()

	fmt.Println("listen!!!!!")

	return fmt.Sprintf("0.0.0.0:%s", s.cfg.Port)
}

func (s *Server) Shutdown(timeout time.Duration) error {

	s.grpcServer.GracefulStop()
	log.Println("gRPC sunucusu durduruldu.")
	return s.app.ShutdownWithTimeout(timeout)
}
func (s *Server) FiberApp() *fiber.App {
	return s.app
}
func (s *Server) startGrpc() error {
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", s.cfg.GrpcPort))
	if err != nil {
		return fmt.Errorf("gRPC dinlemede hata: %w", err)
	}
	log.Printf("👂 gRPC sunucusu %s adresinde dinliyor...", s.cfg.GrpcPort)
	// Bloklayan çağrı: Sunucu çalışmaya başlar
	return s.grpcServer.Serve(listen)
}
