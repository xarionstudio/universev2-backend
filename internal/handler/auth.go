package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev/internal/config"
	"universev/internal/dto"
	"universev/internal/repository"
	"universev/internal/service"
	"universev/pkg/response"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(cfg *config.Config, userRepo *repository.UserRepo, roleRepo *repository.RoleRepo) *AuthHandler {
	return &AuthHandler{
		authSvc: service.NewAuthService(cfg, userRepo, roleRepo),
	}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON payload")
	}

	result, err := h.authSvc.Login(req)
	if err != nil {
		msg := err.Error()
		switch msg {
		case "invalid credentials":
			return response.Error(c, fiber.StatusUnauthorized, "Email or password is incorrect")
		case "account is inactive":
			return response.Error(c, fiber.StatusForbidden, "Account is inactive")
		case "invalid email format":
			return sendValidationError(c, "email", "Invalid email format")
		default:
			return response.Error(c, fiber.StatusBadRequest, msg)
		}
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    result.Token,
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
	})

	data := fiber.Map{
		"token": result.Token,
		"user":  result.User,
		"perms": result.Perms,
	}
	return response.Success(c, fiber.StatusOK, "Login successful", data)
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON payload")
	}

	result, err := h.authSvc.Register(req)
	if err != nil {
		msg := err.Error()
		switch msg {
		case "email is already registered":
			return response.Error(c, fiber.StatusConflict, "Email is already registered")
		case "NIK must be exactly 9 digits":
			return sendValidationError(c, "nik", "NIK must be exactly 9 digits")
		case "invalid email format":
			return sendValidationError(c, "email", "Invalid email format")
		case "password must be at least 8 characters and contain both letters and numbers":
			return sendValidationError(c, "password", "Password must be at least 8 characters and contain both letters and numbers")
		case "password must not exceed 72 characters":
			return sendValidationError(c, "password", "Password must not exceed 72 characters")
		case "name must not exceed 100 characters":
			return sendValidationError(c, "name", "Name must not exceed 100 characters")
		case "NIK must not exceed 9 characters":
			return sendValidationError(c, "nik", "NIK must not exceed 9 characters")
		default:
			if len(msg) > 0 && msg[0] >= 'a' && msg[0] <= 'z' {
				return sendValidationError(c, "field", msg)
			}
			return response.Error(c, fiber.StatusInternalServerError, msg)
		}
	}

	return response.Success(c, fiber.StatusCreated, "User registered successfully", result.User)
}

func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	tokenStr := ""
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		parts := splitAuthHeader(authHeader)
		if len(parts) == 2 {
			tokenStr = parts[1]
		}
	}
	if tokenStr == "" {
		tokenStr = c.Cookies("jwt")
	}

	if tokenStr == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Missing authorization token or cookie")
	}

	result, err := h.authSvc.RefreshToken(tokenStr)
	if err != nil {
		msg := err.Error()
		switch msg {
		case "invalid or expired token":
			return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired token")
		case "user not found":
			return response.Error(c, fiber.StatusUnauthorized, "User not found")
		case "account is inactive":
			return response.Error(c, fiber.StatusForbidden, "Account is inactive")
		default:
			return response.Error(c, fiber.StatusInternalServerError, msg)
		}
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    result.Token,
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
	})

	data := fiber.Map{
		"token": result.Token,
		"user":  result.User,
		"perms": result.Perms,
	}
	return response.Success(c, fiber.StatusOK, "Token refreshed successfully", data)
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
	})
	return response.Success(c, fiber.StatusOK, "Logout successful", nil)
}
