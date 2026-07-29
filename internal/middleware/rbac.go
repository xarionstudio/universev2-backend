package middleware

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/pkg"
	"universev2-backend/internal/repository"
)

type RBACMiddleware struct {
	roleRepo *repository.RoleRepo
}

func NewRBACMiddleware(roleRepo *repository.RoleRepo) *RBACMiddleware {
	return &RBACMiddleware{roleRepo: roleRepo}
}

// RequirePermission checks that the authenticated user has the required
// permission level on the given module by querying role_permissions from DB.
func (m *RBACMiddleware) RequirePermission(moduleName string, needLevel string) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, ok := c.Locals("user").(*pkg.JWTCustomClaims)
		if !ok || claims == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized — no valid session",
			})
		}

		if len(claims.Roles) == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions — no roles assigned",
			})
		}

		// Query effective permissions from DB for all user's roles
		var roleIDs []uint
		for _, roleStr := range claims.Roles {
			id, err := strconv.ParseUint(roleStr, 10, 64)
			if err == nil {
				roleIDs = append(roleIDs, uint(id))
			}
		}
		perms, err := m.roleRepo.GetPermissionsForRoles(roleIDs)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load permissions",
			})
		}

		userLevel, exists := perms[moduleName]
		if !exists || userLevel == "none" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions for this module",
			})
		}

		// "manage" trumps "view" — if user has manage, they can do view-level too
		if needLevel == "manage" && userLevel != "manage" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions — manage access required",
			})
		}

		return c.Next()
	}
}
