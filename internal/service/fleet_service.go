package service

import (
	"fmt"
	"strconv"
	"time"

	"universev/internal/dto"
	"universev/internal/export"
	"universev/internal/model"
	internalpkg "universev/internal/pkg"
	"universev/internal/repository"
)

// FleetService handles fleet and unit business logic
type FleetService struct {
	repo *repository.FleetRepo
}

// NewFleetService creates a new FleetService
func NewFleetService(repo *repository.FleetRepo) *FleetService {
	return &FleetService{repo: repo}
}

// GetFleetSettings returns all fleet settings
func (s *FleetService) GetFleetSettings() ([]model.FleetSetting, error) {
	return s.repo.GetFleetSettings()
}

// CreateFleetSetting creates a new fleet setting
func (s *FleetService) CreateFleetSetting(req dto.CreateFleetSettingRequest) (*model.FleetSetting, error) {
	if internalpkg.IsTrimmedEmpty(req.Digger) {
		return nil, fmt.Errorf("digger code is required")
	}
	f := &model.FleetSetting{
		Digger: req.Digger,
		Loc:    req.Loc,
		Bus:    req.Bus,
		Units:  req.Units,
	}
	if err := s.repo.CreateFleetSetting(f); err != nil {
		return nil, fmt.Errorf("failed to create fleet setting: %w", err)
	}
	return f, nil
}

// UpdateFleetSetting updates an existing fleet setting
func (s *FleetService) UpdateFleetSetting(id string, req dto.UpdateFleetSettingRequest) error {
	if internalpkg.IsTrimmedEmpty(id) {
		return fmt.Errorf("ID is required")
	}
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ID")
	}
	if internalpkg.IsTrimmedEmpty(req.Digger) {
		return fmt.Errorf("digger code is required")
	}
	f := &model.FleetSetting{
		Digger: req.Digger,
		Loc:    req.Loc,
		Bus:    req.Bus,
		Units:  req.Units,
	}
	return s.repo.UpdateFleetSetting(uint(uid), f)
}

// DeleteFleetSetting deletes a fleet setting
func (s *FleetService) DeleteFleetSetting(id string) error {
	if internalpkg.IsTrimmedEmpty(id) {
		return fmt.Errorf("ID is required")
	}
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ID")
	}
	return s.repo.DeleteFleetSetting(uint(uid))
}

// GetAllocations returns fleet allocations
func (s *FleetService) GetAllocations(date, shift string) (model.FleetAllocResponse, error) {
	allocs, err := s.repo.GetAllocations(date, shift)
	if err != nil {
		return nil, err
	}
	result := make(model.FleetAllocResponse)
	for _, a := range allocs {
		if result[a.Date] == nil {
			result[a.Date] = make(map[string]map[string]string)
		}
		result[a.Date][a.Shift] = a.Units
	}
	return result, nil
}

// AutoAllocate performs auto allocation
func (s *FleetService) AutoAllocate(req dto.AutoAllocateRequest) error {
	if internalpkg.IsTrimmedEmpty(req.Date) {
		return fmt.Errorf("date is required")
	}
	if internalpkg.IsTrimmedEmpty(req.Shift) {
		return fmt.Errorf("shift is required")
	}
	return s.repo.AutoAllocate(req.Date, req.Shift)
}

// GetUnitStatuses returns all unit statuses
func (s *FleetService) GetUnitStatuses() ([]model.Unit, error) {
	return s.repo.GetUnitStatuses()
}

// UpdateUnitStatus updates a unit's status
func (s *FleetService) UpdateUnitStatus(code string, req dto.UpdateUnitStatusRequest) error {
	if internalpkg.IsTrimmedEmpty(code) {
		return fmt.Errorf("unit code is required")
	}
	if internalpkg.IsTrimmedEmpty(req.Status) {
		return fmt.Errorf("status is required")
	}
	if req.Status != string(model.UnitStatusReady) && req.Status != string(model.UnitStatusBreakdown) && req.Status != string(model.UnitStatusStandby) {
		return fmt.Errorf("status must be one of: ready, breakdown, standby")
	}

	if err := s.repo.UpdateUnitStatus(code, model.UnitStatus(req.Status), req.Note); err != nil {
		return err
	}
	nowStr := time.Now().Format("02 Jan 15:04")
	_ = s.repo.AddUnitHistory(code, nowStr, req.Status, req.Note, req.Status)
	return nil
}

// ReportUnitBreakdown reports a unit breakdown
func (s *FleetService) ReportUnitBreakdown(code string, req dto.ReportBreakdownRequest) error {
	if internalpkg.IsTrimmedEmpty(code) {
		return fmt.Errorf("unit code is required")
	}
	if internalpkg.IsTrimmedEmpty(req.Reason) {
		return fmt.Errorf("breakdown reason is required")
	}
	if err := s.repo.UpdateUnitStatus(code, model.UnitStatusBreakdown, req.Reason); err != nil {
		return err
	}
	nowStr := time.Now().Format("02 Jan 15:04")
	_ = s.repo.AddUnitHistory(code, nowStr, "Breakdown", req.Reason, "breakdown")
	return nil
}

// GetUnitHistory returns unit history
func (s *FleetService) GetUnitHistory(code string) ([]model.UnitHist, error) {
	if internalpkg.IsTrimmedEmpty(code) {
		return nil, fmt.Errorf("unit code is required")
	}
	return s.repo.GetUnitHistory(code)
}

// GetUnitDB returns all unit DB entries
func (s *FleetService) GetUnitDB() ([]model.UnitDb, error) {
	return s.repo.GetUnitDB()
}

// CreateUnitDB creates a new unit DB entry
func (s *FleetService) CreateUnitDB(req dto.CreateUnitDBRequest) (*model.UnitDb, error) {
	if internalpkg.IsTrimmedEmpty(req.Code) {
		return nil, fmt.Errorf("unit code is required")
	}
	u := &model.UnitDb{
		Code:      req.Code,
		EGI:       req.EGI,
		Product:   req.Product,
		Cls:       req.Cls,
		Category:  req.Category,
		Area:      req.Area,
		Active:    req.Active,
		Standby:   req.Standby,
		Breakdown: req.Breakdown,
		Loc:       req.Loc,
		Upd:       req.Upd,
		By:        req.By,
	}
	if err := s.repo.CreateUnitDB(u); err != nil {
		return nil, fmt.Errorf("failed to create unit DB: %w", err)
	}
	return u, nil
}

// UpdateUnitDB updates a unit DB entry
func (s *FleetService) UpdateUnitDB(req dto.UpdateUnitDBRequest) error {
	if internalpkg.IsTrimmedEmpty(req.Code) {
		return fmt.Errorf("unit code is required")
	}
	u := &model.UnitDb{
		Code:      req.Code,
		EGI:       req.EGI,
		Product:   req.Product,
		Cls:       req.Cls,
		Category:  req.Category,
		Area:      req.Area,
		Active:    req.Active,
		Standby:   req.Standby,
		Breakdown: req.Breakdown,
		Loc:       req.Loc,
		Upd:       req.Upd,
		By:        req.By,
	}
	return s.repo.UpdateUnitDB(u)
}

// DeleteUnitDB deletes a unit DB entry by unit code
func (s *FleetService) DeleteUnitDB(code string) error {
	if internalpkg.IsTrimmedEmpty(code) {
		return fmt.Errorf("unit code is required")
	}
	return s.repo.DeleteUnitDBByCode(code)
}

// ImportUnitDBFromExcel parses Excel file and saves unit DB entries
func (s *FleetService) ImportUnitDBFromExcel(data []byte) (imported int, skipped int, rowErrors []string, err error) {
	units, err := export.ParseUnitDBExcel(data)
	if err != nil {
		return 0, 0, nil, err
	}
	if len(units) == 0 {
		return 0, 0, nil, fmt.Errorf("no valid unit records found in file")
	}
	return s.repo.BulkCreateUnitDB(units)
}

// SaveAllocation saves manual allocation (assign/release) for a date+shift
func (s *FleetService) SaveAllocation(req dto.SaveAllocationRequest) error {
	if internalpkg.IsTrimmedEmpty(req.Date) {
		return fmt.Errorf("date is required")
	}
	if internalpkg.IsTrimmedEmpty(req.Shift) {
		return fmt.Errorf("shift is required")
	}
	return s.repo.SaveAllocation(req.Date, req.Shift, req.Units)
}
