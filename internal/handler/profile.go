package handler

import (
	"github.com/gofiber/fiber/v3"
	"universev2-backend/internal/pkg"
	"universev2-backend/pkg/response"
)

type ProfileHandler struct{}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{}
}

func (h *ProfileHandler) GetProfile(c fiber.Ctx) error {
	claims, ok := c.Locals("user").(*pkg.JWTCustomClaims)
	if !ok || claims == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	data := fiber.Map{
		"id":    claims.UserID,
		"email": claims.Email,
		"roles": claims.Roles,
		// Data mock
		"name": "First Angel Paustine",
		"nik":  "503264133",
	}

	return response.Success(c, fiber.StatusOK, "Success fetch profile", data)
}

func (h *ProfileHandler) UpdateProfile(c fiber.Ctx) error {
	return response.Success(c, fiber.StatusOK, "Profile updated successfully", nil)
}

func (h *ProfileHandler) UpdatePassword(c fiber.Ctx) error {
	return response.Success(c, fiber.StatusOK, "Password updated successfully", nil)
}
