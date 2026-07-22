package handler

import (
	"github.com/gofiber/fiber/v3"
	"universev2-backend/pkg/response"
)

type FingerprintHandler struct{}

func NewFingerprintHandler() *FingerprintHandler {
	return &FingerprintHandler{}
}

func (h *FingerprintHandler) GetDeviceStatus(c fiber.Ctx) error {
	devices := []fiber.Map{
		{
			"id":        "FP-01",
			"name":      "FP-01 · Office",
			"port":      "/dev/ttyUSB0",
			"baud":      115200,
			"status":    "connected",
			"isOnline":  true,
			"lastPing":  "2s ago",
		},
		{
			"id":        "FP-02",
			"name":      "FP-02 · Gate utara",
			"port":      "/dev/ttyUSB1",
			"baud":      115200,
			"status":    "connected",
			"isOnline":  true,
			"lastPing":  "5s ago",
		},
	}

	return response.Success(c, fiber.StatusOK, "Success fetch devices", devices)
}
