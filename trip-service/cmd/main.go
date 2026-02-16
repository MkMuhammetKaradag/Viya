package main

import (
	"trip-service/app"
	"trip-service/config"
)

func main() {
	cfg, _ := config.Read()
	app, err := app.NewApp(cfg)
	if err != nil {
		panic(err)
	}
	if err := app.Run(); err != nil {
		panic(err)
	}
}
