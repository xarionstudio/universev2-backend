package handler

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"

	"universev/internal/model"
	"universev/internal/service"
	"universev/pkg/response"
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

	// Parse raw JSON into map
	var body map[string]interface{}
	if err := c.Bind().JSON(&body); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	name, _ := body["name"].(string)
	if isTrimmedEmpty(name) {
		return sendValidationError(c, "name", "Name is required")
	}

	// Build per-category struct from generic map
	entry := buildEntryFromMap(cat, body)
	if entry == nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid category")
	}

	created, err := h.masterSvc.BulkCreate(cat, []interface{}{entry})
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create master entry: "+err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Master entry created", created)
}

func (h *MasterHandler) UpdateMasterEntry(c fiber.Ctx) error {
	code := c.Params("id")
	if isTrimmedEmpty(code) {
		return response.Error(c, fiber.StatusBadRequest, "Entry code is required")
	}
	cat := c.Params("category")
	if isTrimmedEmpty(cat) {
		return response.Error(c, fiber.StatusBadRequest, "Category parameter is required")
	}

	var body map[string]interface{}
	if err := c.Bind().JSON(&body); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Build updates map with only the fields that were sent
	updates := buildUpdatesMap(cat, body)
	if len(updates) == 0 {
		return response.Error(c, fiber.StatusBadRequest, "No valid fields to update")
	}

	if err := h.masterSvc.UpdateEntry(cat, code, updates); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update master entry: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Master entry updated", nil)
}

func (h *MasterHandler) DeleteMasterEntry(c fiber.Ctx) error {
	code := c.Params("id")
	if isTrimmedEmpty(code) {
		return response.Error(c, fiber.StatusBadRequest, "Entry code is required")
	}
	cat := c.Params("category")
	if isTrimmedEmpty(cat) {
		return response.Error(c, fiber.StatusBadRequest, "Category parameter is required")
	}

	if err := h.masterSvc.BulkDelete(cat, []string{code}); err != nil {
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

// ── Helpers ───────────────────────────────────────────────────────────────────

func buildEntryFromMap(cat string, body map[string]interface{}) interface{} {
	name, _ := body["name"].(string)
	active := true
	if v, ok := body["active"].(bool); ok {
		active = v
	}

	switch cat {
	case "egi":
		code, _ := body["code"].(string)
		return &model.MasterEGIType{Code: code, Name: name, Active: active}
	case "product":
		code, _ := body["code"].(string)
		return &model.MasterProduct{Code: code, Name: name, Active: active}
	case "eqclass":
		code, _ := body["code"].(string)
		desc, _ := body["description"].(string)
		return &model.MasterEqClass{Code: code, Name: name, Description: desc, Active: active}
	case "area":
		code, _ := body["code"].(string)
		category, _ := body["category"].(string)
		return &model.MasterArea{Code: code, Name: name, Category: category, Active: active}
	case "tempudo":
		code, _ := body["code"].(string)
		loc, _ := body["location"].(string)
		pickup, _ := body["pickupType"].(string)
		return &model.MasterTempudo{Code: code, Name: name, Location: loc, PickupType: pickup, Active: active}
	case "bus":
		code, _ := body["code"].(string)
		egiType, _ := body["egiType"].(string)
		depTime, _ := body["departureTime"].(string)
		return &model.MasterBus{Code: code, Name: name, EGIType: egiType, DepartureTime: depTime, Active: active}
	case "lokasiex":
		code, _ := body["code"].(string)
		busCode, _ := body["busCode"].(string)
		tempudoCode, _ := body["tempudoCode"].(string)
		return &model.MasterLocationEx{Code: code, Name: name, BusCode: busCode, TempudoCode: tempudoCode, Active: active}
	case "mess":
		code, _ := body["code"].(string)
		block, _ := body["block"].(string)
		return &model.MasterMess{Code: code, Name: name, Block: block, Active: active}
	case "runtext":
		code, _ := body["code"].(string)
		target, _ := body["targetDisplay"].(string)
		color, _ := body["textColor"].(string)
		return &model.MasterRunningText{Code: code, Name: name, TargetDisplay: target, TextColor: color, Active: active}
	}
	return nil
}

func buildUpdatesMap(cat string, body map[string]interface{}) map[string]interface{} {
	updates := make(map[string]interface{})

	if v, ok := body["name"].(string); ok {
		updates["name"] = v
	}
	if v, ok := body["active"].(bool); ok {
		updates["is_active"] = v
	}

	switch cat {
	case "egi", "product":
		if v, ok := body["code"].(string); ok {
			updates["code"] = v
		}
	case "eqclass":
		if v, ok := body["description"].(string); ok {
			updates["description"] = v
		}
	case "area":
		if v, ok := body["category"].(string); ok {
			updates["category"] = v
		}
	case "tempudo":
		if v, ok := body["location"].(string); ok {
			updates["location"] = v
		}
		if v, ok := body["pickupType"].(string); ok {
			updates["pickup_type"] = v
		}
	case "bus":
		if v, ok := body["egiType"].(string); ok {
			updates["egi_type"] = v
		}
		if v, ok := body["departureTime"].(string); ok {
			updates["departure_time"] = v
		}
	case "lokasiex":
		if v, ok := body["busCode"].(string); ok {
			updates["bus_code"] = v
		}
		if v, ok := body["tempudoCode"].(string); ok {
			updates["tempudo_code"] = v
		}
	case "mess":
		if v, ok := body["block"].(string); ok {
			updates["block"] = v
		}
	case "runtext":
		if v, ok := body["targetDisplay"].(string); ok {
			updates["target_display"] = v
		}
		if v, ok := body["textColor"].(string); ok {
			updates["text_color"] = v
		}
	}

	return updates
}
