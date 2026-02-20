// trip-service/internal/graceful/shutdown.go
package graceful

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trip-service/infrastructure/worker"
	"trip-service/internal/domain"

	"github.com/gofiber/fiber/v3"
)

func WaitForShutdown(app *fiber.App, processor *worker.TaskProcessor, repo domain.TripRepository, timeout time.Duration) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\n--- Shutdown signal received ---")

	// 1. Önce sunucuyu kapat (Yeni istek gelmesin)
	if err := app.ShutdownWithTimeout(timeout); err != nil {
		fmt.Printf("Fiber shutdown error: %v\n", err)
	}

	// 2. Worker'ı durdur (Devam eden işler tamamlansın)

	fmt.Println("Shutting down worker processor...")
	processor.Stop()

	fmt.Println("Server gracefully stopped. Goodbye! 👋")
}
