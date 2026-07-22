package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/pkg"
	"universev2-backend/internal/repository"
	"universev2-backend/pkg/response"
)

type ProfileHandler struct {
	userRepo *repository.UserRepo
}

func NewProfileHandler(userRepo *repository.UserRepo) *ProfileHandler {
	return &ProfileHandler{userRepo: userRepo}
}

func (h *ProfileHandler) GetProfile(c fiber.Ctx) error {
	claims, ok := c.Locals("user").(*pkg.JWTCustomClaims)
	if !ok || claims == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	if h.userRepo == nil {
		return response.Error(c, fiber.StatusInternalServerError, "User repository unavailable")
	}

	user, err := h.userRepo.GetByID(claims.UserID)
	if err != nil || user == nil {
		return response.Error(c, fiber.StatusNotFound, "User not found")
	}

	return response.Success(c, fiber.StatusOK, "Success fetch profile", user)
}

func (h *ProfileHandler) UpdateProfile(c fiber.Ctx) error {
	claims, ok := c.Locals("user").(*pkg.JWTCustomClaims)
	if !ok || claims == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if isTrimmedEmpty(req.Name) {
		return sendValidationError(c, "name", "Name is required")
	}

	if h.userRepo != nil {
		user, err := h.userRepo.GetByID(claims.UserID)
		if err == nil && user != nil {
			user.Name = req.Name
			_ = h.userRepo.Update(claims.UserID, user)
		}
	}

	return response.Success(c, fiber.StatusOK, "Profile updated successfully", nil)
}

type UpdatePasswordReq struct {
	OldPassword     string `json:"oldPassword"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

func (h *ProfileHandler) UpdatePassword(c fiber.Ctx) error {
	claims, ok := c.Locals("user").(*pkg.JWTCustomClaims)
	if !ok || claims == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	var req UpdatePasswordReq
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req.OldPassword == "" {
		return sendValidationError(c, "oldPassword", "Old password is required")
	}
	if !isPasswordStrong(req.NewPassword) {
		return sendValidationError(c, "newPassword", "Password must be at least 8 characters and contain letters and numbers")
	}
	if req.NewPassword != req.ConfirmPassword {
		return sendValidationError(c, "confirmPassword", "Passwords do not match")
	}

	if h.userRepo == nil {
		return response.Error(c, fiber.StatusInternalServerError, "User repository unavailable")
	}

	user, err := h.userRepo.GetByID(claims.UserID)
	if err != nil || user == nil {
		return response.Error(c, fiber.StatusNotFound, "User not found")
	}

	// Verify old password
	oldHashFE := hashPasswordFE(req.OldPassword, user.PasswordSalt)
	oldHashLegacy := hashPasswordLegacy(req.OldPassword, user.PasswordSalt)
	if oldHashFE != user.PasswordHash && oldHashLegacy != user.PasswordHash {
		return response.Error(c, fiber.StatusBadRequest, "Incorrect old password")
	}

	newSalt := generateSalt()
	newHash := hashPasswordFE(req.NewPassword, newSalt)

	if err := h.userRepo.UpdatePassword(claims.UserID, newHash, newSalt); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update password: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Password updated successfully", nil)
}
