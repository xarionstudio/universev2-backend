package handler

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"

	"universev/internal/model"
	"universev/internal/repository"
	"universev/internal/service"
	"universev/pkg/response"
)

type SettingsHandler struct {
	settingsSvc *service.SettingsService
	repo        *repository.SettingsRepo
	uploadDir   string
}

func NewSettingsHandler(repo *repository.SettingsRepo, uploadDir string) *SettingsHandler {
	return &SettingsHandler{
		settingsSvc: service.NewSettingsService(repo),
		repo:        repo,
		uploadDir:   uploadDir,
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

// UploadLogo handles logo file upload for company branding.
// POST /api/settings/logo  (multipart/form-data, field: "file")
func (h *SettingsHandler) UploadLogo(c fiber.Ctx) error {
	return h.uploadBrandFile(c, "logo")
}

// UploadFavicon handles favicon file upload for browser tab icon.
// POST /api/settings/favicon  (multipart/form-data, field: "file")
func (h *SettingsHandler) UploadFavicon(c fiber.Ctx) error {
	return h.uploadBrandFile(c, "favicon")
}

// uploadBrandFile is a shared helper for logo & favicon uploads.
func (h *SettingsHandler) uploadBrandFile(c fiber.Ctx, kind string) error {
	file, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "File is required (field: 'file')")
	}

	// Sanitize filename to prevent path traversal
	safeFilename := strings.ReplaceAll(file.Filename, "..", "")
	safeFilename = strings.ReplaceAll(safeFilename, "/", "")
	safeFilename = strings.ReplaceAll(safeFilename, "\\", "")
	ext := strings.ToLower(filepath.Ext(safeFilename))

	// Validate by MIME and/or extension (browsers often send empty/octet-stream for .ico)
	contentType := file.Header.Get("Content-Type")
	allowedMIME := map[string]bool{
		"image/jpeg":               true,
		"image/jpg":                true,
		"image/png":                true,
		"image/svg+xml":            true,
		"image/x-icon":             true,
		"image/vnd.microsoft.icon": true,
		"image/webp":               true,
		"application/octet-stream": true, // allow if extension is valid
	}
	allowedExt := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".svg": true, ".ico": true, ".webp": true,
	}
	if !allowedExt[ext] && !allowedMIME[contentType] {
		return response.Error(c, fiber.StatusBadRequest, "Only image files (PNG, JPG, SVG, ICO, WEBP) are allowed")
	}
	if ext != "" && !allowedExt[ext] {
		return response.Error(c, fiber.StatusBadRequest, "Only image files (PNG, JPG, SVG, ICO, WEBP) are allowed")
	}

	// Validate file size (max 5MB)
	const maxFileSize = 5 << 20 // 5MB
	if file.Size > maxFileSize {
		return response.Error(c, fiber.StatusBadRequest, "File size exceeds maximum limit of 5MB")
	}

	// Save to uploads/branding directory
	brandDir := filepath.Join(h.uploadDir, "branding")
	if err := os.MkdirAll(brandDir, 0755); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create upload directory: "+err.Error())
	}

	// Use fixed filename so old file gets replaced (avoid stale cache/DB bugs)
	if ext == "" {
		switch contentType {
		case "image/png":
			ext = ".png"
		case "image/svg+xml":
			ext = ".svg"
		case "image/x-icon", "image/vnd.microsoft.icon":
			ext = ".ico"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".jpg"
		}
	}
	fileName := kind + ext
	filePath := filepath.Join(brandDir, fileName)
	if err := c.SaveFile(file, filePath); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to save file: "+err.Error())
	}

	// Build a browser-accessible URL
	relPath, err := filepath.Rel(h.uploadDir, filePath)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to resolve file path: "+err.Error())
	}
	fileURL := "/uploads/" + filepath.ToSlash(relPath)

	// Only logo updates companyLogo in app_settings.
	// Favicon is served from a fixed path and must not overwrite the logo URL.
	if kind == "logo" {
		settings := model.AppSettings{CompanyLogo: fileURL}
		if err := h.settingsSvc.UpdateAppSettings(settings); err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "Failed to save branding URL: "+err.Error())
		}
	}

	return response.Success(c, fiber.StatusOK, "File uploaded successfully", fiber.Map{
		"url": fileURL,
	})
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
	code := c.Params("code")
	if isTrimmedEmpty(code) {
		return response.Error(c, fiber.StatusBadRequest, "Display code is required")
	}

	_ = h.repo.UpdateHeartbeat(code, "Sekarang")

	// Read actual device status from DB
	dev, err := h.repo.GetDeviceByCode(code)
	if err != nil {
		// Device not found — return default online status
		return response.Success(c, fiber.StatusOK, "Success ping heartbeat", fiber.Map{
			"online": true,
			"hb":     "1 dtk lalu",
		})
	}

	hb := "1 dtk lalu"
	if dev.Heartbeat != "" {
		hb = dev.Heartbeat
	}
	return response.Success(c, fiber.StatusOK, "Success ping heartbeat", fiber.Map{
		"online": dev.Online,
		"hb":     hb,
	})
}

// Business Rules

func (h *SettingsHandler) GetAllBusinessRules(c fiber.Ctx) error {
	rules, err := h.settingsSvc.GetAllBusinessRules()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch business rules: "+err.Error())
	}

	// Parse JSON rules into map
	result := make([]model.BusinessRulesResponse, 0, len(rules))
	for _, rule := range rules {
		var rulesMap map[string]interface{}
		if err := json.Unmarshal([]byte(rule.Rules), &rulesMap); err != nil {
			continue
		}
		result = append(result, model.BusinessRulesResponse{
			Category: rule.Category,
			Rules:    rulesMap,
		})
	}

	return response.Success(c, fiber.StatusOK, "Success fetch business rules", result)
}

func (h *SettingsHandler) GetBusinessRule(c fiber.Ctx) error {
	category := c.Params("category")
	if isTrimmedEmpty(category) {
		return response.Error(c, fiber.StatusBadRequest, "Category is required")
	}

	rule, err := h.settingsSvc.GetBusinessRuleByCategory(category)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch business rule: "+err.Error())
	}

	var rulesMap map[string]interface{}
	if err := json.Unmarshal([]byte(rule.Rules), &rulesMap); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to parse business rule")
	}

	return response.Success(c, fiber.StatusOK, "Success fetch business rule", model.BusinessRulesResponse{
		Category: rule.Category,
		Rules:    rulesMap,
	})
}

func (h *SettingsHandler) UpsertBusinessRule(c fiber.Ctx) error {
	category := c.Params("category")
	if isTrimmedEmpty(category) {
		return response.Error(c, fiber.StatusBadRequest, "Category is required")
	}

	var rulesMap map[string]interface{}
	if err := c.Bind().JSON(&rulesMap); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	rulesJSON, err := json.Marshal(rulesMap)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid rules format")
	}

	// Get user info from context (if available)
	updatedBy := "system"
	if username, ok := c.Locals("username").(string); ok && username != "" {
		updatedBy = username
	}

	if err := h.settingsSvc.UpsertBusinessRule(category, string(rulesJSON), updatedBy); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to upsert business rule: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Business rule updated successfully", nil)
}
