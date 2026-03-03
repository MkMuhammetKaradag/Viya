package main

import "auth-service/internal/app"

func main() {

	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}
	if err := app.Start(); err != nil {
		panic(err)
	}
}
