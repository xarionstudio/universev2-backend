package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev/internal/dto"
	internalpkg "universev/internal/pkg"
	"universev/internal/repository"
	"universev/internal/service"
	"universev/pkg/response"
)

type ProfileHandler struct {
	profileSvc *service.ProfileService
}

func NewProfileHandler(userRepo *repository.UserRepo) *ProfileHandler {
	return &ProfileHandler{
		profileSvc: service.NewProfileService(userRepo),
	}
}

func (h *ProfileHandler) GetProfile(c fiber.Ctx) error {
	claims, ok := c.Locals("user").(*internalpkg.JWTCustomClaims)
	if !ok || claims == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	user, err := h.profileSvc.GetProfile(claims.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "User not found")
	}

	return response.Success(c, fiber.StatusOK, "Success fetch profile", user)
}

func (h *ProfileHandler) UpdateProfile(c fiber.Ctx) error {
	claims, ok := c.Locals("user").(*internalpkg.JWTCustomClaims)
	if !ok || claims == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	var req dto.UpdateProfileRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.profileSvc.UpdateProfile(claims.UserID, req); err != nil {
		msg := err.Error()
		switch msg {
		case "name is required":
			return sendValidationError(c, "name", "Name is required")
		case "invalid email format":
			return sendValidationError(c, "email", "Invalid email format")
		case "email is already in use":
			return response.Error(c, fiber.StatusConflict, "Email is already in use")
		default:
			return response.Error(c, fiber.StatusInternalServerError, msg)
		}
	}

	return response.Success(c, fiber.StatusOK, "Profile updated successfully", nil)
}

func (h *ProfileHandler) UpdatePassword(c fiber.Ctx) error {
	claims, ok := c.Locals("user").(*internalpkg.JWTCustomClaims)
	if !ok || claims == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	var req dto.UpdatePasswordRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.profileSvc.UpdatePassword(claims.UserID, req); err != nil {
		msg := err.Error()
		switch msg {
		case "old password is required":
			return sendValidationError(c, "oldPassword", "Old password is required")
		case "password must be at least 8 characters and contain letters and numbers":
			return sendValidationError(c, "newPassword", "Password must be at least 8 characters and contain letters and numbers")
		case "passwords do not match":
			return sendValidationError(c, "confirmPassword", "Passwords do not match")
		case "incorrect old password":
			return response.Error(c, fiber.StatusBadRequest, "Incorrect old password")
		default:
			return response.Error(c, fiber.StatusInternalServerError, msg)
		}
	}

	return response.Success(c, fiber.StatusOK, "Password updated successfully", nil)
}
