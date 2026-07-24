package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
	"universev2-backend/internal/service"
	"universev2-backend/pkg/response"
)

type RoleHandler struct {
	roleSvc *service.RoleService
}

func NewRoleHandler(repo *repository.RoleRepo) *RoleHandler {
	return &RoleHandler{
		roleSvc: service.NewRoleService(repo),
	}
}

func (h *RoleHandler) GetRoles(c fiber.Ctx) error {
	roles, err := h.roleSvc.GetRoles()
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

	if err := h.roleSvc.CreateRole(&role); err != nil {
		msg := err.Error()
		if msg == "role name is required" {
			return sendValidationError(c, "name", msg)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create role: "+msg)
	}
	return response.Success(c, fiber.StatusCreated, "Role created successfully", role)
}

func (h *RoleHandler) UpdateRole(c fiber.Ctx) error {
	id := c.Params("id")
	var role model.Role
	if err := c.Bind().JSON(&role); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.roleSvc.UpdateRole(id, &role); err != nil {
		msg := err.Error()
		switch msg {
		case "role not found":
			return response.Error(c, fiber.StatusNotFound, "Role not found")
		case "role ID is required":
			return response.Error(c, fiber.StatusBadRequest, "Role ID is required")
		case "system locked roles cannot be modified":
			return response.Error(c, fiber.StatusForbidden, msg)
		case "role name is required":
			return sendValidationError(c, "name", msg)
		default:
			return response.Error(c, fiber.StatusInternalServerError, "Failed to update role: "+msg)
		}
	}
	return response.Success(c, fiber.StatusOK, "Role updated successfully", nil)
}

func (h *RoleHandler) DeleteRole(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.roleSvc.DeleteRole(id); err != nil {
		msg := err.Error()
		switch msg {
		case "role not found":
			return response.Error(c, fiber.StatusNotFound, "Role not found")
		case "role ID is required":
			return response.Error(c, fiber.StatusBadRequest, "Role ID is required")
		case "system locked roles cannot be deleted":
			return response.Error(c, fiber.StatusForbidden, msg)
		case "cannot delete role assigned to active users":
			return response.Error(c, fiber.StatusBadRequest, msg)
		default:
			return response.Error(c, fiber.StatusInternalServerError, "Failed to delete role: "+msg)
		}
	}
	return response.Success(c, fiber.StatusOK, "Role deleted successfully", nil)
}
