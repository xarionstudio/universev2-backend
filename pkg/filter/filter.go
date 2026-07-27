package filter

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// Params holds all possible filter parameters parsed from a request.
type Params struct {
	Search   string            // Full-text search across configured columns
	Status   string            // Exact match on status column
	Dept     string            // Exact match on dept column
	NIK      string            // Exact match on nik column
	DateFrom string            // Range start on configured date column (YYYY-MM-DD)
	DateTo   string            // Range end on configured date column (YYYY-MM-DD)
	Month    string            // Exact month match (YYYY-MM)
	Fields   map[string]string // Dynamic field=value pairs  (field[] query params)
	Logic    string            // "and" (default) | "or"
}

// ParseFromCtx reads standard filter query params from a Fiber context.
func ParseFromCtx(c fiber.Ctx) Params {
	fields := make(map[string]string)
	// Allow ?fields[column]=value style dynamic filters
	c.Request().URI().QueryArgs().VisitAll(func(key, val []byte) {
		k := string(key)
		if strings.HasPrefix(k, "fields[") && strings.HasSuffix(k, "]") {
			col := k[7 : len(k)-1]
			if col != "" {
				fields[col] = string(val)
			}
		}
	})

	logic := strings.ToLower(c.Query("logic", "and"))
	if logic != "or" {
		logic = "and"
	}

	return Params{
		Search:   c.Query("search"),
		Status:   c.Query("status"),
		Dept:     c.Query("dept"),
		NIK:      c.Query("nik"),
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Month:    c.Query("month"),
		Fields:   fields,
		Logic:    logic,
	}
}

// Options configures how Apply maps filter params to DB columns.
type Options struct {
	// SearchColumns lists table-qualified columns used for ILIKE search.
	// Example: []string{"employees.name", "employees.nik"}
	SearchColumns []string
	// DateColumn is the table-qualified column used for date range filters.
	// Example: "employees.join_date"
	DateColumn string
	// StatusColumn overrides default "status" column name if needed.
	StatusColumn string
	// DeptColumn overrides default "dept" column name if needed.
	DeptColumn string
	// MonthColumn is the column used for month exact match.
	MonthColumn string
}

// Apply returns a *gorm.DB with all filter conditions applied.
// It respects the Logic field ("and"/"or") for combining conditions.
func Apply(q *gorm.DB, p Params, opts Options) *gorm.DB {
	var conditions []string
	var args []interface{}

	// Search across multiple columns (always OR between columns)
	if p.Search != "" && len(opts.SearchColumns) > 0 {
		var parts []string
		for _, col := range opts.SearchColumns {
			parts = append(parts, col+" ILIKE ?")
			args = append(args, "%"+p.Search+"%")
		}
		conditions = append(conditions, "("+strings.Join(parts, " OR ")+")")
	}

	// Status filter
	if p.Status != "" {
		col := opts.StatusColumn
		if col == "" {
			col = "status"
		}
		conditions = append(conditions, col+" = ?")
		args = append(args, p.Status)
	}

	// Dept filter
	if p.Dept != "" {
		col := opts.DeptColumn
		if col == "" {
			col = "dept"
		}
		conditions = append(conditions, col+" = ?")
		args = append(args, p.Dept)
	}

	// NIK filter (exact)
	if p.NIK != "" {
		conditions = append(conditions, "employee_nik = ?")
		args = append(args, p.NIK)
	}

	// Date range filter
	if p.DateFrom != "" && opts.DateColumn != "" {
		conditions = append(conditions, opts.DateColumn+" >= ?")
		args = append(args, p.DateFrom)
	}
	if p.DateTo != "" && opts.DateColumn != "" {
		conditions = append(conditions, opts.DateColumn+" <= ?")
		args = append(args, p.DateTo)
	}

	// Month exact match
	if p.Month != "" && opts.MonthColumn != "" {
		conditions = append(conditions, opts.MonthColumn+" = ?")
		args = append(args, p.Month)
	}

	// Dynamic field filters
	for col, val := range p.Fields {
		if col != "" && val != "" {
			conditions = append(conditions, col+" = ?")
			args = append(args, val)
		}
	}

	if len(conditions) == 0 {
		return q
	}

	joiner := " AND "
	if p.Logic == "or" {
		joiner = " OR "
	}

	combined := strings.Join(conditions, joiner)
	return q.Where(combined, args...)
}
