package handler

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/dto"
	"universev2-backend/internal/repository"
	"universev2-backend/internal/service"
	"universev2-backend/pkg/response"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userRepo *repository.UserRepo) *UserHandler {
	return &UserHandler{
		userSvc: service.NewUserService(userRepo, nil),
	}
}

func (h *UserHandler) GetUsers(c fiber.Ctx) error {
	users, err := h.userSvc.GetUsers()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch users: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch users", users)
}

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := h.userSvc.CreateUser(req)
	if err != nil {
		msg := err.Error()
		switch msg {
		case "name is required":
			return sendValidationError(c, "name", "Name is required")
		case "invalid email format":
			return sendValidationError(c, "email", "Invalid email format")
		case "at least one role is required":
			return sendValidationError(c, "roles", "At least one role is required")
		case "email is already in use":
			return response.Error(c, fiber.StatusConflict, "Email is already in use")
		default:
			return response.Error(c, fiber.StatusInternalServerError, msg)
		}
	}
	return response.Success(c, fiber.StatusCreated, "User created successfully", user)
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.userSvc.UpdateUser(id, req); err != nil {
		msg := err.Error()
		switch msg {
		case "user not found":
			return response.Error(c, fiber.StatusNotFound, "User not found")
		case "name is required":
			return sendValidationError(c, "name", "Name is required")
		case "invalid email format":
			return sendValidationError(c, "email", "Invalid email format")
		case "at least one role is required":
			return sendValidationError(c, "roles", "At least one role is required")
		default:
			return response.Error(c, fiber.StatusInternalServerError, msg)
		}
	}
	return response.Success(c, fiber.StatusOK, "User updated successfully", nil)
}

func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.userSvc.DeleteUser(id); err != nil {
		return response.Error(c, fiber.StatusNotFound, "User not found")
	}
	return response.Success(c, fiber.StatusOK, "User deleted successfully", nil)
}

// ImportUsers godoc
// POST /api/users/import
// Content-Type: multipart/form-data
// Field: file (.xlsx)
func (h *UserHandler) ImportUsers(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Excel file is required (field: 'file')")
	}

	fname := strings.ToLower(file.Filename)
	if !strings.HasSuffix(fname, ".xlsx") && !strings.HasSuffix(fname, ".xls") {
		return response.Error(c, fiber.StatusBadRequest, "Only .xlsx or .xls files are accepted")
	}

	const maxSize = 10 << 20 // 10MB
	if file.Size > maxSize {
		return response.Error(c, fiber.StatusBadRequest, "File size exceeds 10MB limit")
	}

	f, err := file.Open()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to open uploaded file")
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to read file")
	}

	imported, skipped, err := h.userSvc.ImportUsersFromExcel(data)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Import failed: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK, fmt.Sprintf("Import completed: %d users imported, %d skipped/duplicates", imported, skipped), fiber.Map{
		"imported": imported,
		"skipped":  skipped,
	})
}

// ToggleUserStatus godoc
// PATCH /api/users/:id/status
// Body: { "active": true/false }
func (h *UserHandler) ToggleUserStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Active bool `json:"active"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.userSvc.ToggleUserStatus(id, req.Active); err != nil {
		if err.Error() == "user not found" {
			return response.Error(c, fiber.StatusNotFound, "User not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to toggle user status: "+err.Error())
	}

	statusMsg := "deactivated"
	if req.Active {
		statusMsg = "activated"
	}
	return response.Success(c, fiber.StatusOK, fmt.Sprintf("User status %s successfully", statusMsg), nil)
}

// ExportUsers godoc
// GET /api/users/export
// Download CSV file of users
func (h *UserHandler) ExportUsers(c fiber.Ctx) error {
	csvData, err := h.userSvc.ExportUsersCSV()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to generate CSV export: "+err.Error())
	}

	fileName := fmt.Sprintf("users_%s.csv", time.Now().Format("20060102"))
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	return c.Send(csvData)
}


