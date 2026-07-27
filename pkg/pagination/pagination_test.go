package pagination_test

import (
	"testing"

	"universev2-backend/pkg/pagination"
)

func TestParse(t *testing.T) {
	tests := []struct {
		pageStr    string
		perPageStr string
		wantPage   int
		wantLimit  int
	}{
		{"", "", 1, 20},
		{"2", "50", 2, 50},
		{"-1", "0", 1, 20},
		{"3", "999", 3, 200}, // caps at MaxPerPage (200)
	}

	for _, tt := range tests {
		p := pagination.Parse(tt.pageStr, tt.perPageStr)
		if p.Page != tt.wantPage || p.PerPage != tt.wantLimit {
			t.Errorf("Parse(%q, %q) = {%d, %d}; want {%d, %d}",
				tt.pageStr, tt.perPageStr, p.Page, p.PerPage, tt.wantPage, tt.wantLimit)
		}
	}
}

func TestTotalPages(t *testing.T) {
	tests := []struct {
		total     int64
		perPage   int
		wantPages int
	}{
		{0, 20, 0},
		{20, 20, 1},
		{21, 20, 2},
		{50, 10, 5},
	}

	for _, tt := range tests {
		got := pagination.TotalPages(tt.total, tt.perPage)
		if got != tt.wantPages {
			t.Errorf("TotalPages(%d, %d) = %d; want %d", tt.total, tt.perPage, got, tt.wantPages)
		}
	}
}
