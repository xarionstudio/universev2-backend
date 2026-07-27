package repository

import (
	"gorm.io/gorm"

	"universev2-backend/internal/model"
	"universev2-backend/pkg/filter"
	"universev2-backend/pkg/pagination"
)

type FTWRepo struct {
	db *gorm.DB
}

func NewFTWRepo(db *gorm.DB) *FTWRepo {
	return &FTWRepo{db: db}
}

// GetTodayLogs returns FTW records for a specific date (non-paginated, kept for backward compat).
func (r *FTWRepo) GetTodayLogs(date string) ([]model.FTWRecord, error) {
	var logs []model.FTWRecord
	err := r.db.
		Select("ftw_logs.*, employees.name, employees.dept").
		Joins("LEFT JOIN employees ON ftw_logs.employee_nik = employees.nik").
		Where("ftw_logs.log_date = ?", date).
		Order("ftw_logs.submitted_at DESC").
		Find(&logs).Error
	return logs, err
}

// GetLogsPaginated returns a paginated + filtered page of FTW records.
// Supported filter fields: Status (exact), NIK (exact), DateFrom/DateTo on log_date, Dept via join.
func (r *FTWRepo) GetLogsPaginated(f filter.Params, p pagination.Params) ([]model.FTWRecord, int64, error) {
	var logs []model.FTWRecord
	var total int64

	q := r.db.Model(&model.FTWRecord{}).
		Select("ftw_logs.*, employees.name, employees.dept").
		Joins("LEFT JOIN employees ON ftw_logs.employee_nik = employees.nik")

	// Apply shared filter (NIK handled as employee_nik via filter.Apply)
	q = filter.Apply(q, f, filter.Options{
		SearchColumns: []string{"employees.name", "ftw_logs.employee_nik"},
		DateColumn:    "ftw_logs.log_date",
		StatusColumn:  "ftw_logs.status",
		DeptColumn:    "employees.dept",
	})

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := q.Order("ftw_logs.submitted_at DESC").
		Limit(p.PerPage).Offset(p.Offset()).
		Find(&logs).Error
	return logs, total, err
}

func (r *FTWRepo) Submit(rec *model.FTWRecord) error {
	ev := model.EvaluateFTW(rec.SleepMin)
	rec.St = ev.Status
	rec.RestHours = ev.RestHours
	rec.CanWork = ev.CanWork
	return r.db.Create(rec).Error
}

func (r *FTWRepo) GetHistory(nik string, days int) ([]model.FTWRecord, error) {
	var logs []model.FTWRecord
	q := r.db.Where("employee_nik = ?", nik).Order("log_date DESC")
	if days > 0 {
		q = q.Limit(days)
	}
	err := q.Find(&logs).Error
	return logs, err
}

// GetAllFiltered returns all FTW records matching a filter (used for Excel export, no page limit).
func (r *FTWRepo) GetAllFiltered(f filter.Params) ([]model.FTWRecord, error) {
	var logs []model.FTWRecord
	q := r.db.Model(&model.FTWRecord{}).
		Select("ftw_logs.*, employees.name, employees.dept").
		Joins("LEFT JOIN employees ON ftw_logs.employee_nik = employees.nik")

	q = filter.Apply(q, f, filter.Options{
		SearchColumns: []string{"employees.name", "ftw_logs.employee_nik"},
		DateColumn:    "ftw_logs.log_date",
		StatusColumn:  "ftw_logs.status",
		DeptColumn:    "employees.dept",
	})

	err := q.Order("ftw_logs.log_date DESC, employees.name ASC").Find(&logs).Error
	return logs, err
}
