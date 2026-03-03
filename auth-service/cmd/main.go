package main

import (
	"auth-service/internal/app"
	"auth-service/internal/config"
	"fmt"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("config load error: %w", err))
	}
	app, err := app.NewApp(cfg)
	if err != nil {
		panic(err)
	}
	if err := app.Start(); err != nil {
		panic(err)
	}
}
