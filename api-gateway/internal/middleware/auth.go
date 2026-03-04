package middleware

import (
	"api-gateway/internal/session"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware(protectedPrefixes []string, sessionManager *session.SessionManager) fiber.Handler {
	return func(c fiber.Ctx) error {
		requestPath := c.Path()
		isProtected := false
		for _, prefix := range protectedPrefixes {
			if strings.HasPrefix(requestPath, prefix) {
				isProtected = true
				break
			}
		}

		if !isProtected {
			return c.Next()
		}
		authValue := c.Cookies("session_id")
		if authHeader := c.Get("Authorization"); strings.HasPrefix(authHeader, "Bearer ") {
			authValue = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if authValue == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Authentication required"})
		}

		// 3. Redis/Session kontrolü
		sess, err := sessionManager.GetSession(c.Context(), authValue)
		if err != nil {
			fmt.Println("error", err)
			return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired session"})
		}

		// 4. User ID'yi hem context'e hem Header'a koy (Proxy için kritik!)
		c.Locals("user_id", sess.UserID)
		c.Request().Header.Set("X-User-ID", sess.UserID)
		return c.Next()
	}
}
