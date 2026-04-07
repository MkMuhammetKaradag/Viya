package main

import (
	"fmt"
	"log"
	"social-service/internal/app"
	"social-service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}
	fmt.Println("config", cfg)
	app, err := app.NewApp(cfg)
	if err != nil {
		panic(err)
	}
	if err := app.Start(); err != nil {
		panic(err)
	}
}
