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
	if err := h.repo.Create(&role); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create role: "+err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Role created successfully", role)
}

func (h *RoleHandler) UpdateRole(c fiber.Ctx) error {
	id := c.Params("id")
	var role model.Role
	if err := c.Bind().JSON(&role); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if err := h.repo.Update(id, &role); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update role: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Role updated successfully", nil)
}

func (h *RoleHandler) DeleteRole(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.repo.Delete(id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete role: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Role deleted successfully", nil)
}
