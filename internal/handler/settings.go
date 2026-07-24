package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
	"universev2-backend/internal/service"
	"universev2-backend/pkg/response"
)

type SettingsHandler struct {
	settingsSvc *service.SettingsService
	repo        *repository.SettingsRepo
}

func NewSettingsHandler(repo *repository.SettingsRepo) *SettingsHandler {
	return &SettingsHandler{
		settingsSvc: service.NewSettingsService(repo),
		repo:        repo,
	}
}

func (h *SettingsHandler) GetSettings(c fiber.Ctx) error {
	settings, err := h.settingsSvc.GetAppSettings()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch settings: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch settings", settings)
}

func (h *SettingsHandler) UpdateSettings(c fiber.Ctx) error {
	var settings model.AppSettings
	if err := c.Bind().JSON(&settings); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.settingsSvc.UpdateAppSettings(settings); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update settings: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Settings updated successfully", nil)
}

func (h *SettingsHandler) GetAudioSchedules(c fiber.Ctx) error {
	audios, err := h.settingsSvc.GetAudioSchedules()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch audio schedules: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch audio schedules", audios)
}

func (h *SettingsHandler) CreateAudioSchedule(c fiber.Ctx) error {
	var a model.AudioSchedule
	if err := c.Bind().JSON(&a); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.settingsSvc.CreateAudioSchedule(&a); err != nil {
		msg := err.Error()
		switch msg {
		case "audio schedule title is required":
			return sendValidationError(c, "title", msg)
		case "trigger time is required":
			return sendValidationError(c, "when", msg)
		default:
			return response.Error(c, fiber.StatusInternalServerError, "Failed to create audio schedule: "+msg)
		}
	}
	return response.Success(c, fiber.StatusCreated, "Audio schedule created successfully", a)
}

func (h *SettingsHandler) UpdateAudioSchedule(c fiber.Ctx) error {
	id := c.Params("id")
	var a model.AudioSchedule
	if err := c.Bind().JSON(&a); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.settingsSvc.UpdateAudioSchedule(id, &a); err != nil {
		msg := err.Error()
		switch msg {
		case "audio schedule ID is required":
			return response.Error(c, fiber.StatusBadRequest, msg)
		case "audio schedule title is required":
			return sendValidationError(c, "title", msg)
		case "trigger time is required":
			return sendValidationError(c, "when", msg)
		default:
			return response.Error(c, fiber.StatusInternalServerError, "Failed to update audio schedule: "+msg)
		}
	}
	return response.Success(c, fiber.StatusOK, "Audio schedule updated successfully", nil)
}

func (h *SettingsHandler) DeleteAudioSchedule(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.settingsSvc.DeleteAudioSchedule(id); err != nil {
		msg := err.Error()
		if msg == "audio schedule ID is required" {
			return response.Error(c, fiber.StatusBadRequest, msg)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete audio schedule: "+msg)
	}
	return response.Success(c, fiber.StatusOK, "Audio schedule deleted successfully", nil)
}

func (h *SettingsHandler) GetDisplays(c fiber.Ctx) error {
	kind := c.Query("kind")
	displays, err := h.settingsSvc.GetDisplays(kind)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch displays: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch displays", displays)
}

func (h *SettingsHandler) CreateDisplay(c fiber.Ctx) error {
	var d model.DisplayDevice
	if err := c.Bind().JSON(&d); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.settingsSvc.CreateDisplay(&d); err != nil {
		msg := err.Error()
		switch msg {
		case "display name is required":
			return sendValidationError(c, "name", msg)
		case "display location is required":
			return sendValidationError(c, "loc", msg)
		default:
			return response.Error(c, fiber.StatusInternalServerError, "Failed to create display: "+msg)
		}
	}
	return response.Success(c, fiber.StatusCreated, "Display created successfully", d)
}

func (h *SettingsHandler) UpdateDisplay(c fiber.Ctx) error {
	id := c.Params("id")
	var d model.DisplayDevice
	if err := c.Bind().JSON(&d); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.settingsSvc.UpdateDisplay(id, &d); err != nil {
		msg := err.Error()
		switch msg {
		case "display ID is required":
			return response.Error(c, fiber.StatusBadRequest, msg)
		case "display name is required":
			return sendValidationError(c, "name", msg)
		case "display location is required":
			return sendValidationError(c, "loc", msg)
		default:
			return response.Error(c, fiber.StatusInternalServerError, "Failed to update display: "+msg)
		}
	}
	return response.Success(c, fiber.StatusOK, "Display updated successfully", nil)
}

func (h *SettingsHandler) DeleteDisplay(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.settingsSvc.DeleteDisplay(id); err != nil {
		msg := err.Error()
		if msg == "display ID is required" {
			return response.Error(c, fiber.StatusBadRequest, msg)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete display: "+msg)
	}
	return response.Success(c, fiber.StatusOK, "Display deleted successfully", nil)
}

func (h *SettingsHandler) GetDisplayHeartbeat(c fiber.Ctx) error {
	id := c.Params("id")
	if isTrimmedEmpty(id) {
		return response.Error(c, fiber.StatusBadRequest, "Display ID is required")
	}

	_ = h.repo.UpdateHeartbeat(id, "Sekarang")
	data := fiber.Map{
		"online": true,
		"hb":     "1 dtk lalu",
	}
	return response.Success(c, fiber.StatusOK, "Success ping heartbeat", data)
}
