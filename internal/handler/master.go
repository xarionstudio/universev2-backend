package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/model"
	"universev2-backend/internal/service"
	"universev2-backend/pkg/response"
)

type MasterHandler struct {
	masterSvc *service.MasterService
}

func NewMasterHandler(masterSvc *service.MasterService) *MasterHandler {
	return &MasterHandler{masterSvc: masterSvc}
}

func (h *MasterHandler) GetMasterByCategory(c fiber.Ctx) error {
	cat := c.Params("category")
	if isTrimmedEmpty(cat) {
		return response.Error(c, fiber.StatusBadRequest, "Category parameter is required")
	}

	page, perPage := service.ParsePaginationParams(
		c.Query("page"),
		c.Query("perPage"),
	)
	search := c.Query("search")

	result, err := h.masterSvc.GetMasterByCategory(cat, page, perPage, search)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch master entries: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch master data", result)
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

	created, err := h.masterSvc.BulkCreate(cat, []model.MdEntry{entry})
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create master entry: "+err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Master entry created", created[0])
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

	entry.ID = id
	_, err := h.masterSvc.BulkCreate(string(entry.Cat), []model.MdEntry{entry})
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update master entry: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Master entry updated", nil)
}

func (h *MasterHandler) DeleteMasterEntry(c fiber.Ctx) error {
	id := c.Params("id")
	if isTrimmedEmpty(id) {
		return response.Error(c, fiber.StatusBadRequest, "Entry ID is required")
	}

	if err := h.masterSvc.BulkDelete([]string{id}); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete master entry: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Master entry deleted", nil)
}
