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

// ── MasterRepo ────────────────────────────────────────────────────────────────

type MasterRepo struct {
	db *gorm.DB
}

func NewMasterRepo(db *gorm.DB) *MasterRepo {
	return &MasterRepo{db: db}
}

// GetByCategory returns all entries for a category as interface{} (slice of specific type)
func (r *MasterRepo) GetByCategory(cat string) (interface{}, error) {
	switch cat {
	case "egi":
		var entries []model.MasterEGIType
		err := r.db.Order("name ASC").Find(&entries).Error
		return entries, err
	case "product":
		var entries []model.MasterProduct
		err := r.db.Order("name ASC").Find(&entries).Error
		return entries, err
	case "eqclass":
		var entries []model.MasterEqClass
		err := r.db.Order("name ASC").Find(&entries).Error
		return entries, err
	case "area":
		var entries []model.MasterArea
		err := r.db.Order("name ASC").Find(&entries).Error
		return entries, err
	case "tempudo":
		var entries []model.MasterTempudo
		err := r.db.Order("name ASC").Find(&entries).Error
		return entries, err
	case "bus":
		var entries []model.MasterBus
		err := r.db.Order("name ASC").Find(&entries).Error
		return entries, err
	case "lokasiex":
		var entries []model.MasterLocationEx
		err := r.db.Order("name ASC").Find(&entries).Error
		return entries, err
	case "mess":
		var entries []model.MasterMess
		err := r.db.Order("name ASC").Find(&entries).Error
		return entries, err
	case "runtext":
		var entries []model.MasterRunningText
		err := r.db.Order("name ASC").Find(&entries).Error
		return entries, err
	}
	return nil, nil
}

// Create inserts a new entry for the given category
func (r *MasterRepo) Create(cat string, entry interface{}) error {
	return r.db.Create(entry).Error
}

// Update updates an entry by code for the given category
func (r *MasterRepo) Update(cat string, code string, updates map[string]interface{}) error {
	tbl := categoryTableMap(cat)
	if tbl == nil {
		return nil
	}
	return r.db.Model(tbl).Where("code = ?", code).Updates(updates).Error
}

// Delete deletes an entry by code for the given category
func (r *MasterRepo) Delete(cat string, code string) error {
	tbl := categoryTableMap(cat)
	if tbl == nil {
		return nil
	}
	return r.db.Where("code = ?", code).Delete(tbl).Error
}

// categoryTableMap maps category key to its model struct for GORM
func categoryTableMap(cat string) interface{} {
	switch cat {
	case "egi":
		return &model.MasterEGIType{}
	case "product":
		return &model.MasterProduct{}
	case "eqclass":
		return &model.MasterEqClass{}
	case "area":
		return &model.MasterArea{}
	case "tempudo":
		return &model.MasterTempudo{}
	case "bus":
		return &model.MasterBus{}
	case "lokasiex":
		return &model.MasterLocationEx{}
	case "mess":
		return &model.MasterMess{}
	case "runtext":
		return &model.MasterRunningText{}
	default:
		return nil
	}
}
