package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev/internal/service"
	"universev/pkg/response"
)

type DisplayHandler struct {
	displaySvc *service.DisplayService
}

func NewDisplayHandler(displaySvc *service.DisplayService) *DisplayHandler {
	return &DisplayHandler{displaySvc: displaySvc}
}

// GetDisplayAttendance returns attendance data for TV display
func (h *DisplayHandler) GetDisplayAttendance(c fiber.Ctx) error {
	rows, err := h.displaySvc.GetDisplayAttendance()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch attendance: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success", rows)
}

// GetDisplayFTW returns fit-to-work data for TV display
func (h *DisplayHandler) GetDisplayFTW(c fiber.Ctx) error {
	rows, err := h.displaySvc.GetDisplayFTW()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch FTW: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success", rows)
}

// GetDisplayFleet returns fleet data for TV display
func (h *DisplayHandler) GetDisplayFleet(c fiber.Ctx) error {
	fleetID := c.Query("fleetId")
	shift := c.Query("shift")
	date := c.Query("date")
	rows, err := h.displaySvc.GetDisplayFleet(fleetID, shift, date)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch fleet: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success", rows)
}

// GetDisplayMonitor returns monitor display data with fleet rotation
func (h *DisplayHandler) GetDisplayMonitor(c fiber.Ctx) error {
	monitorID := c.Query("monitor")
	rows, err := h.displaySvc.GetDisplayMonitor(monitorID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch monitor: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success", rows)
}

// GetDisplayFingerprint returns fingerprint device status for TV display
func (h *DisplayHandler) GetDisplayFingerprint(c fiber.Ctx) error {
	rows, err := h.displaySvc.GetDisplayFingerprint()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch fingerprint devices: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success", rows)
}
