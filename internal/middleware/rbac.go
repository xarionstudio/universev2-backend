package middleware

import (
	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/pkg"
)

// RequirePermission checks that the authenticated user has the required
// permission level on the given module. In this initial implementation
// the check is simplified; a production version would query the DB for
// the user's effective permissions via their role assignments.
func RequirePermission(moduleName string, needLevel string) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, ok := c.Locals("user").(*pkg.JWTCustomClaims)
		if !ok || claims == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized — no valid session",
			})
		}

		// Superadmin bypass — role "r1" has full access to everything
		for _, r := range claims.Roles {
			if r == "r1" {
				return c.Next()
			}
		}

		// Admin role "r2" — manage access on operational modules, none on users/settings
		for _, r := range claims.Roles {
			if r == "r2" {
				if moduleName == "users" || moduleName == "settings" {
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
						"error": "Insufficient permissions for this module",
					})
				}
				return c.Next()
			}
		}

		// Viewer role "r3" — view-only, block manage
		if needLevel == "manage" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions — view-only access",
			})
		}

		return c.Next()
	}
}
