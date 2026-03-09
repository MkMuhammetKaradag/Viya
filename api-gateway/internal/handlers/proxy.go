package handlers

import (
	"api-gateway/internal/service"
	"api-gateway/internal/session"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"
)

type ProxyHandler struct {
	sessionManager *session.SessionManager
	Registry       *service.ServiceRegistry
}

func NewProxyHandler(registry *service.ServiceRegistry, sessionManager *session.SessionManager) *ProxyHandler {
	return &ProxyHandler{
		Registry:       registry,
		sessionManager: sessionManager,
	}
}

func (h *ProxyHandler) Handle(c fiber.Ctx) error {
	if c.Method() == "OPTIONS" {
		
		return c.SendStatus(fiber.StatusNoContent)
	}
	path := c.Path()
	svc, ok := h.Registry.GetByPath(path)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Service not found"})
	}

	baseURL, ok := svc.GetNextBackend()
	if !ok {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Service unhealthy"})
	}

	// 3. Hedef URL'i oluştur (Örn: http://localhost:8081 + /api/v1/trips/...)
	targetURL := baseURL + path
	if queryString := string(c.Request().URI().QueryString()); queryString != "" {
		targetURL += "?" + queryString
	}

	// Setup Request Headers

	if userID, ok := c.Locals("user_id").(string); ok && userID != "" {
		c.Request().Header.Set("X-User-ID", userID)
	}

	c.Request().Header.Set("X-Forwarded-For", c.IP())
	log.Printf("🔀 Proxying to: %s", targetURL)

	if err := proxy.Do(c, targetURL); err != nil {
		log.Printf("❌ Proxy error [%s]: %v", targetURL, err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Backend service error"})
	}
	// err := proxy.Do(c, targetURL)
	c.Response().Header.Set("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Response().Header.Set("Access-Control-Allow-Credentials", "true")

	// if err != nil {
	// 	return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Backend error"})
	// }
	return nil
}
