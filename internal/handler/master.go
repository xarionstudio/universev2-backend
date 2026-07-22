package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
	"universev2-backend/pkg/response"
)

type MasterHandler struct {
	repo *repository.MasterRepo
}

func NewMasterHandler(repo *repository.MasterRepo) *MasterHandler {
	return &MasterHandler{repo: repo}
}

func (h *MasterHandler) GetMasterByCategory(c fiber.Ctx) error {
	cat := c.Params("category")
	if isTrimmedEmpty(cat) {
		return response.Error(c, fiber.StatusBadRequest, "Category parameter is required")
	}

	entries, err := h.repo.GetByCategory(cat)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch master entries: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch master data", entries)
}

func (h *MasterHandler) CreateMasterEntry(c fiber.Ctx) error {
	cat := c.Params("category")
	if isTrimmedEmpty(cat) {
		return response.Error(c, fiber.StatusBadRequest, "Category parameter is required")
	}

	var entry model.MdEntry
	if err := c.Bind().JSON(&entry); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(entry.Name) {
		return sendValidationError(c, "name", "Name is required")
	}

	entry.Cat = model.MdCat(cat)
	if err := h.repo.Create(&entry); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create master entry: "+err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Master entry created", entry)
}

func (h *MasterHandler) UpdateMasterEntry(c fiber.Ctx) error {
	id := c.Params("id")
	if isTrimmedEmpty(id) {
		return response.Error(c, fiber.StatusBadRequest, "Entry ID is required")
	}

	var entry model.MdEntry
	if err := c.Bind().JSON(&entry); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(entry.Name) {
		return sendValidationError(c, "name", "Name is required")
	}

	if err := h.repo.Update(id, &entry); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update master entry: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Master entry updated", nil)
}

func (h *MasterHandler) DeleteMasterEntry(c fiber.Ctx) error {
	id := c.Params("id")
	if isTrimmedEmpty(id) {
		return response.Error(c, fiber.StatusBadRequest, "Entry ID is required")
	}

	if err := h.repo.Delete(id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete master entry: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Master entry deleted", nil)
}
