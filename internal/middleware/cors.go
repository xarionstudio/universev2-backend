package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"universev2-backend/internal/config"
)

// CORSMiddleware configures CORS to allow frontend origins from config
func CORSMiddleware(cfg *config.Config) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     []string{cfg.CORSAllowedOrigins},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours preflight cache
	})
}
