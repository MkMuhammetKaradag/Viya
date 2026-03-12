package middleware

import (
	"api-gateway/internal/grpc_client"
	"api-gateway/internal/pb"
	"api-gateway/internal/session"
	"fmt"
	"strings"
	"time"

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

		resp, err := grpc_client.ValidateAndRotateSession(c.Context(), &pb.ValidateRequest{
			SessionId: authValue,
			Ip:        c.IP(),
			UserAgent: c.Get("User-Agent"),
		})

		if err != nil || !resp.Valid {
			c.Cookie(&fiber.Cookie{
				Name:     "session_id",
				Value:    "",
				MaxAge:   -1,
				Expires:  time.Now().Add(-1 * time.Hour),
				HTTPOnly: true,
				Secure:   false,
				Path:     "/",
				SameSite: "Lax",
			})
			return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired session"})
		}

		if resp.NewSessionId != "" {
			c.Locals("new_session_id", resp.NewSessionId)
			c.Cookie(&fiber.Cookie{
				Name:     "session_id",
				Value:    resp.NewSessionId,
				HTTPOnly: true,
				Secure:   false, // Üretimde true olmalı
				Expires:  time.Now().Add(24 * time.Hour),
				SameSite: "Lax",
				Path:     "/",
			})
			fmt.Println("🔄 Session yenilendi:", resp.NewSessionId)
		}
		// 3. Redis/Session kontrolü
		// sess, err := sessionManager.GetSession(c.Context(), authValue)
		// if err != nil {
		// 	fmt.Println("error", err)
		// 	return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired session"})
		// }

		// 4. User ID'yi hem context'e hem Header'a koy (Proxy için kritik!)
		fmt.Println("resp:", resp)
		c.Locals("user_id", resp.UserId)
		c.Request().Header.Set("X-User-ID", resp.UserId)
		return c.Next()
	}
}
