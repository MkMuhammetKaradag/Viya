package handlers

import (
	"api-gateway/internal/service"
	"api-gateway/internal/session"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"

	"github.com/valyala/fasthttp"
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

	userID := c.Locals("user_id").(string)
	c.Request().Header.Set("X-User-ID", userID)
	c.Request().Header.Set("X-Forwarded-For", c.IP())
	log.Printf("🔀 Proxying to: %s", targetURL)

	if err := proxy.Do(c, targetURL, &fasthttp.Client{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}); err != nil {

		log.Printf("❌ Proxy error [%s]: %v", targetURL, err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Backend service error"})
	}

	return nil
}
