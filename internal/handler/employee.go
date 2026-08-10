package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"

	"universev/internal/dto"
	"universev/internal/model"
	"universev/internal/repository"
	"universev/internal/service"
	"universev/pkg/filter"
	"universev/pkg/pagination"
	"universev/pkg/response"
)

type EmployeeHandler struct {
	empSvc    *service.EmployeeService
	uploadDir string
}

func NewEmployeeHandler(repo *repository.EmployeeRepo, uploadDir string) *EmployeeHandler {
	return &EmployeeHandler{
		empSvc:    service.NewEmployeeService(repo),
		uploadDir: uploadDir,
	}
}

// GetEmployees godoc
// GET /api/employees
// Query params: page, perPage, search, status, dept, date_from, date_to, logic
func (h *EmployeeHandler) GetEmployees(c fiber.Ctx) error {
	f := filter.ParseFromCtx(c)
	p := pagination.Parse(c.Query("page"), c.Query("perPage"))

	employees, total, err := h.empSvc.GetEmployeesPaginated(f, p)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch employees: "+err.Error())
	}

	meta := pagination.BuildMeta(p, total)
	return response.SuccessPaged(c, fiber.StatusOK, "Success fetch employees", response.PagedData{
		Items:      employees,
		Pagination: meta,
	})
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

	// Support both shapes: { competencies: [...] } and direct array
	var body struct {
		Competencies []model.Competency `json:"competencies"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		// Try direct array binding as fallback
		var comps []model.Competency
		if err2 := c.Bind().JSON(&comps); err2 != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
		}
		body.Competencies = comps
	}

	if err := h.empSvc.UpdateCompetencies(nik, body.Competencies); err != nil {
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
	photoDir := filepath.Join(h.uploadDir, "photos")
	if err := os.MkdirAll(photoDir, 0755); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create upload directory: "+err.Error())
	}
	photoPath := filepath.Join(photoDir, nik+"_"+safeFilename)
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

// ImportEmployees godoc
// POST /api/employees/import
// Content-Type: multipart/form-data
// Field: file (.xlsx)
func (h *EmployeeHandler) ImportEmployees(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Excel file is required (field: 'file')")
	}

	// Validate extension
	fname := strings.ToLower(file.Filename)
	if !strings.HasSuffix(fname, ".xlsx") && !strings.HasSuffix(fname, ".xls") {
		return response.Error(c, fiber.StatusBadRequest, "Only .xlsx or .xls files are accepted")
	}

	// Validate file size (max 10MB)
	const maxSize = 10 << 20
	if file.Size > maxSize {
		return response.Error(c, fiber.StatusBadRequest, "File size exceeds 10MB limit")
	}

	// Read file bytes
	f, err := file.Open()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to open uploaded file")
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to read file")
	}

	imported, skipped, rowErrors, err := h.empSvc.ImportEmployeesFromExcel(data)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Import failed: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK, fmt.Sprintf("Import completed: %d imported, %d skipped", imported, skipped), fiber.Map{
		"imported": imported,
		"skipped":  skipped,
		"errors":   rowErrors,
	})
}

// ExportEmployees godoc
// GET /api/employees/export
// Accepts same filter params as GetEmployees; returns an xlsx download.
func (h *EmployeeHandler) ExportEmployees(c fiber.Ctx) error {
	f := filter.ParseFromCtx(c)

	xlsxData, err := h.empSvc.ExportEmployeesToExcel(f)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to generate export: "+err.Error())
	}

	fileName := fmt.Sprintf("employees_export_%s.xlsx", time.Now().Format("20060102"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	return c.Send(xlsxData)
}
