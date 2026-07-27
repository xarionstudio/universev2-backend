package handler

import (
	"fmt"
	"io"
	"strings"
	"time"

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

// ImportMaster godoc
// POST /api/master/:category/import
// Content-Type: multipart/form-data
// Field: file (.xlsx)
func (h *MasterHandler) ImportMaster(c fiber.Ctx) error {
	cat := c.Params("category")
	if isTrimmedEmpty(cat) {
		return response.Error(c, fiber.StatusBadRequest, "Category parameter is required")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Excel file is required (field: 'file')")
	}

	fname := strings.ToLower(file.Filename)
	if !strings.HasSuffix(fname, ".xlsx") && !strings.HasSuffix(fname, ".xls") {
		return response.Error(c, fiber.StatusBadRequest, "Only .xlsx or .xls files are accepted")
	}

	const maxSize = 10 << 20 // 10 MB
	if file.Size > maxSize {
		return response.Error(c, fiber.StatusBadRequest, "File size exceeds 10MB limit")
	}

	f, err := file.Open()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to open uploaded file")
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to read file")
	}

	imported, err := h.masterSvc.ImportFromExcel(cat, data)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Import failed: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK,
		fmt.Sprintf("Import completed: %d entries imported into category '%s'", imported, cat),
		fiber.Map{"imported": imported, "category": cat},
	)
}

// ExportMaster godoc
// GET /api/master/:category/export
// Returns an xlsx file download of all entries in the given category.
func (h *MasterHandler) ExportMaster(c fiber.Ctx) error {
	cat := c.Params("category")
	if isTrimmedEmpty(cat) {
		return response.Error(c, fiber.StatusBadRequest, "Category parameter is required")
	}

	xlsxData, err := h.masterSvc.ExportToExcel(cat)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to generate export: "+err.Error())
	}

	fileName := fmt.Sprintf("master_%s_%s.xlsx", cat, time.Now().Format("20060102"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	return c.Send(xlsxData)
}
