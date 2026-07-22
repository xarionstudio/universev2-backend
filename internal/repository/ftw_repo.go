package repository

import (
	"gorm.io/gorm"

	"universev2-backend/internal/model"
)

type FTWRepo struct {
	db *gorm.DB
}

func NewFTWRepo(db *gorm.DB) *FTWRepo {
	return &FTWRepo{db: db}
}

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
