package pagination

import "strconv"

const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 200
)

// Params holds parsed pagination query parameters.
type Params struct {
	Page    int
	PerPage int
}

// Parse extracts page and perPage from raw query-string values.
// It applies sensible defaults and a maximum cap.
func Parse(pageStr, perPageStr string) Params {
	page := DefaultPage
	perPage := DefaultPerPage

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 {
		if pp > MaxPerPage {
			pp = MaxPerPage
		}
		perPage = pp
	}
	return Params{Page: page, PerPage: perPage}
}

// Offset returns the SQL OFFSET value for the current page.
func (p Params) Offset() int {
	return (p.Page - 1) * p.PerPage
}

// TotalPages calculates the number of pages given a total item count.
func TotalPages(total int64, perPage int) int {
	if perPage <= 0 {
		return 0
	}
	pages := int(total) / perPage
	if int(total)%perPage > 0 {
		pages++
	}
	return pages
}

// Meta is the pagination metadata included in API responses.
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"perPage"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

// BuildMeta assembles a Meta from params and a total count.
func BuildMeta(p Params, total int64) Meta {
	return Meta{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      total,
		TotalPages: TotalPages(total, p.PerPage),
	}
}
