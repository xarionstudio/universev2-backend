package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/service"
	"universev2-backend/pkg/response"
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
	rows, err := h.displaySvc.GetDisplayFleet(fleetID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch fleet: "+err.Error())
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
