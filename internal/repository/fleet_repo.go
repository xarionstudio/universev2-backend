package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"universev2-backend/internal/model"
)

type FleetRepo struct {
	db *gorm.DB
}

func NewFleetRepo(db *gorm.DB) *FleetRepo {
	return &FleetRepo{db: db}
}

func (r *FleetRepo) GetUnitStatuses() ([]model.Unit, error) {
	var rows []model.UnitStatusRow
	if err := r.db.Order("unit_code ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	var units []model.Unit
	for _, row := range rows {
		var histRows []model.UnitHistoryRow
		r.db.Where("unit_code = ?", row.Code).Order("id DESC").Limit(10).Find(&histRows)

		var hist []model.UnitHist
		for _, h := range histRows {
			hist = append(hist, model.UnitHist{h.HistWhen, h.HistWhat, h.HistWhy, h.HistStatus})
		}

		units = append(units, model.Unit{
			Code:   row.Code,
			Type:   row.TypeInfo,
			Status: row.Status,
			Loc:    row.Location,
			Upd:    row.UpdatedNote,
			Hist:   hist,
		})
	}

	return units, nil
}

func (r *FleetRepo) GetUnitHistory(code string) ([]model.UnitHist, error) {
	var rows []model.UnitHistoryRow
	err := r.db.Where("unit_code = ?", code).Order("id DESC").Limit(20).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	var hist []model.UnitHist
	for _, h := range rows {
		hist = append(hist, model.UnitHist{h.HistWhen, h.HistWhat, h.HistWhy, h.HistStatus})
	}
	return hist, nil
}

func (r *FleetRepo) UpdateUnitStatus(code string, status model.UnitStatus, note string) error {
	return r.db.Model(&model.UnitStatusRow{}).Where("unit_code = ?", code).
		Updates(map[string]interface{}{"status": string(status), "updated_note": note}).Error
}

func (r *FleetRepo) AddUnitHistory(code, when, what, why, status string) error {
	return r.db.Create(&model.UnitHistoryRow{
		UnitCode: code, HistWhen: when, HistWhat: what, HistWhy: why, HistStatus: status,
	}).Error
}

func (r *FleetRepo) GetFleetSettings() ([]model.FleetSetting, error) {
	var fleets []model.FleetSetting
	if err := r.db.Order("id ASC").Find(&fleets).Error; err != nil {
		return nil, err
	}

	for i, f := range fleets {
		var fsu []model.FleetSettingUnit
		r.db.Where("fleet_setting_id = ?", f.ID).Find(&fsu)
		var codes []string
		for _, u := range fsu {
			codes = append(codes, u.UnitCode)
		}
		fleets[i].Units = codes
	}

	return fleets, nil
}

func (r *FleetRepo) CreateFleetSetting(f *model.FleetSetting) error {
	if err := r.db.Create(f).Error; err != nil {
		return err
	}
	for _, code := range f.Units {
		r.db.Create(&model.FleetSettingUnit{FleetSettingID: f.ID, UnitCode: code})
	}
	return nil
}

func (r *FleetRepo) UpdateFleetSetting(id uint, f *model.FleetSetting) error {
	if err := r.db.Model(&model.FleetSetting{}).Where("id = ?", id).Updates(f).Error; err != nil {
		return err
	}
	r.db.Where("fleet_setting_id = ?", id).Delete(&model.FleetSettingUnit{})
	for _, code := range f.Units {
		r.db.Create(&model.FleetSettingUnit{FleetSettingID: id, UnitCode: code})
	}
	return nil
}

func (r *FleetRepo) DeleteFleetSetting(id uint) error {
	r.db.Where("fleet_setting_id = ?", id).Delete(&model.FleetSettingUnit{})
	return r.db.Where("id = ?", id).Delete(&model.FleetSetting{}).Error
}

func (r *FleetRepo) GetAllocations(date, shift string) ([]model.FleetAlloc, error) {
	var allocs []model.FleetAlloc
	q := r.db.Model(&model.FleetAlloc{})
	if date != "" {
		q = q.Where("alloc_date = ?", date)
	}
	if shift != "" {
		q = q.Where("shift = ?", shift)
	}
	err := q.Order("id ASC").Find(&allocs).Error
	if err != nil {
		return nil, err
	}

	// Load operators for each allocation
	for i := range allocs {
		var ops []model.FleetAllocOperator
		r.db.Where("allocation_id = ?", allocs[i].ID).Find(&ops)
		units := make(map[string]string)
		for _, op := range ops {
			units[op.UnitCode] = op.OperatorNIK
		}
		allocs[i].Units = units
	}

	return allocs, nil
}

// GetAllocationOperators returns operator NIK per unit code for a date+shift
func (r *FleetRepo) GetAllocationOperators(date, shift string) (map[string]string, error) {
	var allocs []model.FleetAlloc
	q := r.db.Model(&model.FleetAlloc{})
	if date != "" {
		q = q.Where("alloc_date = ?", date)
	}
	if shift != "" {
		q = q.Where("shift = ?", shift)
	}
	if err := q.Find(&allocs).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, a := range allocs {
		var ops []model.FleetAllocOperator
		r.db.Where("allocation_id = ?", a.ID).Find(&ops)
		for _, op := range ops {
			result[op.UnitCode] = op.OperatorNIK
		}
	}
	return result, nil
}

// GetOperatorNameByNIK returns employee name for a NIK
func (r *FleetRepo) GetOperatorNameByNIK(nik string) string {
	var emp model.Employee
	if err := r.db.Select("name").Where("nik = ?", nik).First(&emp).Error; err != nil {
		return ""
	}
	return emp.Name
}

// AutoAllocate performs auto allocation with operator matching
func (r *FleetRepo) AutoAllocate(date, shift string) error {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	if shift == "" {
		shift = "pagi"
	}

	fleets, err := r.GetFleetSettings()
	if err != nil {
		return err
	}

	// Load available units (active, not breakdown, not standby)
	var units []model.UnitDb
	if err := r.db.Where("is_active = ? AND is_breakdown = ? AND is_standby = ?", true, false, false).Find(&units).Error; err != nil {
		return err
	}
	unitByCode := make(map[string]model.UnitDb)
	for _, u := range units {
		unitByCode[u.Code] = u
	}

	// Load active operators with competencies
	type EmpComp struct {
		NIK       string
		Name      string
		ClassName string
	}
	var empComps []EmpComp
	if err := r.db.Table("employees").
		Select("employees.nik, employees.name, employee_competencies.class_name").
		Joins("JOIN employee_competencies ON employee_competencies.employee_id = employees.id").
		Where("employees.status = ?", "aktif").
		Scan(&empComps).Error; err != nil {
		return err
	}

	// Build operator competency map: nik -> []class
	opClasses := make(map[string][]string)
	for _, e := range empComps {
		if _, ok := opClasses[e.NIK]; !ok {
			opClasses[e.NIK] = []string{}
		}
		opClasses[e.NIK] = append(opClasses[e.NIK], e.ClassName)
	}

	// Delete existing allocation for this date+shift to re-allocate cleanly
	r.db.Where("alloc_date = ? AND shift = ?", date, shift).Delete(&model.FleetAlloc{})

	usedOps := make(map[string]bool)

	for _, f := range fleets {
		if !f.Active {
			continue
		}
		alloc := model.FleetAlloc{
			Date:   date,
			Shift:  shift,
			FlID:   f.ID,
			Digger: f.Digger,
			Loc:    f.Loc,
			Bus:    f.Bus,
		}
		if err := r.db.Create(&alloc).Error; err != nil {
			return err
		}

		// Assign operators to each unit in the fleet
		unitCodes := append([]string{f.Digger}, f.Units...)
		for _, code := range unitCodes {
			u, ok := unitByCode[code]
			if !ok {
				continue
			}
			// Find an operator whose competency matches the unit's EGI type
			tegi := typeiEgi(u.EGI)
			nik := ""
			for opNik, classes := range opClasses {
				if usedOps[opNik] {
					continue
				}
				for _, cls := range classes {
					if cls == tegi {
						nik = opNik
						break
					}
				}
				if nik != "" {
					break
				}
			}
			if nik == "" {
				continue
			}
			usedOps[nik] = true
			r.db.Create(&model.FleetAllocOperator{
				AllocationID: alloc.ID,
				UnitCode:     code,
				OperatorNIK:  nik,
			})
		}
	}
	return nil
}

// typeiEgi maps unit EGI model to competency class name (same logic as FE typeOfEgi)
func typeiEgi(egi string) string {
	e := egi
	switch {
	case containsAny(e, "WT", "-WT", "SYM3256", "FM260"):
		return "WATER TRUCK"
	case containsAny(e, "FT"):
		return "FUEL TRUCK"
	case containsAny(e, "ST", "CT", "XTREM", "WB4300"):
		return "SUPPORT TRUCK"
	case containsAny(e, "MH", "-MH"):
		return "MANHAUL"
	case containsAny(e, "F84G"):
		return "BUS"
	case containsAny(e, "TRITON", "PAJERO", "COLDDIESEL", "FE71"):
		return "LIGHT VEHICLE"
	case containsAny(e, "785", "777"):
		return "HD 785 / 777"
	case containsAny(e, "465", "773", "TR60"):
		return "HD 465 / 773"
	case containsAny(e, "R100E"):
		return "VOLVO"
	case containsAny(e, "SKT130"):
		return "SKT130"
	case containsAny(e, "SKT105"):
		return "SKT105"
	case containsAny(e, "SYZ440", "SYZ320"):
		return "SANY SYZ 440"
	case containsAny(e, "2600"):
		return "PC 2600"
	case containsAny(e, "2000"):
		return "PC 2000"
	case containsAny(e, "1250"):
		return "PC 1250"
	case containsAny(e, "1200"):
		return "PC 1200"
	case containsAny(e, "6020"):
		return "PC 6020"
	case containsAny(e, "870", "SY750"):
		return "PC 870"
	case containsAny(e, "470"):
		return "PC 470"
	case containsAny(e, "ZX350", "350", "SY365"):
		return "PC 350"
	case containsAny(e, "ZX210", "ZX200", "PC200", "200-LA", "SY215", "215"):
		return "PC 200"
	case containsAny(e, "D375", "D9"):
		return "D9-375"
	case containsAny(e, "D155", "D8T", "D360"):
		return "D8-155"
	case containsAny(e, "D6", "D85", "D260"):
		return "D6-D85SS"
	case containsAny(e, "GD825", "GD755", "16GC", "14M"):
		return "GRADER"
	case containsAny(e, "P410", "P460"):
		return "SCANIA P410"
	case containsAny(e, "DM30"):
		return "DRILL"
	case containsAny(e, "BW2"):
		return "COMPACTOR"
	default:
		return "SPARE"
	}
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if len(sub) > 0 && len(s) >= len(sub) && indexOf(s, sub) >= 0 {
			return true
		}
	}
	return false
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func (r *FleetRepo) GetUnitDB() ([]model.UnitDb, error) {
	var units []model.UnitDb
	err := r.db.Order("code ASC").Find(&units).Error
	return units, err
}

func (r *FleetRepo) CreateUnitDB(u *model.UnitDb) error {
	return r.db.Create(u).Error
}

func (r *FleetRepo) UpdateUnitDB(u *model.UnitDb) error {
	return r.db.Model(&model.UnitDb{}).Where("id = ?", u.ID).Updates(u).Error
}

func (r *FleetRepo) DeleteUnitDB(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.UnitDb{}).Error
}

func (r *FleetRepo) BulkCreateUnitDB(units []model.UnitDb) (imported int, skipped int, rowErrors []string, err error) {
	for _, u := range units {
		var count int64
		r.db.Model(&model.UnitDb{}).Where("code = ?", u.Code).Count(&count)
		if count > 0 {
			r.db.Model(&model.UnitDb{}).Where("code = ?", u.Code).Updates(map[string]interface{}{
				"class_name": u.Cls,
				"egi":        u.EGI,
				"product":    u.Product,
				"work_area":  u.Area,
				"location":   u.Loc,
				"upd_date":   u.Upd,
				"upd_by":     u.By,
			})
			skipped++
			rowErrors = append(rowErrors, fmt.Sprintf("Code %q: already exists, updated record", u.Code))
		} else {
			if err := r.db.Create(&u).Error; err == nil {
				imported++
			} else {
				skipped++
				rowErrors = append(rowErrors, fmt.Sprintf("Code %q: %v", u.Code, err))
			}
		}
	}
	return imported, skipped, rowErrors, nil
}
