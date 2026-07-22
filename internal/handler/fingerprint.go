package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/config"
	"universev2-backend/internal/repository"
	"universev2-backend/pkg/response"
)

type FingerprintHandler struct {
	cfg          *config.Config
	settingsRepo *repository.SettingsRepo
}

func NewFingerprintHandler(cfg *config.Config, settingsRepo *repository.SettingsRepo) *FingerprintHandler {
	return &FingerprintHandler{
		cfg:          cfg,
		settingsRepo: settingsRepo,
	}
}

func (h *FingerprintHandler) GetDeviceStatus(c fiber.Ctx) error {
	devices := []fiber.Map{}

	if h.settingsRepo != nil {
		ds, err := h.settingsRepo.GetDisplays("finger")
		if err == nil {
			for _, d := range ds {
				devices = append(devices, fiber.Map{
					"id":       d.ID,
					"name":     d.Name,
					"port":     h.cfg.FingerprintPort,
					"baud":     h.cfg.FingerprintBaud,
					"status":   "connected",
					"isOnline": d.Online,
					"lastPing": d.Heartbeat,
				})
			}
		}
	}

	return response.Success(c, fiber.StatusOK, "Success fetch devices", devices)
}
