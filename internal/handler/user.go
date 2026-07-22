package handler

import (
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

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	var user model.User
	if err := c.Bind().JSON(&user); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if !isValidEmail(user.Email) {
		return sendValidationError(c, "email", "Invalid email format")
	}
	if len(user.Roles) == 0 {
		return sendValidationError(c, "roles", "At least one role is required")
	}
	if h.repo.ExistsByEmail(user.Email) {
		return response.Error(c, fiber.StatusConflict, "Email is already in use")
	}

	if err := h.repo.Create(&user); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create user: "+err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "User created successfully", user)
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

	var user model.User
	if err := c.Bind().JSON(&user); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if !isValidEmail(user.Email) {
		return sendValidationError(c, "email", "Invalid email format")
	}
	if len(user.Roles) == 0 {
		return sendValidationError(c, "roles", "At least one role is required")
	}

	if err := h.repo.Update(id, &user); err != nil {
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
