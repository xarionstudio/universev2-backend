package handler

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
	"universev2-backend/pkg/response"
)

type UserHandler struct {
	repo *repository.UserRepo
}

func NewUserHandler(repo *repository.UserRepo) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) GetUsers(c fiber.Ctx) error {
	users, err := h.repo.GetAll()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch users: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch users", users)
}

type CreateUserRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	NIK      string   `json:"nik"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(req.Name) {
		return sendValidationError(c, "name", "Name is required")
	}
	if !isValidEmail(req.Email) {
		return sendValidationError(c, "email", "Invalid email format")
	}
	if len(req.Roles) == 0 {
		return sendValidationError(c, "roles", "At least one role is required")
	}
	if h.repo.ExistsByEmail(req.Email) {
		return response.Error(c, fiber.StatusConflict, "Email is already in use")
	}

	// Only allow setting whitelisted fields to prevent mass assignment
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
		Roles:        req.Roles,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := h.repo.Create(user); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create user: "+err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "User created successfully", user)
}

type UpdateUserRequest struct {
	Name  string   `json:"name"`
	Email string   `json:"email"`
	NIK   string   `json:"nik"`
	Roles []string `json:"roles"`
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	id := c.Params("id")
	if isTrimmedEmpty(id) {
		return response.Error(c, fiber.StatusBadRequest, "User ID is required")
	}

	existing, err := h.repo.GetByID(id)
	if err != nil || existing == nil {
		return response.Error(c, fiber.StatusNotFound, "User not found")
	}

	var req UpdateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(req.Name) {
		return sendValidationError(c, "name", "Name is required")
	}
	if !isValidEmail(req.Email) {
		return sendValidationError(c, "email", "Invalid email format")
	}
	if len(req.Roles) == 0 {
		return sendValidationError(c, "roles", "At least one role is required")
	}

	// Only update whitelisted fields to prevent mass assignment
	existing.Name = req.Name
	existing.Email = req.Email
	if req.NIK != "" {
		existing.NIK = &req.NIK
	}
	existing.Roles = req.Roles

	if err := h.repo.Update(id, existing); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update user: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "User updated successfully", nil)
}

func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")
	if isTrimmedEmpty(id) {
		return response.Error(c, fiber.StatusBadRequest, "User ID is required")
	}

	existing, err := h.repo.GetByID(id)
	if err != nil || existing == nil {
		return response.Error(c, fiber.StatusNotFound, "User not found")
	}

	if err := h.repo.Delete(id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete user: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "User deleted successfully", nil)
}
