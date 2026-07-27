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
		r.db.Select("name, dept, pos").Where("nik = ?", s.EmployeeNIK).First(&emp)

		var badgeRows []model.PrestasiBadge
		r.db.Where("employee_nik = ?", s.EmployeeNIK).Find(&badgeRows)
		var badges []string
		for _, b := range badgeRows {
			badges = append(badges, b.BadgeKey)
		}

		attRate := 0.0
		if s.AttCount > 0 && s.TotalScheduledDays > 0 {
			attRate = float64(s.AttCount) / float64(s.TotalScheduledDays)
		}
		sleepRate := 0.0
		if s.SleepOkCount > 0 && s.TotalScheduledDays > 0 {
			sleepRate = float64(s.SleepOkCount) / float64(s.TotalScheduledDays)
		}

		records = append(records, model.PrestasiRecord{
			NIK:           s.EmployeeNIK,
			Name:          emp.Name,
			Dept:          emp.Dept,
			Pos:           emp.Pos,
			Foto:          emp.Foto,
			Rank:          s.Rank,
			Points:        s.TotalPoints,
			BestStreak:    s.StreakDays,
			CurrentStreak: s.CurrentStreak,
			AttCount:      s.AttCount,
			SleepPct:      s.SleepPct,
			AttRate:       attRate,
			SleepRate:     sleepRate,
			AvgSleepMin:   s.AvgSleepMin,
			Badges:        badges,
			LateCount:     s.LateCount,
			PenaltyDays:   s.PenaltyCount,
			QualifiedDays: s.QualifiedDays,
			ScheduledDays: s.TotalScheduledDays,
			CoverDays:     s.CoverDays,
		})
	}

	return records, nil
}

func (r *PrestasiRepo) GetOperatorHistory(nik string, days int) ([]model.PrestasiHistoryEntry, error) {
	var entries []model.PrestasiHistoryEntry
	q := r.db.Where("employee_nik = ?", nik)
	if days > 0 {
		q = q.Where("period_days = ?", days)
	}
	err := q.Order("record_date ASC").Find(&entries).Error
	return entries, err
}

func (r *PrestasiRepo) GetAllEmployeeNIKs() ([]model.Employee, error) {
	var emps []model.Employee
	err := r.db.Select("nik, name, dept, pos, foto").Order("nik ASC").Find(&emps).Error
	return emps, err
}

func (r *PrestasiRepo) SavePrestasiData(periodDays int, scores []model.PrestasiScore, history []model.PrestasiHistoryEntry, badges []model.PrestasiBadge) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing entries for the given period to overwrite with recalculation
		if err := tx.Where("period_days = ?", periodDays).Delete(&model.PrestasiScore{}).Error; err != nil {
			return err
		}
		if err := tx.Where("period_days = ?", periodDays).Delete(&model.PrestasiHistoryEntry{}).Error; err != nil {
			return err
		}
		for _, s := range scores {
			s.PeriodDays = periodDays
			if err := tx.Create(&s).Error; err != nil {
				return err
			}
		}
		for _, h := range history {
			h.PeriodDays = periodDays
			if err := tx.Create(&h).Error; err != nil {
				return err
			}
		}
		for _, b := range badges {
			tx.Where("employee_nik = ? AND badge_key = ?", b.EmployeeNIK, b.BadgeKey).FirstOrCreate(&b)
		}
		return nil
	})
}

func (r *PrestasiRepo) GetAttendanceRecord(nik string, dateStr string) (*model.AttendanceRow, error) {
	var att model.AttendanceRow
	err := r.db.Where("employee_nik = ? AND attendance_date = ?", nik, dateStr).First(&att).Error
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (r *PrestasiRepo) GetFTWRecord(nik string, dateStr string) (*model.FTWRecord, error) {
	var ftw model.FTWRecord
	err := r.db.Where("employee_nik = ? AND log_date = ?", nik, dateStr).First(&ftw).Error
	if err != nil {
		return nil, err
	}
	return &ftw, nil
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
