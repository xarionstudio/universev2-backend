package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
	"universev2-backend/pkg/response"
)

type RoleHandler struct {
	repo *repository.RoleRepo
}

func NewRoleHandler(repo *repository.RoleRepo) *RoleHandler {
	return &RoleHandler{repo: repo}
}

func (h *RoleHandler) GetRoles(c fiber.Ctx) error {
	roles, err := h.repo.GetAll()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch roles: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch roles", roles)
}

func (h *RoleHandler) CreateRole(c fiber.Ctx) error {
	var role model.Role
	if err := c.Bind().JSON(&role); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(role.Name) {
		return sendValidationError(c, "name", "Role name is required")
	}

	if err := h.repo.Create(&role); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create role: "+err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Role created successfully", role)
}

func (h *RoleHandler) UpdateRole(c fiber.Ctx) error {
	id := c.Params("id")
	if isTrimmedEmpty(id) {
		return response.Error(c, fiber.StatusBadRequest, "Role ID is required")
	}

	existing, err := h.repo.GetByID(id)
	if err != nil || existing == nil {
		return response.Error(c, fiber.StatusNotFound, "Role not found")
	}

	if existing.IsLocked {
		return response.Error(c, fiber.StatusForbidden, "System locked roles cannot be modified")
	}

	var role model.Role
	if err := c.Bind().JSON(&role); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(role.Name) {
		return sendValidationError(c, "name", "Role name is required")
	}

	if err := h.repo.Update(id, &role); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update role: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Role updated successfully", nil)
}

func (h *RoleHandler) DeleteRole(c fiber.Ctx) error {
	id := c.Params("id")
	if isTrimmedEmpty(id) {
		return response.Error(c, fiber.StatusBadRequest, "Role ID is required")
	}

	existing, err := h.repo.GetByID(id)
	if err != nil || existing == nil {
		return response.Error(c, fiber.StatusNotFound, "Role not found")
	}

	if existing.IsLocked {
		return response.Error(c, fiber.StatusForbidden, "System locked roles cannot be deleted")
	}

	userCount, err := h.repo.CountUsersByRoleID(id)
	if err == nil && userCount > 0 {
		return response.Error(c, fiber.StatusBadRequest, "Cannot delete role assigned to active users")
	}

	if err := h.repo.Delete(id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete role: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Role deleted successfully", nil)
}
