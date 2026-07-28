package service

import (
	"fmt"
	"strings"

	"universev2-backend/internal/export"
	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
	"universev2-backend/pkg/pagination"
)

// MasterService handles master data business logic
type MasterService struct {
	repo *repository.MasterRepo
}

// NewMasterService creates a new MasterService
func NewMasterService(repo *repository.MasterRepo) *MasterService {
	return &MasterService{repo: repo}
}

// MasterListResponse represents paginated master data response
type MasterListResponse struct {
	Entries    interface{} `json:"entries"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"perPage"`
	TotalPages int         `json:"totalPages"`
	Category   string      `json:"category"`
}

// GetMasterByCategory returns master entries with pagination and search
func (s *MasterService) GetMasterByCategory(category string, page, perPage int, search string) (*MasterListResponse, error) {
	raw, err := s.repo.GetByCategory(category)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return &MasterListResponse{
			Entries:    []interface{}{},
			Total:      0,
			Page:       page,
			PerPage:    perPage,
			TotalPages: 0,
			Category:   category,
		}, nil
	}

	// Count total
	total := countEntries(raw)

	// Apply search filter
	if search != "" {
		raw = filterEntries(raw, search)
		total = countEntries(raw)
	}

	// Pagination
	if perPage <= 0 {
		perPage = 10
	}
	if page <= 0 {
		page = 1
	}

	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	start := (page - 1) * perPage
	end := start + perPage
	if start > total {
		start = 0
		end = perPage
	}
	if end > total {
		end = total
	}

	// Slice the entries
	sliced := sliceEntries(raw, start, end)

	return &MasterListResponse{
		Entries:    sliced,
		Total:      int64(total),
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		Category:   category,
	}, nil
}

// BulkCreate creates multiple master entries at once
func (s *MasterService) BulkCreate(category string, entries []interface{}) (interface{}, error) {
	var created []interface{}
	for _, entry := range entries {
		if err := s.repo.Create(category, entry); err != nil {
			return nil, err
		}
		created = append(created, entry)
	}
	return created, nil
}

// UpdateEntry updates a single master entry
func (s *MasterService) UpdateEntry(category string, code string, updates map[string]interface{}) error {
	return s.repo.Update(category, code, updates)
}

// BulkDelete deletes multiple master entries by codes
func (s *MasterService) BulkDelete(category string, codes []string) error {
	for _, code := range codes {
		if err := s.repo.Delete(category, code); err != nil {
			return err
		}
	}
	return nil
}

// ImportFromExcel bulk-imports master entries from xlsx bytes.
func (s *MasterService) ImportFromExcel(category string, data []byte) (int, error) {
	entries, err := export.ParseMasterExcel(data)
	if err != nil {
		return 0, fmt.Errorf("failed to parse excel: %w", err)
	}

	imported := 0
	for _, entry := range entries {
		if entry.ID == "" {
			entry.ID = fmt.Sprintf("%s-%s", category, strings.ToLower(strings.ReplaceAll(entry.Name, " ", "-")))
		}
		// Map to per-category struct
		model := mapExcelToModel(category, entry)
		if model == nil {
			continue
		}
		if err := s.repo.Create(category, model); err != nil {
			continue
		}
		imported++
	}
	return imported, nil
}

// ExportToExcel exports all master entries for a category as xlsx bytes.
func (s *MasterService) ExportToExcel(category string) ([]byte, error) {
	raw, err := s.repo.GetByCategory(category)
	if err != nil {
		return nil, err
	}
	return export.GenerateMasterExcel(category, raw)
}

// ParsePaginationParams extracts page and perPage from query params.
func ParsePaginationParams(pageStr, perPageStr string) (int, int) {
	p := pagination.Parse(pageStr, perPageStr)
	return p.Page, p.PerPage
}

// ── Internal helpers ──────────────────────────────────────────────────────────

func countEntries(raw interface{}) int {
	if raw == nil {
		return 0
	}
	switch v := raw.(type) {
	case []model.MasterEGIType:
		return len(v)
	case []model.MasterProduct:
		return len(v)
	case []model.MasterEqClass:
		return len(v)
	case []model.MasterArea:
		return len(v)
	case []model.MasterTempudo:
		return len(v)
	case []model.MasterBus:
		return len(v)
	case []model.MasterLocationEx:
		return len(v)
	case []model.MasterMess:
		return len(v)
	case []model.MasterRunningText:
		return len(v)
	}
	return 0
}

func filterEntries(raw interface{}, search string) interface{} {
	needle := strings.ToLower(search)
	switch v := raw.(type) {
	case []model.MasterEGIType:
		var filtered []model.MasterEGIType
		for _, e := range v {
			if containsIgnoreCase(e.Name, needle) || containsIgnoreCase(e.Code, needle) {
				filtered = append(filtered, e)
			}
		}
		return filtered
	case []model.MasterProduct:
		var filtered []model.MasterProduct
		for _, e := range v {
			if containsIgnoreCase(e.Name, needle) || containsIgnoreCase(e.Code, needle) {
				filtered = append(filtered, e)
			}
		}
		return filtered
	case []model.MasterEqClass:
		var filtered []model.MasterEqClass
		for _, e := range v {
			if containsIgnoreCase(e.Name, needle) || containsIgnoreCase(e.Code, needle) || containsIgnoreCase(e.Description, needle) {
				filtered = append(filtered, e)
			}
		}
		return filtered
	case []model.MasterArea:
		var filtered []model.MasterArea
		for _, e := range v {
			if containsIgnoreCase(e.Name, needle) || containsIgnoreCase(e.Code, needle) || containsIgnoreCase(e.Category, needle) {
				filtered = append(filtered, e)
			}
		}
		return filtered
	case []model.MasterTempudo:
		var filtered []model.MasterTempudo
		for _, e := range v {
			if containsIgnoreCase(e.Name, needle) || containsIgnoreCase(e.Code, needle) || containsIgnoreCase(e.Location, needle) || containsIgnoreCase(e.PickupType, needle) {
				filtered = append(filtered, e)
			}
		}
		return filtered
	case []model.MasterBus:
		var filtered []model.MasterBus
		for _, e := range v {
			if containsIgnoreCase(e.Name, needle) || containsIgnoreCase(e.Code, needle) || containsIgnoreCase(e.EGIType, needle) || containsIgnoreCase(e.DepartureTime, needle) {
				filtered = append(filtered, e)
			}
		}
		return filtered
	case []model.MasterLocationEx:
		var filtered []model.MasterLocationEx
		for _, e := range v {
			if containsIgnoreCase(e.Name, needle) || containsIgnoreCase(e.Code, needle) || containsIgnoreCase(e.BusCode, needle) || containsIgnoreCase(e.TempudoCode, needle) {
				filtered = append(filtered, e)
			}
		}
		return filtered
	case []model.MasterMess:
		var filtered []model.MasterMess
		for _, e := range v {
			if containsIgnoreCase(e.Name, needle) || containsIgnoreCase(e.Code, needle) || containsIgnoreCase(e.Block, needle) {
				filtered = append(filtered, e)
			}
		}
		return filtered
	case []model.MasterRunningText:
		var filtered []model.MasterRunningText
		for _, e := range v {
			if containsIgnoreCase(e.Name, needle) || containsIgnoreCase(e.Code, needle) || containsIgnoreCase(e.TargetDisplay, needle) || containsIgnoreCase(e.TextColor, needle) {
				filtered = append(filtered, e)
			}
		}
		return filtered
	}
	return raw
}

func sliceEntries(raw interface{}, start, end int) interface{} {
	switch v := raw.(type) {
	case []model.MasterEGIType:
		if start >= len(v) {
			return []model.MasterEGIType{}
		}
		if end > len(v) {
			end = len(v)
		}
		return v[start:end]
	case []model.MasterProduct:
		if start >= len(v) {
			return []model.MasterProduct{}
		}
		if end > len(v) {
			end = len(v)
		}
		return v[start:end]
	case []model.MasterEqClass:
		if start >= len(v) {
			return []model.MasterEqClass{}
		}
		if end > len(v) {
			end = len(v)
		}
		return v[start:end]
	case []model.MasterArea:
		if start >= len(v) {
			return []model.MasterArea{}
		}
		if end > len(v) {
			end = len(v)
		}
		return v[start:end]
	case []model.MasterTempudo:
		if start >= len(v) {
			return []model.MasterTempudo{}
		}
		if end > len(v) {
			end = len(v)
		}
		return v[start:end]
	case []model.MasterBus:
		if start >= len(v) {
			return []model.MasterBus{}
		}
		if end > len(v) {
			end = len(v)
		}
		return v[start:end]
	case []model.MasterLocationEx:
		if start >= len(v) {
			return []model.MasterLocationEx{}
		}
		if end > len(v) {
			end = len(v)
		}
		return v[start:end]
	case []model.MasterMess:
		if start >= len(v) {
			return []model.MasterMess{}
		}
		if end > len(v) {
			end = len(v)
		}
		return v[start:end]
	case []model.MasterRunningText:
		if start >= len(v) {
			return []model.MasterRunningText{}
		}
		if end > len(v) {
			end = len(v)
		}
		return v[start:end]
	}
	return raw
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// mapExcelToModel creates a per-category struct from an Excel parsed entry
func mapExcelToModel(cat string, entry export.MasterExcelEntry) interface{} {
	switch cat {
	case "egi":
		return &model.MasterEGIType{Code: entry.ID, Name: entry.Name, Active: entry.Active}
	case "product":
		return &model.MasterProduct{Code: entry.ID, Name: entry.Name, Active: entry.Active}
	case "eqclass":
		return &model.MasterEqClass{Code: entry.ID, Name: entry.Name, Description: entry.A, Active: entry.Active}
	case "area":
		return &model.MasterArea{Code: entry.ID, Name: entry.Name, Category: entry.A, Active: entry.Active}
	case "tempudo":
		return &model.MasterTempudo{Code: entry.ID, Name: entry.Name, Location: entry.A, PickupType: entry.B, Active: entry.Active}
	case "bus":
		return &model.MasterBus{Code: entry.ID, Name: entry.Name, EGIType: entry.A, DepartureTime: entry.B, Active: entry.Active}
	case "lokasiex":
		return &model.MasterLocationEx{Code: entry.ID, Name: entry.Name, BusCode: entry.A, TempudoCode: entry.B, Active: entry.Active}
	case "mess":
		return &model.MasterMess{Code: entry.ID, Name: entry.Name, Block: entry.A, Active: entry.Active}
	case "runtext":
		return &model.MasterRunningText{Code: entry.ID, Name: entry.Name, TargetDisplay: entry.A, TextColor: entry.B, Active: entry.Active}
	}
	return nil
}
