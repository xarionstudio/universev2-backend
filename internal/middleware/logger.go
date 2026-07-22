package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
)

// LoggerMiddleware logs each request with method, path, status, and latency
func LoggerMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Path()

		color := "\033[32m" // green
		if status >= 400 && status < 500 {
			color = "\033[33m" // yellow
		} else if status >= 500 {
			color = "\033[31m" // red
		}
		reset := "\033[0m"

		fmt.Printf("%s[%d]%s %s %s — %s\n", color, status, reset, method, path, latency.Round(time.Millisecond))

		return err
	}
}
