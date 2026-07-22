package repository

import (
	"time"

	"gorm.io/gorm"

	"universev2-backend/internal/model"
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

	var row model.AttendanceRow
	err := r.db.Where("employee_nik = ? AND attendance_date = ?", nik, today).First(&row).Error
	if err == gorm.ErrRecordNotFound {
		row = model.AttendanceRow{
			NIK:   nik,
			Date:  today,
			Code:  "D",
			In:    nowTime,
			InM:   machine,
			St:    "hadir",
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

	var row model.AttendanceRow
	err := r.db.Where("employee_nik = ? AND attendance_date = ?", nik, today).First(&row).Error
	if err == gorm.ErrRecordNotFound {
		row = model.AttendanceRow{
			NIK:   nik,
			Date:  today,
			Code:  "D",
			Out:   nowTime,
			OutM:  machine,
			St:    "hadir",
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
