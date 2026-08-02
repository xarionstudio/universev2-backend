package repository

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"universev2-backend/internal/model"
	"universev2-backend/pkg/filter"
	"universev2-backend/pkg/pagination"
)

type RosterRepo struct {
	db *gorm.DB
}

func NewRosterRepo(db *gorm.DB) *RosterRepo {
	return &RosterRepo{db: db}
}

func (r *RosterRepo) GetRosterByID(id string) (*model.RosterMeta, error) {
	var meta model.RosterMeta
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	err = r.db.Where("id = ?", uint(uid)).First(&meta).Error
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

func (r *RosterRepo) GetExportRosterData(fileId string, deptFilter string) ([]model.RosterExportRow, error) {
	type dbRow struct {
		NIK          string
		Name         string
		Dept         string
		Pos          string
		ScheduleDate time.Time
		ShiftCode    string
	}
	var dbRows []dbRow
	q := r.db.Table("roster_schedules").
		Select("employees.nik, employees.name, employees.dept, employees.pos, roster_schedules.schedule_date, roster_schedules.shift_code").
		Joins("JOIN employees ON roster_schedules.employee_nik = employees.nik")

	if fileId != "" {
		fid, err := strconv.ParseUint(fileId, 10, 64)
		if err == nil {
			q = q.Where("roster_schedules.roster_file_id = ?", uint(fid))
		}
	}

	err := q.Order("employees.name ASC, roster_schedules.schedule_date ASC").
		Scan(&dbRows).Error

	empMap := make(map[string]*model.RosterExportRow)
	var empOrder []string

	if err == nil && len(dbRows) > 0 {
		for _, row := range dbRows {
			if _, exists := empMap[row.NIK]; !exists {
				empMap[row.NIK] = &model.RosterExportRow{
					NIK:       row.NIK,
					Name:      row.Name,
					Dept:      row.Dept,
					Pos:       row.Pos,
					Schedules: make(map[int]string),
				}
				empOrder = append(empOrder, row.NIK)
			}
			day := row.ScheduleDate.Day()
			empMap[row.NIK].Schedules[day] = row.ShiftCode
		}
	} else {
		// Fallback: If no custom schedules uploaded for fileId, fetch active employees from DB
		var employees []model.Employee
		q2 := r.db.Model(&model.Employee{}).Where("status = ?", "aktif")
		if deptFilter != "" && deptFilter != "All Department" {
			q2 = q2.Where("dept = ?", deptFilter)
		}
		if err := q2.Order("name ASC").Find(&employees).Error; err != nil {
			return nil, err
		}
		for _, emp := range employees {
			empMap[emp.NIK] = &model.RosterExportRow{
				NIK:       emp.NIK,
				Name:      emp.Name,
				Dept:      emp.Dept,
				Pos:       emp.Pos,
				Schedules: make(map[int]string),
			}
			empOrder = append(empOrder, emp.NIK)
		}
	}

	result := make([]model.RosterExportRow, 0, len(empOrder))
	for _, nik := range empOrder {
		result = append(result, *empMap[nik])
	}
	return result, nil
}

func (r *RosterRepo) GetRosters(dept string) ([]model.RosterMeta, error) {
	var metas []model.RosterMeta
	q := r.db.Model(&model.RosterMeta{})
	if dept != "" {
		q = q.Where("dept = ?", dept)
	}
	err := q.Order("created_at DESC").Find(&metas).Error
	return metas, err
}

func (r *RosterRepo) GetRevisions(status string) ([]model.RosterRevision, error) {
	var revisions []model.RosterRevision
	q := r.db.
		Select("roster_revisions.*, employees.name").
		Joins("LEFT JOIN employees ON roster_revisions.employee_nik = employees.nik")
	if status != "" {
		q = q.Where("roster_revisions.status = ?", status)
	}
	err := q.Order("roster_revisions.created_at DESC").Find(&revisions).Error
	return revisions, err
}

func (r *RosterRepo) CountPendingRevisions() (int, error) {
	var count int64
	err := r.db.Model(&model.RosterRevision{}).Where("status = ?", "pending").Count(&count).Error
	return int(count), err
}

func (r *RosterRepo) ApproveRevision(id int, byId, byEn string) error {
	return r.db.Model(&model.RosterRevision{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": "approved", "by_id": byId, "by_en": byEn}).Error
}

func (r *RosterRepo) RejectRevision(id int, byId, byEn string) error {
	return r.db.Model(&model.RosterRevision{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": "rejected", "by_id": byId, "by_en": byEn}).Error
}

func (r *RosterRepo) DeleteRevision(id int) error {
	return r.db.Where("id = ?", id).Delete(&model.RosterRevision{}).Error
}

func (r *RosterRepo) CreateRoster(meta *model.RosterMeta) error {
	return r.db.Create(meta).Error
}

func (r *RosterRepo) CreateRevision(rev *model.RosterRevision) error {
	return r.db.Create(rev).Error
}

func (r *RosterRepo) GetAttendance(date string) ([]model.AttendanceRow, error) {
	var rows []model.AttendanceRow
	err := r.db.
		Select("attendance_logs.*, employees.name, employees.dept").
		Joins("LEFT JOIN employees ON attendance_logs.employee_nik = employees.nik").
		Where("attendance_logs.attendance_date = ?", date).
		Order("employees.name ASC").
		Find(&rows).Error
	return rows, err
}

func (r *RosterRepo) GetSchedulesByFile(fileId string) ([]model.RosterSchedule, error) {
	var schedules []model.RosterSchedule
	fid, err := strconv.ParseUint(fileId, 10, 64)
	if err != nil {
		return nil, err
	}
	err = r.db.Where("roster_file_id = ?", uint(fid)).Order("employee_nik, schedule_date").Find(&schedules).Error
	return schedules, err
}

func (r *RosterRepo) GetRosterDetail(key string) (fiber.Map, error) {
	type detailRow struct {
		NIK       string `json:"nik"`
		Name      string `json:"name"`
		Dept      string `json:"dept"`
		Pos       string `json:"pos"`
		Date      string `json:"date"`
		ShiftCode string `json:"shiftCode"`
	}
	var rows []detailRow
	q := r.db.Table("roster_schedules").
		Select("employees.nik, employees.name, employees.dept, employees.pos, roster_schedules.schedule_date as date, roster_schedules.shift_code").
		Joins("JOIN employees ON roster_schedules.employee_nik = employees.nik")

	if key != "" {
		fid, err := strconv.ParseUint(key, 10, 64)
		if err == nil {
			q = q.Where("roster_schedules.roster_file_id = ?", uint(fid))
		}
	}

	err := q.Order("employees.name ASC, roster_schedules.schedule_date ASC").
		Scan(&rows).Error

	if err != nil || len(rows) == 0 {
		return fiber.Map{
			"key":   key,
			"days":  []string{},
			"rows":  []fiber.Map{},
			"total": 0,
		}, nil
	}

	// Collect unique dates
	dateSet := make(map[string]bool)
	for _, r := range rows {
		dateSet[r.Date] = true
	}
	days := make([]string, 0, len(dateSet))
	for d := range dateSet {
		days = append(days, d)
	}

	// Group by employee
	empMap := make(map[string]*fiber.Map)
	var empOrder []string
	for _, r := range rows {
		if _, exists := empMap[r.NIK]; !exists {
			empMap[r.NIK] = &fiber.Map{
				"nik":   r.NIK,
				"name":  r.Name,
				"dept":  r.Dept,
				"pos":   r.Pos,
				"codes": []fiber.Map{},
			}
			empOrder = append(empOrder, r.NIK)
		}
	}

	// Build codes per employee per date
	empCodes := make(map[string]map[string]string)
	for _, r := range rows {
		if empCodes[r.NIK] == nil {
			empCodes[r.NIK] = make(map[string]string)
		}
		empCodes[r.NIK][r.Date] = r.ShiftCode
	}

	resultRows := make([]fiber.Map, 0, len(empOrder))
	for _, nik := range empOrder {
		codes := make([]fiber.Map, 0, len(days))
		for _, d := range days {
			code := empCodes[nik][d]
			if code == "" {
				code = "—"
			}
			codes = append(codes, fiber.Map{"date": d, "code": code})
		}
		resultRows = append(resultRows, fiber.Map{
			"nik":   nik,
			"name":  (*empMap[nik])["name"],
			"dept":  (*empMap[nik])["dept"],
			"pos":   (*empMap[nik])["pos"],
			"codes": codes,
		})
	}

	return fiber.Map{
		"key":   key,
		"days":  days,
		"rows":  resultRows,
		"total": len(resultRows),
	}, nil
}

// GetRostersPaginated returns a filtered, paginated list of roster metadata.
// Supported filters: Dept, Status (exact), Month (exact), DateFrom/DateTo on created_at.
func (r *RosterRepo) GetRostersPaginated(f filter.Params, p pagination.Params) ([]model.RosterMeta, int64, error) {
	var metas []model.RosterMeta
	var total int64

	q := r.db.Model(&model.RosterMeta{})
	q = filter.Apply(q, f, filter.Options{
		DateColumn:   "roster_files.created_at",
		StatusColumn: "roster_files.status",
		DeptColumn:   "roster_files.dept",
		MonthColumn:  "roster_files.month_period",
	})

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := q.Order("created_at DESC").Limit(p.PerPage).Offset(p.Offset()).Find(&metas).Error
	return metas, total, err
}

// GetRevisionsPaginated returns a filtered, paginated list of roster revisions.
// Supported filters: Status (exact), NIK, DateFrom/DateTo on created_at, search on employee name.
func (r *RosterRepo) GetRevisionsPaginated(f filter.Params, p pagination.Params) ([]model.RosterRevision, int64, error) {
	var revisions []model.RosterRevision
	var total int64

	q := r.db.Model(&model.RosterRevision{}).
		Select("roster_revisions.*, employees.name").
		Joins("LEFT JOIN employees ON roster_revisions.employee_nik = employees.nik")

	q = filter.Apply(q, f, filter.Options{
		SearchColumns: []string{"employees.name", "roster_revisions.employee_nik"},
		DateColumn:    "roster_revisions.created_at",
		StatusColumn:  "roster_revisions.status",
	})

	// Explicit NIK scoping on revisions table
	if f.NIK != "" {
		q = q.Where("roster_revisions.employee_nik = ?", f.NIK)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := q.Order("roster_revisions.created_at DESC").Limit(p.PerPage).Offset(p.Offset()).Find(&revisions).Error
	return revisions, total, err
}
