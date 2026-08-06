package handler

import (
	"fmt"
	"io"
	"strings"

	"github.com/gofiber/fiber/v3"

	"universev/internal/dto"
	"universev/internal/repository"
	"universev/internal/service"
	"universev/pkg/response"
)

type FleetHandler struct {
	fleetSvc *service.FleetService
}

func NewFleetHandler(repo *repository.FleetRepo) *FleetHandler {
	return &FleetHandler{
		fleetSvc: service.NewFleetService(repo),
	}
}

func (h *FleetHandler) GetUnitStatuses(c fiber.Ctx) error {
	units, err := h.fleetSvc.GetUnitStatuses()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch unit statuses: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch unit statuses", units)
}

func (h *FleetHandler) UpdateUnitStatus(c fiber.Ctx) error {
	code := c.Params("code")
	var req dto.UpdateUnitStatusRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.fleetSvc.UpdateUnitStatus(code, req); err != nil {
		msg := err.Error()
		switch msg {
		case "unit code is required":
			return response.Error(c, fiber.StatusBadRequest, "Unit code is required")
		case "status is required":
			return sendValidationError(c, "status", "Status is required")
		case "status must be one of: ready, breakdown, standby":
			return sendValidationError(c, "status", msg)
		default:
			return response.Error(c, fiber.StatusInternalServerError, msg)
		}
	}
	return response.Success(c, fiber.StatusOK, "Unit status updated successfully", nil)
}

func (h *FleetHandler) ReportUnitBreakdown(c fiber.Ctx) error {
	code := c.Params("code")
	var req dto.ReportBreakdownRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.fleetSvc.ReportUnitBreakdown(code, req); err != nil {
		msg := err.Error()
		switch msg {
		case "unit code is required":
			return response.Error(c, fiber.StatusBadRequest, "Unit code is required")
		case "breakdown reason is required":
			return sendValidationError(c, "reason", "Breakdown reason is required")
		default:
			return response.Error(c, fiber.StatusInternalServerError, msg)
		}
	}
	return response.Success(c, fiber.StatusOK, "Unit breakdown reported with required reason", nil)
}

func (h *FleetHandler) GetUnitHistory(c fiber.Ctx) error {
	code := c.Params("code")
	history, err := h.fleetSvc.GetUnitHistory(code)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Unit code is required")
	}
	return response.Success(c, fiber.StatusOK, "Success fetch unit history", history)
}

func (h *FleetHandler) GetFleetSettings(c fiber.Ctx) error {
	fleets, err := h.fleetSvc.GetFleetSettings()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch fleet settings: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch fleet settings", fleets)
}

func (h *FleetHandler) CreateFleetSetting(c fiber.Ctx) error {
	var req dto.CreateFleetSettingRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	f, err := h.fleetSvc.CreateFleetSetting(req)
	if err != nil {
		msg := err.Error()
		if msg == "digger code is required" {
			return sendValidationError(c, "digger", "Digger code is required")
		}
		return response.Error(c, fiber.StatusInternalServerError, msg)
	}
	return response.Success(c, fiber.StatusCreated, "Fleet setting created successfully", f)
}

func (h *FleetHandler) UpdateFleetSetting(c fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateFleetSettingRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.fleetSvc.UpdateFleetSetting(id, req); err != nil {
		msg := err.Error()
		switch msg {
		case "ID is required":
			return response.Error(c, fiber.StatusBadRequest, "ID is required")
		case "digger code is required":
			return sendValidationError(c, "digger", "Digger code is required")
		default:
			return response.Error(c, fiber.StatusInternalServerError, msg)
		}
	}
	return response.Success(c, fiber.StatusOK, "Fleet setting updated successfully", nil)
}

func (h *FleetHandler) DeleteFleetSetting(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.fleetSvc.DeleteFleetSetting(id); err != nil {
		msg := err.Error()
		switch msg {
		case "ID is required":
			return response.Error(c, fiber.StatusBadRequest, "ID is required")
		case "invalid ID":
			return response.Error(c, fiber.StatusBadRequest, "Invalid ID")
		default:
			return response.Error(c, fiber.StatusInternalServerError, msg)
		}
	}
	return response.Success(c, fiber.StatusOK, "Fleet setting deleted successfully", nil)
}

func (h *FleetHandler) GetAllocations(c fiber.Ctx) error {
	date := c.Query("date")
	shift := c.Query("shift")
	result, err := h.fleetSvc.GetAllocations(date, shift)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch allocations: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch allocations", result)
}

func (h *FleetHandler) AutoAllocate(c fiber.Ctx) error {
	var req dto.AutoAllocateRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.fleetSvc.AutoAllocate(req); err != nil {
		msg := err.Error()
		switch msg {
		case "date is required":
			return sendValidationError(c, "date", "Date is required")
		case "shift is required":
			return sendValidationError(c, "shift", "Shift is required")
		default:
			return response.Error(c, fiber.StatusInternalServerError, msg)
		}
	}
	return response.Success(c, fiber.StatusOK, "Auto allocation completed successfully", nil)
}

func (h *FleetHandler) GetUnitDB(c fiber.Ctx) error {
	units, err := h.fleetSvc.GetUnitDB()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch unit DB: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch unit db", units)
}

func (h *FleetHandler) CreateUnitDB(c fiber.Ctx) error {
	var req dto.CreateUnitDBRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	u, err := h.fleetSvc.CreateUnitDB(req)
	if err != nil {
		msg := err.Error()
		if msg == "unit code is required" {
			return sendValidationError(c, "code", "Unit code is required")
		}
		return response.Error(c, fiber.StatusInternalServerError, msg)
	}
	return response.Success(c, fiber.StatusCreated, "Unit DB created", u)
}

func (h *FleetHandler) UpdateUnitDB(c fiber.Ctx) error {
	var req dto.UpdateUnitDBRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.fleetSvc.UpdateUnitDB(req); err != nil {
		msg := err.Error()
		if msg == "unit code is required" {
			return sendValidationError(c, "code", "Unit code is required")
		}
		return response.Error(c, fiber.StatusInternalServerError, msg)
	}
	return response.Success(c, fiber.StatusOK, "Unit DB updated", nil)
}

func (h *FleetHandler) DeleteUnitDB(c fiber.Ctx) error {
	id := c.Query("id")
	if err := h.fleetSvc.DeleteUnitDB(id); err != nil {
		msg := err.Error()
		switch msg {
		case "unit ID is required":
			return response.Error(c, fiber.StatusBadRequest, "Unit ID is required")
		case "invalid unit ID":
			return response.Error(c, fiber.StatusBadRequest, "Invalid unit ID")
		default:
			return response.Error(c, fiber.StatusInternalServerError, msg)
		}
	}
	return response.Success(c, fiber.StatusOK, "Unit DB deleted", nil)
}

// ImportUnitDB godoc
// POST /api/units/db/import
// Content-Type: multipart/form-data
// Field: file (.xlsx)
func (h *FleetHandler) ImportUnitDB(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Excel file is required (field: 'file')")
	}

	fname := strings.ToLower(file.Filename)
	if !strings.HasSuffix(fname, ".xlsx") && !strings.HasSuffix(fname, ".xls") {
		return response.Error(c, fiber.StatusBadRequest, "Only .xlsx or .xls files are accepted")
	}

	const maxSize = 10 << 20 // 10MB
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

	imported, skipped, rowErrors, err := h.fleetSvc.ImportUnitDBFromExcel(data)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Import failed: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK, fmt.Sprintf("Import completed: %d imported, %d updated/skipped", imported, skipped), fiber.Map{
		"imported": imported,
		"skipped":  skipped,
		"errors":   rowErrors,
	})
}
