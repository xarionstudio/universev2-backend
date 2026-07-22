package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
	"universev2-backend/pkg/response"
)

type SettingsHandler struct {
	repo *repository.SettingsRepo
}

func NewSettingsHandler(repo *repository.SettingsRepo) *SettingsHandler {
	return &SettingsHandler{repo: repo}
}

func (h *SettingsHandler) GetSettings(c fiber.Ctx) error {
	settings, err := h.repo.GetAppSettings()
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

	if isTrimmedEmpty(settings.AppName) {
		settings.AppName = "universev2-backend"
	}

	if err := h.repo.UpdateAppSettings(settings); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update settings: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Settings updated successfully", nil)
}

func (h *SettingsHandler) GetAudioSchedules(c fiber.Ctx) error {
	audios, err := h.repo.GetAudioSchedules()
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

	if isTrimmedEmpty(a.Title) {
		return sendValidationError(c, "title", "Audio schedule title is required")
	}
	if isTrimmedEmpty(a.When) {
		return sendValidationError(c, "when", "Trigger time is required")
	}

	if err := h.repo.CreateAudioSchedule(&a); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create audio schedule: "+err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Audio schedule created successfully", a)
}

func (h *SettingsHandler) UpdateAudioSchedule(c fiber.Ctx) error {
	id := c.Params("id")
	if isTrimmedEmpty(id) {
		return response.Error(c, fiber.StatusBadRequest, "Audio schedule ID is required")
	}

	var a model.AudioSchedule
	if err := c.Bind().JSON(&a); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(a.Title) {
		return sendValidationError(c, "title", "Audio schedule title is required")
	}
	if isTrimmedEmpty(a.When) {
		return sendValidationError(c, "when", "Trigger time is required")
	}

	if err := h.repo.UpdateAudioSchedule(id, &a); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update audio schedule: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Audio schedule updated successfully", nil)
}

func (h *SettingsHandler) DeleteAudioSchedule(c fiber.Ctx) error {
	id := c.Params("id")
	if isTrimmedEmpty(id) {
		return response.Error(c, fiber.StatusBadRequest, "Audio schedule ID is required")
	}

	if err := h.repo.DeleteAudioSchedule(id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete audio schedule: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Audio schedule deleted successfully", nil)
}

func (h *SettingsHandler) GetDisplays(c fiber.Ctx) error {
	kind := c.Query("kind")
	displays, err := h.repo.GetDisplays(kind)
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

	if isTrimmedEmpty(d.Name) {
		return sendValidationError(c, "name", "Display name is required")
	}
	if isTrimmedEmpty(d.Loc) {
		return sendValidationError(c, "loc", "Display location is required")
	}

	if err := h.repo.CreateDisplay(&d); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create display: "+err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Display created successfully", d)
}

func (h *SettingsHandler) UpdateDisplay(c fiber.Ctx) error {
	id := c.Params("id")
	if isTrimmedEmpty(id) {
		return response.Error(c, fiber.StatusBadRequest, "Display ID is required")
	}

	var d model.DisplayDevice
	if err := c.Bind().JSON(&d); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(d.Name) {
		return sendValidationError(c, "name", "Display name is required")
	}
	if isTrimmedEmpty(d.Loc) {
		return sendValidationError(c, "loc", "Display location is required")
	}

	if err := h.repo.UpdateDisplay(id, &d); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update display: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Display updated successfully", nil)
}

func (h *SettingsHandler) DeleteDisplay(c fiber.Ctx) error {
	id := c.Params("id")
	if isTrimmedEmpty(id) {
		return response.Error(c, fiber.StatusBadRequest, "Display ID is required")
	}

	if err := h.repo.DeleteDisplay(id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete display: "+err.Error())
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
