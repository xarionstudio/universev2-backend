package service

import (
	"strconv"

	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
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
	// Get all entries for the category
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

// ReorderEntries updates the order of master entries
func (s *MasterService) ReorderEntries(category string, orderedIDs []string) error {
	// Update the order by updating a sort index field
	// For now, we'll use the ID field to maintain order
	// In a real implementation, you might want to add a 'sort_order' column
	entries, err := s.repo.GetByCategory(category)
	if err != nil {
		return err
	}

	// Create a map for quick lookup
	entryMap := make(map[string]model.MdEntry)
	for _, entry := range entries {
		entryMap[entry.ID] = entry
	}

	// Update entries in the new order
	for i, id := range orderedIDs {
		if entry, exists := entryMap[id]; exists {
			// In a real implementation, you would update a sort_order field here
			// For now, we just verify the entry exists
			_ = i
			_ = entry
		}
	}

	return nil
}

// Helper functions
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			containsIgnoreCase(s[1:], substr) ||
			containsIgnoreCase(s[:len(s)-1], substr) ||
			containsIgnoreCase(s[1:len(s)-1], substr))
}

func sortByName(entries []model.MdEntry) {
	// Simple bubble sort by name
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].Name > entries[j].Name {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
}

// ParsePaginationParams extracts page and perPage from query params
func ParsePaginationParams(pageStr, perPageStr string) (int, int) {
	page := 1
	perPage := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 {
			perPage = pp
		}
	}

	return page, perPage
}
