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
	if err := r.db.Order("code ASC").Find(&rows).Error; err != nil {
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
	return r.db.Model(&model.UnitStatusRow{}).Where("code = ?", code).
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

func (r *FleetRepo) UpdateFleetSetting(id string, f *model.FleetSetting) error {
	if err := r.db.Model(&model.FleetSetting{}).Where("id = ?", id).Updates(f).Error; err != nil {
		return err
	}
	r.db.Where("fleet_setting_id = ?", id).Delete(&model.FleetSettingUnit{})
	for _, code := range f.Units {
		r.db.Create(&model.FleetSettingUnit{FleetSettingID: id, UnitCode: code})
	}
	return nil
}

func (r *FleetRepo) DeleteFleetSetting(id string) error {
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
	return allocs, err
}

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

	for _, f := range fleets {
		if !f.Active {
			continue
		}
		allocID := fmt.Sprintf("alloc-%s-%s-%s", date, shift, f.ID)
		alloc := model.FleetAlloc{
			ID:     allocID,
			Date:   date,
			Shift:  shift,
			FlID:   f.ID,
			Digger: f.Digger,
			Loc:    f.Loc,
			Bus:    f.Bus,
		}
		r.db.Where("id = ?", allocID).FirstOrCreate(&alloc)
	}
	return nil
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

func (r *FleetRepo) DeleteUnitDB(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.UnitDb{}).Error
}
