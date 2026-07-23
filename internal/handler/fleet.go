package handler

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
	"universev2-backend/pkg/response"
)

type FleetHandler struct {
	repo *repository.FleetRepo
}

func NewFleetHandler(repo *repository.FleetRepo) *FleetHandler {
	return &FleetHandler{repo: repo}
}

func (h *FleetHandler) GetUnitStatuses(c fiber.Ctx) error {
	units, err := h.repo.GetUnitStatuses()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch unit statuses: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch unit statuses", units)
}

func (h *FleetHandler) UpdateUnitStatus(c fiber.Ctx) error {
	code := c.Params("code")
	if isTrimmedEmpty(code) {
		return response.Error(c, fiber.StatusBadRequest, "Unit code is required")
	}

	var req struct {
		Status model.UnitStatus `json:"status"`
		Note   string           `json:"note"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(string(req.Status)) {
		return sendValidationError(c, "status", "Status is required")
	}
	if req.Status != model.UnitStatusReady && req.Status != model.UnitStatusBreakdown && req.Status != model.UnitStatusStandby {
		return sendValidationError(c, "status", "Status must be one of: ready, breakdown, standby")
	}

	if err := h.repo.UpdateUnitStatus(code, req.Status, req.Note); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update unit status: "+err.Error())
	}

	nowStr := time.Now().Format("02 Jan 15:04")
	_ = h.repo.AddUnitHistory(code, nowStr, string(req.Status), req.Note, string(req.Status))

	return response.Success(c, fiber.StatusOK, "Unit status updated successfully", nil)
}

func (h *FleetHandler) ReportUnitBreakdown(c fiber.Ctx) error {
	code := c.Params("code")
	if isTrimmedEmpty(code) {
		return response.Error(c, fiber.StatusBadRequest, "Unit code is required")
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(req.Reason) {
		return sendValidationError(c, "reason", "Breakdown reason is required")
	}

	if err := h.repo.UpdateUnitStatus(code, model.UnitStatusBreakdown, req.Reason); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to report breakdown: "+err.Error())
	}

	nowStr := time.Now().Format("02 Jan 15:04")
	_ = h.repo.AddUnitHistory(code, nowStr, "Breakdown", req.Reason, "breakdown")

	return response.Success(c, fiber.StatusOK, "Unit breakdown reported with required reason", nil)
}

func (h *FleetHandler) GetUnitHistory(c fiber.Ctx) error {
	code := c.Params("code")
	if isTrimmedEmpty(code) {
		return response.Error(c, fiber.StatusBadRequest, "Unit code is required")
	}

	history, err := h.repo.GetUnitHistory(code)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch unit history: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch unit history", history)
}

func (h *FleetHandler) GetFleetSettings(c fiber.Ctx) error {
	fleets, err := h.repo.GetFleetSettings()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch fleet settings: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch fleet settings", fleets)
}

func (h *FleetHandler) CreateFleetSetting(c fiber.Ctx) error {
	var f model.FleetSetting
	if err := c.Bind().JSON(&f); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(f.Digger) {
		return sendValidationError(c, "digger", "Digger code is required")
	}

	if err := h.repo.CreateFleetSetting(&f); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create fleet setting: "+err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Fleet setting created successfully", f)
}

func (h *FleetHandler) UpdateFleetSetting(c fiber.Ctx) error {
	id := c.Params("id")
	if isTrimmedEmpty(id) {
		return response.Error(c, fiber.StatusBadRequest, "ID is required")
	}

	var f model.FleetSetting
	if err := c.Bind().JSON(&f); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(f.Digger) {
		return sendValidationError(c, "digger", "Digger code is required")
	}

	if err := h.repo.UpdateFleetSetting(id, &f); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update fleet setting: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Fleet setting updated successfully", nil)
}

func (h *FleetHandler) DeleteFleetSetting(c fiber.Ctx) error {
	id := c.Params("id")
	if isTrimmedEmpty(id) {
		return response.Error(c, fiber.StatusBadRequest, "ID is required")
	}

	if err := h.repo.DeleteFleetSetting(id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete fleet setting: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Fleet setting deleted successfully", nil)
}

func (h *FleetHandler) GetAllocations(c fiber.Ctx) error {
	date := c.Query("date")
	shift := c.Query("shift")
	allocs, err := h.repo.GetAllocations(date, shift)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch allocations: "+err.Error())
	}

	// Transform to FE FaAlloc format: { "2026-07-23": { "pagi": { "EX7001": "503264133", ... }, "malam": { ... } } }
	result := make(model.FleetAllocResponse)
	for _, a := range allocs {
		if result[a.Date] == nil {
			result[a.Date] = make(map[string]map[string]string)
		}
		result[a.Date][a.Shift] = a.Units
	}
	return response.Success(c, fiber.StatusOK, "Success fetch allocations", result)
}

func (h *FleetHandler) AutoAllocate(c fiber.Ctx) error {
	var req struct {
		Date  string `json:"date"`
		Shift string `json:"shift"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(req.Date) {
		return sendValidationError(c, "date", "Date is required")
	}
	if isTrimmedEmpty(req.Shift) {
		return sendValidationError(c, "shift", "Shift is required")
	}

	if err := h.repo.AutoAllocate(req.Date, req.Shift); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to auto allocate fleets: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Auto allocation completed successfully", nil)
}

func (h *FleetHandler) GetUnitDB(c fiber.Ctx) error {
	units, err := h.repo.GetUnitDB()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch unit DB: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch unit db", units)
}

func (h *FleetHandler) CreateUnitDB(c fiber.Ctx) error {
	var u model.UnitDb
	if err := c.Bind().JSON(&u); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(u.Code) {
		return sendValidationError(c, "code", "Unit code is required")
	}

	if err := h.repo.CreateUnitDB(&u); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create unit DB: "+err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Unit DB created", u)
}

func (h *FleetHandler) UpdateUnitDB(c fiber.Ctx) error {
	var u model.UnitDb
	if err := c.Bind().JSON(&u); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(u.Code) {
		return sendValidationError(c, "code", "Unit code is required")
	}

	if err := h.repo.UpdateUnitDB(&u); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update unit DB: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Unit DB updated", nil)
}

func (h *FleetHandler) DeleteUnitDB(c fiber.Ctx) error {
	id := c.Query("id")
	if isTrimmedEmpty(id) {
		return response.Error(c, fiber.StatusBadRequest, "Unit ID is required")
	}

	if err := h.repo.DeleteUnitDB(id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete unit DB: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Unit DB deleted", nil)
}
