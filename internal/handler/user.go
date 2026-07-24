package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/dto"
	"universev2-backend/internal/repository"
	"universev2-backend/internal/service"
	"universev2-backend/pkg/response"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userRepo *repository.UserRepo) *UserHandler {
	return &UserHandler{
		userSvc: service.NewUserService(userRepo, nil),
	}
}

func (h *UserHandler) GetUsers(c fiber.Ctx) error {
	users, err := h.userSvc.GetUsers()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch users: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch users", users)
}

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := h.userSvc.CreateUser(req)
	if err != nil {
		msg := err.Error()
		switch msg {
		case "name is required":
			return sendValidationError(c, "name", "Name is required")
		case "invalid email format":
			return sendValidationError(c, "email", "Invalid email format")
		case "at least one role is required":
			return sendValidationError(c, "roles", "At least one role is required")
		case "email is already in use":
			return response.Error(c, fiber.StatusConflict, "Email is already in use")
		default:
			return response.Error(c, fiber.StatusInternalServerError, msg)
		}
	}
	return response.Success(c, fiber.StatusCreated, "User created successfully", user)
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.userSvc.UpdateUser(id, req); err != nil {
		msg := err.Error()
		switch msg {
		case "user not found":
			return response.Error(c, fiber.StatusNotFound, "User not found")
		case "name is required":
			return sendValidationError(c, "name", "Name is required")
		case "invalid email format":
			return sendValidationError(c, "email", "Invalid email format")
		case "at least one role is required":
			return sendValidationError(c, "roles", "At least one role is required")
		default:
			return response.Error(c, fiber.StatusInternalServerError, msg)
		}
	}
	return response.Success(c, fiber.StatusOK, "User updated successfully", nil)
}

func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.userSvc.DeleteUser(id); err != nil {
		return response.Error(c, fiber.StatusNotFound, "User not found")
	}
	return response.Success(c, fiber.StatusOK, "User deleted successfully", nil)
}
