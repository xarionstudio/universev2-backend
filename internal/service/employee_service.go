package service

import (
	"fmt"

	"universev2-backend/internal/dto"
	"universev2-backend/internal/export"
	"universev2-backend/internal/model"
	internalpkg "universev2-backend/internal/pkg"
	"universev2-backend/internal/repository"
	"universev2-backend/pkg/filter"
	"universev2-backend/pkg/pagination"
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

	if req.Status != "" && req.Status != "aktif" && req.Status != "cuti" && req.Status != "nonaktif" {
		return nil, fmt.Errorf("invalid status: must be aktif, cuti, or nonaktif")
	}

	newEmp := &model.Employee{
		NIK: req.NIK, Name: req.Name, Dept: req.Dept, Pos: req.Pos,
		Simper: req.Simper, SimperExp: req.SimperExp, Status: req.Status,
		Company: req.Company, Equip: req.Equip, Join: req.Join, Exp: req.Exp,
		License: req.License, MCU: req.MCU, Medis: req.Medis, Blood: req.Blood,
		BPJS: req.BPJS, Mess: req.Mess, Kamar: req.Kamar, HP: req.HP,
		Emergency: req.Emergency, Foto: req.Foto,
	}
	if err := s.repo.Create(newEmp); err != nil {
		return nil, fmt.Errorf("failed to create employee: %w", err)
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

	if req.Status != "" && req.Status != "aktif" && req.Status != "cuti" && req.Status != "nonaktif" {
		return fmt.Errorf("invalid status: must be aktif, cuti, or nonaktif")
	}

	emp := &model.Employee{
		Name: req.Name, Dept: req.Dept, Pos: req.Pos,
		Simper: req.Simper, SimperExp: req.SimperExp, Status: req.Status,
		Company: req.Company, Equip: req.Equip, Join: req.Join, Exp: req.Exp,
		License: req.License, MCU: req.MCU, Medis: req.Medis, Blood: req.Blood,
		BPJS: req.BPJS, Mess: req.Mess, Kamar: req.Kamar, HP: req.HP,
		Emergency: req.Emergency, Foto: req.Foto,
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

// GetEmployeesPaginated returns paginated employees with advanced filter support.
func (s *EmployeeService) GetEmployeesPaginated(f filter.Params, p pagination.Params) ([]model.Employee, int64, error) {
	return s.repo.ListPaginated(f, p)
}

// ImportEmployeesFromExcel reads an xlsx file and bulk-creates employees.
// Returns counts of imported, skipped records and any row-level errors.
func (s *EmployeeService) ImportEmployeesFromExcel(data []byte) (imported, skipped int, rowErrors []string, err error) {
	employees, parseErr := export.ParseEmployeeExcel(data)
	if parseErr != nil {
		return 0, 0, nil, fmt.Errorf("failed to parse excel: %w", parseErr)
	}

	for _, emp := range employees {
		if !internalpkg.IsValidNIK(emp.NIK) {
			skipped++
			rowErrors = append(rowErrors, fmt.Sprintf("NIK %q: invalid (must be 9 digits)", emp.NIK))
			continue
		}
		if internalpkg.IsTrimmedEmpty(emp.Name) {
			skipped++
			rowErrors = append(rowErrors, fmt.Sprintf("NIK %q: name is required", emp.NIK))
			continue
		}

		existing, _ := s.repo.GetByNIK(emp.NIK)
		if existing != nil {
			skipped++
			rowErrors = append(rowErrors, fmt.Sprintf("NIK %q: already exists, skipped", emp.NIK))
			continue
		}

		e := emp // avoid loop-var capture
		if createErr := s.repo.Create(&e); createErr != nil {
			skipped++
			rowErrors = append(rowErrors, fmt.Sprintf("NIK %q: %v", emp.NIK, createErr))
			continue
		}
		imported++
	}
	return imported, skipped, rowErrors, nil
}

// ExportEmployeesToExcel generates an xlsx for employees matching the given filter.
func (s *EmployeeService) ExportEmployeesToExcel(f filter.Params) ([]byte, error) {
	// Fetch all (no page limit) with a very large perPage
	p := pagination.Params{Page: 1, PerPage: pagination.MaxPerPage}
	employees, _, err := s.repo.ListPaginated(f, p)
	if err != nil {
		return nil, err
	}
	return export.GenerateEmployeeExcel(employees)
}
