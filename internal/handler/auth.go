package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/config"
	"universev2-backend/internal/model"
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

func hashPasswordFE(password, salt string) string {
	data := []byte("sha-256/mock:" + salt + ":" + password)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func hashPasswordLegacy(password, salt string) string {
	hHasher := sha256.New()
	hHasher.Write([]byte(password + salt))
	return hex.EncodeToString(hHasher.Sum(nil))
}

func generateSalt() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
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

	if isTrimmedEmpty(req.Email) || req.Password == "" {
		return response.Error(c, fiber.StatusBadRequest, "Email and password are required")
	}

	if h.userRepo == nil {
		return response.Error(c, fiber.StatusInternalServerError, "User repository unavailable")
	}

	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil || user == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	if !user.IsActive {
		return response.Error(c, fiber.StatusForbidden, "Account is inactive")
	}

	// Verify hash with FE format or fallback to legacy format
	computedHash := hashPasswordFE(req.Password, user.PasswordSalt)
	legacyHash := hashPasswordLegacy(req.Password, user.PasswordSalt)

	if computedHash != user.PasswordHash && legacyHash != user.PasswordHash {
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

type RegisterRequest struct {
	Name     string `json:"name"`
	NIK      string `json:"nik"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Dept     string `json:"dept"`
	Pos      string `json:"pos"`
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req RegisterRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid JSON payload")
	}

	if isTrimmedEmpty(req.Name) {
		return sendValidationError(c, "name", "Name is required")
	}
	if isTrimmedEmpty(req.NIK) {
		return sendValidationError(c, "nik", "NIK is required")
	}
	if isTrimmedEmpty(req.Dept) {
		return sendValidationError(c, "dept", "Department is required")
	}
	if isTrimmedEmpty(req.Pos) {
		return sendValidationError(c, "pos", "Position is required")
	}
	if !isValidEmail(req.Email) {
		return sendValidationError(c, "email", "Invalid email format")
	}
	if !isPasswordStrong(req.Password) {
		return sendValidationError(c, "password", "Password must be at least 8 characters and contain both letters and numbers")
	}

	if h.userRepo.ExistsByEmail(req.Email) {
		return response.Error(c, fiber.StatusConflict, "Email is already registered")
	}

	salt := generateSalt()
	hash := hashPasswordFE(req.Password, salt)
	now := time.Now()

	nikPtr := req.NIK
	user := &model.User{
		ID:           fmt.Sprintf("u-%d", now.UnixNano()),
		Email:        req.Email,
		Name:         req.Name,
		NIK:          &nikPtr,
		PasswordHash: hash,
		PasswordSalt: salt,
		IsActive:     true,
		Roles:        []string{"user"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := h.userRepo.Create(user); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create user: "+err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "User registered successfully", user)
}

func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	return response.Success(c, fiber.StatusOK, "Token refreshed successfully", nil)
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	return response.Success(c, fiber.StatusOK, "Logout successful", nil)
}
