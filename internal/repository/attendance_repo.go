package repository

import (
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"universev/internal/model"
)

type AttendanceRepo struct {
	db *gorm.DB
}

func NewAttendanceRepo(db *gorm.DB) *AttendanceRepo {
	return &AttendanceRepo{db: db}
}

type attendanceJoinResult struct {
	model.AttendanceRow
	EmpName string `gorm:"column:emp_name"`
	EmpDept string `gorm:"column:emp_dept"`
}

func (r *AttendanceRepo) GetLogsByDate(dateStr string) ([]model.AttendanceRow, error) {
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	var results []attendanceJoinResult
	err := r.db.Table("attendance_logs").
		Select("attendance_logs.*, employees.name as emp_name, employees.dept as emp_dept").
		Joins("LEFT JOIN employees ON employees.nik = attendance_logs.employee_nik").
		Where("attendance_logs.attendance_date = ?", dateStr).
		Order("employees.name ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	rows := make([]model.AttendanceRow, len(results))
	for i, res := range results {
		row := res.AttendanceRow
		row.Name = res.EmpName
		row.Dept = res.EmpDept
		rows[i] = row
	}

	return rows, nil
}

func deriveShiftCode() string {
	hour := time.Now().Hour()
	// Day shift: 04:00 - 17:59, Night shift: 18:00 - 03:59
	if hour >= 4 && hour < 18 {
		return "D"
	}
	return "N"
}

func (r *AttendanceRepo) GetLogsRange(from, to string) ([]model.AttendanceRow, error) {
	var results []attendanceJoinResult
	err := r.db.Table("attendance_logs").
		Select("attendance_logs.*, employees.name as emp_name, employees.dept as emp_dept").
		Joins("LEFT JOIN employees ON employees.nik = attendance_logs.employee_nik").
		Where("attendance_logs.attendance_date BETWEEN ? AND ?", from, to).
		Order("attendance_logs.attendance_date DESC, employees.name ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	rows := make([]model.AttendanceRow, len(results))
	for i, res := range results {
		row := res.AttendanceRow
		row.Name = res.EmpName
		row.Dept = res.EmpDept
		rows[i] = row
	}

	return rows, nil
}

func (r *AttendanceRepo) RecordCheckIn(nik, machine string) (*model.AttendanceRow, error) {
	today := time.Now().Format("2006-01-02")
	nowTime := time.Now().Format("15:04")
	sc := deriveShiftCode()

	var row model.AttendanceRow
	err := r.db.Where("employee_nik = ? AND attendance_date = ?", nik, today).First(&row).Error
	if err == gorm.ErrRecordNotFound {
		row = model.AttendanceRow{
			NIK:  nik,
			Date: today,
			Code: sc,
			In:   nowTime,
			InM:  machine,
			St:   "hadir",
		}
		if err := r.db.Create(&row).Error; err != nil {
			return nil, err
		}
	} else if err == nil {
		row.In = nowTime
		row.InM = machine
		row.St = "hadir"
		if err := r.db.Save(&row).Error; err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return &row, nil
}

func (r *AttendanceRepo) RecordCheckOut(nik, machine string) (*model.AttendanceRow, error) {
	today := time.Now().Format("2006-01-02")
	nowTime := time.Now().Format("15:04")
	shiftCode := deriveShiftCode()

	var row model.AttendanceRow
	err := r.db.Where("employee_nik = ? AND attendance_date = ?", nik, today).First(&row).Error
	if err == gorm.ErrRecordNotFound {
		row = model.AttendanceRow{
			NIK:  nik,
			Date: today,
			Code: shiftCode,
			Out:  nowTime,
			OutM: machine,
			St:   "hadir",
		}
		if err := r.db.Create(&row).Error; err != nil {
			return nil, err
		}
	} else if err == nil {
		row.Out = nowTime
		row.OutM = machine
		if err := r.db.Save(&row).Error; err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return &row, nil
}

func deriveShiftCodeForTime(t time.Time) string {
	hour := t.Hour()
	if hour >= 4 && hour < 18 {
		return "D"
	}
	return "N"
}

func (r *AttendanceRepo) RecordScanWithTimestamp(nik, machine string, timestamp time.Time) (*model.AttendanceRow, error) {
	dateStr := timestamp.Format("2006-01-02")
	timeStr := timestamp.Format("15:04")
	sc := deriveShiftCodeForTime(timestamp)

	var row model.AttendanceRow
	err := r.db.Where("employee_nik = ? AND attendance_date = ?", nik, dateStr).First(&row).Error
	if err == gorm.ErrRecordNotFound {
		row = model.AttendanceRow{
			NIK:  nik,
			Date: dateStr,
			Code: sc,
			In:   timeStr,
			InM:  machine,
			St:   "hadir",
		}
		if err := r.db.Create(&row).Error; err != nil {
			return nil, err
		}
	} else if err == nil {
		if row.In == "" {
			row.In = timeStr
			row.InM = machine
			row.St = "hadir"
		} else {
			row.Out = timeStr
			row.OutM = machine
		}
		if err := r.db.Save(&row).Error; err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return &row, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Attendance status engine (single source of truth)
//
// Dashboard, Display TV and the roster attendance page all read rows from
// attendance_logs. Previously those rows only existed after someone checked in
// and always carried st="hadir", so "belum"/"off"/"terlambat" were never shown
// and numbers diverged across modules. These helpers rebuild the attendance
// board for a date from the roster (roster_schedules) + real clock-in/out logs
// and persist the computed status, so every reader agrees.
// ─────────────────────────────────────────────────────────────────────────────

// isWorkingShift reports whether a roster shift code expects a clocked-in work
// shift. Leave/sick/off codes mark the employee as "off" (no check-in expected).
func isWorkingShift(code string) bool {
	switch strings.ToUpper(strings.TrimSpace(code)) {
	case "OFF", "CR", "AL", "LWP", "LWOP", "PH", "PHD",
		"S", "A",
		"MCU", "MCR", "MCUF", "ISM", "OBC", "KRT",
		"TGS", "DNS", "TRV", "TR", "TRS", "IN",
		"TERM", "EOC", "RSG":
		return false
	default:
		return true
	}
}

// laterThan compares "HH:MM" strings; true when t is strictly after ref.
func laterThan(t, ref string) bool {
	parse := func(s string) (int, int, bool) {
		parts := strings.SplitN(strings.TrimSpace(s), ":", 2)
		if len(parts) != 2 {
			return 0, 0, false
		}
		h, e1 := strconv.Atoi(parts[0])
		m, e2 := strconv.Atoi(parts[1])
		if e1 != nil || e2 != nil {
			return 0, 0, false
		}
		return h, m, true
	}
	th, tm, ok1 := parse(t)
	rh, rm, ok2 := parse(ref)
	if !ok1 || !ok2 {
		return false
	}
	return th > rh || (th == rh && tm > rm)
}

// isLateShift decides whether a check-in is late for the given shift code.
// Day shifts (D, R, STB, …) start at 06:00; night shift (N) starts at 18:00.
func isLateShift(code, inTime string) bool {
	if strings.EqualFold(strings.TrimSpace(code), "N") {
		return laterThan(inTime, "18:00")
	}
	return laterThan(inTime, "06:00")
}

// SyncAttendanceBoard ensures attendance_logs has one row per roster-scheduled
// employee for the date with a consistent computed status. It never deletes
// existing rows and never downgrades an existing "unfit" status.
func (r *AttendanceRepo) SyncAttendanceBoard(dateStr string) error {
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	// Roster shift code of each employee scheduled on this date. When a date is
	// covered by several roster files (e.g. month boundary), the latest wins.
	type schedRow struct {
		NIK  string
		Code string
	}
	var rows []schedRow
	err := r.db.Table("roster_schedules").
		Select("roster_schedules.employee_nik AS nik, roster_schedules.shift_code AS code").
		Joins("JOIN roster_files ON roster_files.id = roster_schedules.roster_file_id").
		Where("roster_schedules.schedule_date = ?", dateStr).
		Order("roster_schedules.roster_file_id DESC").
		Scan(&rows).Error
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil // no roster for this date — leave existing rows untouched
	}

	type logRow struct {
		NIK  string
		In   string
		InM  string
		Out  string
		OutM string
		St   string
	}
	var logs []logRow
	if err := r.db.Table("attendance_logs").
		Where("attendance_date = ?", dateStr).
		Scan(&logs).Error; err != nil {
		return err
	}
	logMap := make(map[string]logRow, len(logs))
	for _, l := range logs {
		logMap[l.NIK] = l
	}

	for _, row := range rows {
		lg, _ := logMap[row.NIK]
		st := lg.St
		if !isWorkingShift(row.Code) {
			st = "off"
		} else if st != "unfit" { // never downgrade an existing unfit
			switch {
			case lg.In == "":
				st = "belum"
			case isLateShift(row.Code, lg.In):
				st = "terlambat"
			default:
				st = "hadir"
			}
		}
		if st == "" {
			st = "off"
		}

		updates := map[string]interface{}{
			"status":     st,
			"shift_code": row.Code,
		}
		if lg.In != "" {
			updates["check_in"] = lg.In
			updates["check_in_machine"] = lg.InM
		}
		if lg.Out != "" {
			updates["check_out"] = lg.Out
			updates["check_out_machine"] = lg.OutM
		}

		res := r.db.Model(&model.AttendanceRow{}).
			Where("employee_nik = ? AND attendance_date = ?", row.NIK, dateStr).
			Updates(updates)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			rec := model.AttendanceRow{
				NIK: row.NIK,
				Code: row.Code,
				Date: dateStr,
				St:   st,
				In:   lg.In,
				InM:  lg.InM,
				Out:  lg.Out,
				OutM: lg.OutM,
			}
			if err := r.db.Create(&rec).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// SyncAttendanceRange rebuilds the attendance board for every date in [from, to].
func (r *AttendanceRepo) SyncAttendanceRange(from, to string) error {
	if from == "" || to == "" {
		return nil
	}
	f, err := time.Parse("2006-01-02", from)
	if err != nil {
		return nil
	}
	t, err := time.Parse("2006-01-02", to)
	if err != nil {
		return nil
	}
	if t.Before(f) {
		f, t = t, f
	}
	// Guard against an unbounded loop on a bad range.
	if t.Sub(f) > 400*24*time.Hour {
		t = f.AddDate(0, 0, 400)
	}
	for day := f; !day.After(t); day = day.AddDate(0, 0, 1) {
		_ = r.SyncAttendanceBoard(day.Format("2006-01-02"))
	}
	return nil
}

