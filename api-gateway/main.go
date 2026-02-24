package main

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"
)

func main() {
	app := fiber.New()

	// Auth Middleware (Sadece Gateway'de çalışır)
	//app.Use("/api/v1/trips", authMiddleware)

	// Proxy (İstekleri Trip Service'e yönlendirir)
	app.All("/api/v1/trips/*", func(c fiber.Ctx) error {
		// Auth'dan gelen user_id'yi Header'a ekle
		//userID := c.Locals("user_id").(string)
		c.Request().Header.Set("X-User-ID", "userID")
		fmt.Println("geldi istek ", c.Path())
		// Trip Service'in adresi (Örn: localhost:8081)
		targetURL := "http://localhost:8081" + c.Path()

		// Buradan isteği ileten bir proxy kütüphanesi (fiber/proxy) kullanabilirsin
		return proxy.Do(c, targetURL)
	})

	app.Listen(":8080")
}
