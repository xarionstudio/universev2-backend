package handler

import (
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
	var emp model.Employee
	if err := c.Bind().JSON(&emp); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if err := h.repo.Update(nik, &emp); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update employee: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Employee updated", nil)
}

func (h *EmployeeHandler) DeleteEmployee(c fiber.Ctx) error {
	nik := c.Params("nik")
	if err := h.repo.Delete(nik); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete employee: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Employee deleted", nil)
}
