// trip-service/internal/graceful/shutdown.go
package graceful

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
)

type Releasable interface {
	Close() error
}

func Shutdown(app *fiber.App, timeout time.Duration, resources ...Releasable) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	fmt.Println("--- 🛑 shutdowm signal recived ---")

	if err := app.ShutdownWithTimeout(timeout); err != nil {
		fmt.Printf(" fiber shutdowm error: %v", err)
	}

	for _, res := range resources {
		if res != nil {
			if err := res.Close(); err != nil {
				fmt.Println("resorve close error :%v", err)
			}
		}
	}
	fmt.Println("✅server graccefully stopped. Goodbye!")
}
