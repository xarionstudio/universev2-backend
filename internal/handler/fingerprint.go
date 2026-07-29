package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/config"
	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
	"universev2-backend/internal/worker"
	"universev2-backend/pkg/response"
)

type FingerprintHandler struct {
	cfg     *config.Config
	fpRepo  *repository.FingerprintRepo
	worker  *worker.FingerprintWorker
}

func NewFingerprintHandler(cfg *config.Config, fpRepo *repository.FingerprintRepo, worker *worker.FingerprintWorker) *FingerprintHandler {
	return &FingerprintHandler{
		cfg:    cfg,
		fpRepo: fpRepo,
		worker: worker,
	}
}

// GetDeviceStatus — GET /api/fingerprint/devices
func (h *FingerprintHandler) GetDeviceStatus(c fiber.Ctx) error {
	devices, err := h.fpRepo.GetAllDevices()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch fingerprint devices: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch devices", devices)
}

// CreateDevice — POST /api/fingerprint/devices
func (h *FingerprintHandler) CreateDevice(c fiber.Ctx) error {
	var dev model.FingerprintDevice
	if err := c.Bind().JSON(&dev); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if dev.Code == "" || dev.IPAddress == "" {
		return response.Error(c, fiber.StatusBadRequest, "Device code and IP address are required")
	}

	if dev.Port <= 0 {
		dev.Port = 80
	}

	if err := h.fpRepo.CreateDevice(&dev); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create fingerprint device: "+err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Fingerprint device created successfully", dev)
}

// UpdateDevice — PUT /api/fingerprint/devices/:id
func (h *FingerprintHandler) UpdateDevice(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid device ID")
	}

	existing, err := h.fpRepo.GetDeviceByID(uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Fingerprint device not found")
	}

	var req model.FingerprintDevice
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Code != "" {
		existing.Code = req.Code
	}
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.IPAddress != "" {
		existing.IPAddress = req.IPAddress
	}
	if req.Port > 0 {
		existing.Port = req.Port
	}
	if req.Location != "" {
		existing.Location = req.Location
	}
	existing.IsActive = req.IsActive

	if err := h.fpRepo.UpdateDevice(existing); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update device: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Device updated successfully", existing)
}

// DeleteDevice — DELETE /api/fingerprint/devices/:id
func (h *FingerprintHandler) DeleteDevice(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid device ID")
	}

	if err := h.fpRepo.DeleteDevice(uint(id)); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete device: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Device deleted successfully", nil)
}

// SyncNow — POST /api/fingerprint/sync
func (h *FingerprintHandler) SyncNow(c fiber.Ctx) error {
	if h.worker == nil {
		return response.Error(c, fiber.StatusInternalServerError, "Fingerprint worker is not initialized")
	}

	totalSynced := h.worker.SyncAllDevices()
	return response.Success(c, fiber.StatusOK, "Fingerprint sync completed", fiber.Map{
		"totalSynced": totalSynced,
	})
}
