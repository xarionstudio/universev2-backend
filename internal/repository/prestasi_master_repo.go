package repository

import (
	"gorm.io/gorm"

	"universev2-backend/internal/model"
)

type PrestasiRepo struct {
	db *gorm.DB
}

func NewPrestasiRepo(db *gorm.DB) *PrestasiRepo {
	return &PrestasiRepo{db: db}
}

func (r *PrestasiRepo) GetLeaderboard(periodDays int) ([]model.PrestasiRecord, error) {
	var scores []model.PrestasiScore
	q := r.db.Model(&model.PrestasiScore{})
	if periodDays > 0 {
		q = q.Where("period_days = ?", periodDays)
	}
	if err := q.Order("rank ASC").Find(&scores).Error; err != nil {
		return nil, err
	}

	var records []model.PrestasiRecord
	for _, s := range scores {
		var emp model.Employee
		r.db.Select("name, dept").Where("nik = ?", s.EmployeeNIK).First(&emp)

		var badgeRows []model.PrestasiBadge
		r.db.Where("employee_nik = ?", s.EmployeeNIK).Find(&badgeRows)
		var badges []string
		for _, b := range badgeRows {
			badges = append(badges, b.BadgeKey)
		}

		records = append(records, model.PrestasiRecord{
			NIK:          s.EmployeeNIK,
			Name:         emp.Name,
			Dept:         emp.Dept,
			Rank:         s.Rank,
			PtsTotal:     s.TotalPoints,
			Streak:       s.StreakDays,
			AttCount:     s.AttCount,
			SleepPct:     s.SleepPct,
			Badges:       badges,
			LateCount:    s.LateCount,
			PenaltyCount: s.PenaltyCount,
		})
	}

	return records, nil
}

type MasterRepo struct {
	db *gorm.DB
}

func NewMasterRepo(db *gorm.DB) *MasterRepo {
	return &MasterRepo{db: db}
}

func (r *MasterRepo) GetByCategory(cat string) ([]model.MdEntry, error) {
	var entries []model.MdEntry
	err := r.db.Where("category_key = ?", cat).Order("name ASC").Find(&entries).Error
	return entries, err
}

func (r *MasterRepo) Create(entry *model.MdEntry) error {
	return r.db.Create(entry).Error
}

func (r *MasterRepo) Update(id string, entry *model.MdEntry) error {
	return r.db.Model(&model.MdEntry{}).Where("id = ?", id).Updates(entry).Error
}

func (r *MasterRepo) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.MdEntry{}).Error
}
