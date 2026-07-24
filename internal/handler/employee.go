package handler

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/dto"
	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
	"universev2-backend/internal/service"
	"universev2-backend/pkg/response"
)

type EmployeeHandler struct {
	empSvc *service.EmployeeService
}

func NewEmployeeHandler(repo *repository.EmployeeRepo) *EmployeeHandler {
	return &EmployeeHandler{
		empSvc: service.NewEmployeeService(repo),
	}
}

func (h *EmployeeHandler) GetEmployees(c fiber.Ctx) error {
	dept := c.Query("dept")
	status := c.Query("status")
	search := c.Query("q")

	employees, err := h.empSvc.GetEmployees(dept, status, search)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch employees: "+err.Error())
	}

	meta := &response.Meta{
		Page:      1,
		Limit:     50,
		Total:     len(employees),
		TotalPage: 1,
	}
	return response.SuccessWithMeta(c, fiber.StatusOK, "Success fetch employees", employees, meta)
}

func (h *EmployeeHandler) GetEmployeeByNIK(c fiber.Ctx) error {
	nik := c.Params("nik")
	emp, err := h.empSvc.GetEmployeeByNIK(nik)
	if err != nil {
		msg := err.Error()
		if msg == "NIK must be exactly 9 digits" {
			return sendValidationError(c, "nik", msg)
		}
		return response.Error(c, fiber.StatusNotFound, "Employee not found")
	}
	return response.Success(c, fiber.StatusOK, "Success fetch employee", emp)
}

func (h *EmployeeHandler) CreateEmployee(c fiber.Ctx) error {
	var req dto.CreateEmployeeRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	newEmp, err := h.empSvc.CreateEmployee(req)
	if err != nil {
		msg := err.Error()
		switch msg {
		case "NIK must be exactly 9 digits":
			return sendValidationError(c, "nik", msg)
		case "name is required":
			return sendValidationError(c, "name", "Name is required")
		case "employee with this NIK already exists":
			return response.Error(c, fiber.StatusConflict, msg)
		default:
			return response.Error(c, fiber.StatusInternalServerError, msg)
		}
	}
	return response.Success(c, fiber.StatusCreated, "Employee created", newEmp)
}

func (h *EmployeeHandler) UpdateEmployee(c fiber.Ctx) error {
	nik := c.Params("nik")
	var req dto.UpdateEmployeeRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.empSvc.UpdateEmployee(nik, req); err != nil {
		msg := err.Error()
		switch msg {
		case "NIK must be exactly 9 digits":
			return sendValidationError(c, "nik", msg)
		case "name is required":
			return sendValidationError(c, "name", "Name is required")
		case "employee not found":
			return response.Error(c, fiber.StatusNotFound, "Employee not found")
		default:
			return response.Error(c, fiber.StatusInternalServerError, msg)
		}
	}
	return response.Success(c, fiber.StatusOK, "Employee updated", nil)
}

func (h *EmployeeHandler) DeleteEmployee(c fiber.Ctx) error {
	nik := c.Params("nik")
	if err := h.empSvc.DeleteEmployee(nik); err != nil {
		msg := err.Error()
		if msg == "NIK must be exactly 9 digits" {
			return sendValidationError(c, "nik", msg)
		}
		return response.Error(c, fiber.StatusNotFound, "Employee not found")
	}
	return response.Success(c, fiber.StatusOK, "Employee deleted", nil)
}

func (h *EmployeeHandler) GetCompetencies(c fiber.Ctx) error {
	nik := c.Params("nik")
	comps, err := h.empSvc.GetCompetencies(nik)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch competencies: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch competencies", comps)
}

func (h *EmployeeHandler) UpdateCompetencies(c fiber.Ctx) error {
	nik := c.Params("nik")
	var comps []model.Competency
	if err := c.Bind().JSON(&comps); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if err := h.empSvc.UpdateCompetencies(nik, comps); err != nil {
		msg := err.Error()
		if msg == "NIK must be exactly 9 digits" {
			return sendValidationError(c, "nik", msg)
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update competencies: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Competencies updated", nil)
}

func (h *EmployeeHandler) UploadPhoto(c fiber.Ctx) error {
	nik := c.Params("nik")
	if !isValidNIK(nik) {
		return sendValidationError(c, "nik", "NIK must be exactly 9 digits")
	}

	file, err := c.FormFile("photo")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Photo file is required")
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		return response.Error(c, fiber.StatusBadRequest, "Only JPEG and PNG files are allowed")
	}

	// Validate file size (max 5MB)
	const maxFileSize = 5 << 20 // 5MB
	if file.Size > maxFileSize {
		return response.Error(c, fiber.StatusBadRequest, "File size exceeds maximum limit of 5MB")
	}

	// Sanitize filename to prevent path traversal
	safeFilename := strings.ReplaceAll(file.Filename, "..", "")
	safeFilename = strings.ReplaceAll(safeFilename, "/", "")
	safeFilename = strings.ReplaceAll(safeFilename, "\\", "")

	// Save to uploads directory
	photoPath := "uploads/photos/" + nik + "_" + safeFilename
	if err := c.SaveFile(file, photoPath); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to save photo: "+err.Error())
	}

	photoURL := "/" + photoPath
	if err := h.empSvc.UpdatePhoto(nik, photoURL); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update photo URL: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Photo uploaded successfully", fiber.Map{
		"photoUrl": photoURL,
	})
}
