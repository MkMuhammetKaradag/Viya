package main

import (
	"fmt"
	"trip-service/internal/app"
	"trip-service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("config load error: %w", err))
	}
	fmt.Println("config", cfg)
	app, err := app.NewApp(cfg)
	if err != nil {
		panic(err)
	}
	if err := app.Run(); err != nil {
		panic(err)
	}
}
