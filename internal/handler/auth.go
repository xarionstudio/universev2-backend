package handler

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/config"
	"universev2-backend/internal/pkg"
	"universev2-backend/internal/repository"
	"universev2-backend/pkg/response"
)

type AuthHandler struct {
	cfg      *config.Config
	userRepo *repository.UserRepo
	roleRepo *repository.RoleRepo
}

func NewAuthHandler(cfg *config.Config, userRepo *repository.UserRepo, roleRepo *repository.RoleRepo) *AuthHandler {
	return &AuthHandler{
		cfg:      cfg,
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON payload")
	}

	if req.Email == "" || req.Password == "" {
		return response.Error(c, fiber.StatusBadRequest, "Email and password are required")
	}

	if h.userRepo == nil {
		return response.Error(c, fiber.StatusInternalServerError, "User repository unavailable")
	}

	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil || user == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Hash input password with stored salt
	hHasher := sha256.New()
	hHasher.Write([]byte(req.Password + user.PasswordSalt))
	computedHash := hex.EncodeToString(hHasher.Sum(nil))

	if computedHash != user.PasswordHash {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	tokenStr, err := pkg.GenerateToken(user.ID, user.Email, user.Roles, h.cfg.JWTSecret, h.cfg.JWTExpiration)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	var perms map[string]string
	if h.roleRepo != nil {
		perms, _ = h.roleRepo.GetPermissionsForRoles(user.Roles)
	}

	data := fiber.Map{
		"token": tokenStr,
		"user":  user,
		"perms": perms,
	}

	return response.Success(c, fiber.StatusOK, "Login successful", data)
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	return response.Success(c, fiber.StatusCreated, "User registered successfully", nil)
}

func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	return response.Success(c, fiber.StatusOK, "Token refreshed successfully", nil)
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	return response.Success(c, fiber.StatusOK, "Logout successful", nil)
}
