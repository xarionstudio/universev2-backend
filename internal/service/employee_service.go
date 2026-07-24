package service

import (
	"fmt"

	"universev2-backend/internal/dto"
	"universev2-backend/internal/model"
	internalpkg "universev2-backend/internal/pkg"
	"universev2-backend/internal/repository"
)

// EmployeeService handles employee business logic
type EmployeeService struct {
	repo *repository.EmployeeRepo
}

// NewEmployeeService creates a new EmployeeService
func NewEmployeeService(repo *repository.EmployeeRepo) *EmployeeService {
	return &EmployeeService{repo: repo}
}

// GetEmployees returns a list of employees with optional filters
func (s *EmployeeService) GetEmployees(dept, status, search string) ([]model.Employee, error) {
	return s.repo.List(dept, status, search)
}

// GetEmployeeByNIK returns a single employee by NIK
func (s *EmployeeService) GetEmployeeByNIK(nik string) (*model.Employee, error) {
	if !internalpkg.IsValidNIK(nik) {
		return nil, fmt.Errorf("NIK must be exactly 9 digits")
	}
	emp, err := s.repo.GetByNIK(nik)
	if err != nil {
		return nil, fmt.Errorf("employee not found")
	}
	return emp, nil
}

// CreateEmployee creates a new employee
func (s *EmployeeService) CreateEmployee(req dto.CreateEmployeeRequest) (*model.Employee, error) {
	if !internalpkg.IsValidNIK(req.NIK) {
		return nil, fmt.Errorf("NIK must be exactly 9 digits")
	}
	if internalpkg.IsTrimmedEmpty(req.Name) {
		return nil, fmt.Errorf("name is required")
	}

	existing, _ := s.repo.GetByNIK(req.NIK)
	if existing != nil {
		return nil, fmt.Errorf("employee with this NIK already exists")
	}

	newEmp := &model.Employee{
		NIK: req.NIK, Name: req.Name, Dept: req.Dept, Pos: req.Pos,
		Simper: req.Simper, Status: req.Status, Company: req.Company,
		Equip: req.Equip, HP: req.HP,
	}
	if err := s.repo.Create(newEmp); err != nil {
		return nil, fmt.Errorf("failed to create employee: " + err.Error())
	}
	return newEmp, nil
}

// UpdateEmployee updates an existing employee
func (s *EmployeeService) UpdateEmployee(nik string, req dto.UpdateEmployeeRequest) error {
	if !internalpkg.IsValidNIK(nik) {
		return fmt.Errorf("NIK must be exactly 9 digits")
	}

	existing, err := s.repo.GetByNIK(nik)
	if err != nil || existing == nil {
		return fmt.Errorf("employee not found")
	}

	if internalpkg.IsTrimmedEmpty(req.Name) {
		return fmt.Errorf("name is required")
	}

	emp := &model.Employee{
		Name: req.Name, Dept: req.Dept, Pos: req.Pos,
		Simper: req.Simper, Status: req.Status, Company: req.Company,
		Equip: req.Equip, HP: req.HP,
	}
	return s.repo.Update(nik, emp)
}

// DeleteEmployee deletes an employee by NIK
func (s *EmployeeService) DeleteEmployee(nik string) error {
	if !internalpkg.IsValidNIK(nik) {
		return fmt.Errorf("NIK must be exactly 9 digits")
	}

	existing, err := s.repo.GetByNIK(nik)
	if err != nil || existing == nil {
		return fmt.Errorf("employee not found")
	}

	return s.repo.Delete(nik)
}

// GetCompetencies returns competencies for an employee
func (s *EmployeeService) GetCompetencies(nik string) ([]model.Competency, error) {
	return s.repo.GetCompetencies(nik)
}

// UpdateCompetencies updates competencies for an employee
func (s *EmployeeService) UpdateCompetencies(nik string, comps []model.Competency) error {
	if !internalpkg.IsValidNIK(nik) {
		return fmt.Errorf("NIK must be exactly 9 digits")
	}
	return s.repo.UpdateCompetencies(nik, comps)
}

// UpdatePhoto updates the photo URL for an employee
func (s *EmployeeService) UpdatePhoto(nik string, photoURL string) error {
	return s.repo.UpdatePhoto(nik, photoURL)
}
