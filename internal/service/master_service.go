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
	Entries    []model.MdEntry `json:"entries"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"perPage"`
	TotalPages int             `json:"totalPages"`
}

// GetMasterByCategory returns master entries with pagination and search
func (s *MasterService) GetMasterByCategory(category string, page, perPage int, search string) (*MasterListResponse, error) {
	entries, err := s.repo.GetByCategory(category)
	if err != nil {
		return nil, err
	}

	// Filter by search if provided
	if search != "" {
		filtered := make([]model.MdEntry, 0)
		for _, entry := range entries {
			if containsIgnoreCase(entry.Name, search) ||
				containsIgnoreCase(entry.FieldA, search) ||
				containsIgnoreCase(entry.FieldB, search) {
				filtered = append(filtered, entry)
			}
		}
		entries = filtered
	}

	// Sort by name
	sortByName(entries)

	total := len(entries)

	// Pagination
	if perPage <= 0 {
		perPage = 10
	}
	if page <= 0 {
		page = 1
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

	paginatedEntries := entries[start:end]
	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	return &MasterListResponse{
		Entries:    paginatedEntries,
		Total:      int64(total),
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// BulkCreate creates multiple master entries at once
func (s *MasterService) BulkCreate(category string, entries []model.MdEntry) ([]model.MdEntry, error) {
	created := make([]model.MdEntry, 0, len(entries))
	for _, entry := range entries {
		entry.Cat = model.MdCat(category)
		if err := s.repo.Create(&entry); err != nil {
			return nil, err
		}
		created = append(created, entry)
	}
	return created, nil
}

// BulkDelete deletes multiple master entries by IDs
func (s *MasterService) BulkDelete(ids []string) error {
	for _, id := range ids {
		if err := s.repo.Delete(id); err != nil {
			return err
		}
	}
	return nil
}

// ReorderEntries updates the order of master entries (no-op placeholder)
func (s *MasterService) ReorderEntries(category string, orderedIDs []string) error {
	entries, err := s.repo.GetByCategory(category)
	if err != nil {
		return err
	}
	entryMap := make(map[string]model.MdEntry)
	for _, entry := range entries {
		entryMap[entry.ID] = entry
	}
	for i, id := range orderedIDs {
		if entry, exists := entryMap[id]; exists {
			_ = i
			_ = entry
		}
	}
	return nil
}

// ImportFromExcel bulk-imports master entries from xlsx bytes.
// Returns count of imported records and any error.
func (s *MasterService) ImportFromExcel(category string, data []byte) (int, error) {
	entries, err := export.ParseMasterExcel(data)
	if err != nil {
		return 0, fmt.Errorf("failed to parse excel: %w", err)
	}

	imported := 0
	for _, entry := range entries {
		entry.Cat = model.MdCat(category)
		if entry.ID == "" {
			// Generate a simple ID if not provided
			entry.ID = fmt.Sprintf("%s-%s", category, strings.ToLower(strings.ReplaceAll(entry.Name, " ", "-")))
		}
		if err := s.repo.Create(&entry); err != nil {
			// Skip duplicates, continue
			continue
		}
		imported++
	}
	return imported, nil
}

// ExportToExcel exports all master entries for a category as xlsx bytes.
func (s *MasterService) ExportToExcel(category string) ([]byte, error) {
	entries, err := s.repo.GetByCategory(category)
	if err != nil {
		return nil, err
	}
	return export.GenerateMasterExcel(category, entries)
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// containsIgnoreCase reports whether s contains substr (case-insensitive).
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// sortByName sorts entries alphabetically by Name using a simple insertion sort.
func sortByName(entries []model.MdEntry) {
	for i := 1; i < len(entries); i++ {
		key := entries[i]
		j := i - 1
		for j >= 0 && entries[j].Name > key.Name {
			entries[j+1] = entries[j]
			j--
		}
		entries[j+1] = key
	}
}

// ParsePaginationParams extracts page and perPage from query params.
// Kept for backwards compatibility with existing handler usages.
func ParsePaginationParams(pageStr, perPageStr string) (int, int) {
	p := pagination.Parse(pageStr, perPageStr)
	return p.Page, p.PerPage
}
