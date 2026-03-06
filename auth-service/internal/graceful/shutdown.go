package graceful

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
)

// Releasable kapatılabilir kaynaklar için bir interface
type Releasable interface {
	Close() error
}

func Shutdown(app *fiber.App, timeout time.Duration, resources ...Releasable) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\n--- 🛑 Shutdown signal received ---")

	// 1. HTTP Sunucusunu durdur (Yeni istekleri keser)
	if err := app.ShutdownWithTimeout(timeout); err != nil {
		fmt.Printf("Fiber shutdown error: %v\n", err)
	}

	// 2. Diğer tüm kaynakları sırayla kapat (Postgres, Redis vb.)
	for _, res := range resources {
		if res != nil {
			if err := res.Close(); err != nil {
				fmt.Printf("Resource close error: %v\n", err)
			}
		}
	}

	fmt.Println("✅ Server gracefully stopped. Goodbye!")
}
