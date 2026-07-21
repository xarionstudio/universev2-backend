package middleware

import (
	"github.com/gofiber/fiber/v3"
)

func RBAC(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Forbidden",
			})
		}
		return c.Next()
	}
}
