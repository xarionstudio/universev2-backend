package handler

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
	"universev2-backend/pkg/response"
)

type EmployeeHandler struct {
	repo *repository.EmployeeRepo
}

func NewEmployeeHandler(repo *repository.EmployeeRepo) *EmployeeHandler {
	return &EmployeeHandler{repo: repo}
}

func (h *EmployeeHandler) GetEmployees(c fiber.Ctx) error {
	dept := c.Query("dept")
	status := c.Query("status")
	search := c.Query("q")

	employees, err := h.repo.List(dept, status, search)
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
	if !isValidNIK(nik) {
		return sendValidationError(c, "nik", "NIK must be exactly 9 digits")
	}

	emp, err := h.repo.GetByNIK(nik)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Employee not found")
	}
	return response.Success(c, fiber.StatusOK, "Success fetch employee", emp)
}

func (h *EmployeeHandler) CreateEmployee(c fiber.Ctx) error {
	var emp struct {
		NIK     string `json:"nik"`
		Name    string `json:"name"`
		Dept    string `json:"dept"`
		Pos     string `json:"pos"`
		Simper  string `json:"simper"`
		Status  string `json:"status"`
		Company string `json:"company"`
		Equip   string `json:"equip"`
		HP      string `json:"hp"`
	}
	if err := c.Bind().JSON(&emp); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if !isValidNIK(emp.NIK) {
		return sendValidationError(c, "nik", "NIK must be exactly 9 digits")
	}
	if isTrimmedEmpty(emp.Name) {
		return sendValidationError(c, "name", "Name is required")
	}

	existing, _ := h.repo.GetByNIK(emp.NIK)
	if existing != nil {
		return response.Error(c, fiber.StatusConflict, "Employee with this NIK already exists")
	}

	newEmp := &model.Employee{
		NIK: emp.NIK, Name: emp.Name, Dept: emp.Dept, Pos: emp.Pos,
		Simper: emp.Simper, Status: emp.Status, Company: emp.Company,
		Equip: emp.Equip, HP: emp.HP,
	}
	if err := h.repo.Create(newEmp); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create employee: "+err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Employee created", newEmp)
}

func (h *EmployeeHandler) UpdateEmployee(c fiber.Ctx) error {
	nik := c.Params("nik")
	if !isValidNIK(nik) {
		return sendValidationError(c, "nik", "NIK must be exactly 9 digits")
	}

	existing, err := h.repo.GetByNIK(nik)
	if err != nil || existing == nil {
		return response.Error(c, fiber.StatusNotFound, "Employee not found")
	}

	var emp model.Employee
	if err := c.Bind().JSON(&emp); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(emp.Name) {
		return sendValidationError(c, "name", "Name is required")
	}

	if err := h.repo.Update(nik, &emp); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update employee: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Employee updated", nil)
}

func (h *EmployeeHandler) DeleteEmployee(c fiber.Ctx) error {
	nik := c.Params("nik")
	if !isValidNIK(nik) {
		return sendValidationError(c, "nik", "NIK must be exactly 9 digits")
	}

	existing, err := h.repo.GetByNIK(nik)
	if err != nil || existing == nil {
		return response.Error(c, fiber.StatusNotFound, "Employee not found")
	}

	if err := h.repo.Delete(nik); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete employee: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Employee deleted", nil)
}

func (h *EmployeeHandler) GetCompetencies(c fiber.Ctx) error {
	nik := c.Params("nik")
	comps, err := h.repo.GetCompetencies(nik)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch competencies: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch competencies", comps)
}

func (h *EmployeeHandler) UpdateCompetencies(c fiber.Ctx) error {
	nik := c.Params("nik")
	if !isValidNIK(nik) {
		return sendValidationError(c, "nik", "NIK must be exactly 9 digits")
	}

	var comps []model.Competency
	if err := c.Bind().JSON(&comps); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if err := h.repo.UpdateCompetencies(nik, comps); err != nil {
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
	if err := h.repo.UpdatePhoto(nik, photoURL); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update photo URL: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Photo uploaded successfully", fiber.Map{
		"photoUrl": photoURL,
	})
}
