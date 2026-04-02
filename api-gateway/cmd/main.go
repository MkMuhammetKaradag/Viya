package main

import (
	"api-gateway/internal/app"
	"api-gateway/internal/config"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Konfigürasyon yüklenemedi: %v", err)
	}
	application := app.New(cfg)
	application.RegisterService("trips-service", []string{"http://localhost:8081"}, "/api/v1/trips")
	application.RegisterService("waypoints-service", []string{"http://localhost:8081"}, "/api/v1/waypoints")
	application.RegisterService("auth-service", []string{"http://localhost:8082"}, "/api/v1/auth")
	application.RegisterService("user-service", []string{"http://localhost:8083"}, "/api/v1/users")
	// application.RegisterService("auth-service", []string{"http://localhost:8082"}, "/api/v1")

	if err := application.Start(); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
