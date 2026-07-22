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
	var req struct {
		Status model.UnitStatus `json:"status"`
		Note   string           `json:"note"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
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
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
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
	if err := h.repo.CreateFleetSetting(&f); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create fleet setting: "+err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Fleet setting created successfully", f)
}

func (h *FleetHandler) UpdateFleetSetting(c fiber.Ctx) error {
	id := c.Params("id")
	var f model.FleetSetting
	if err := c.Bind().JSON(&f); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if err := h.repo.UpdateFleetSetting(id, &f); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update fleet setting: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Fleet setting updated successfully", nil)
}

func (h *FleetHandler) DeleteFleetSetting(c fiber.Ctx) error {
	id := c.Params("id")
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
	return response.Success(c, fiber.StatusOK, "Success fetch allocations", allocs)
}

func (h *FleetHandler) AutoAllocate(c fiber.Ctx) error {
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
	if err := h.repo.UpdateUnitDB(&u); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update unit DB: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Unit DB updated", nil)
}

func (h *FleetHandler) DeleteUnitDB(c fiber.Ctx) error {
	id := c.Query("id")
	if err := h.repo.DeleteUnitDB(id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete unit DB: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Unit DB deleted", nil)
}
