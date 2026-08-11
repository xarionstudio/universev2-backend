package repository

import (
	"encoding/json"
	"strconv"

	"gorm.io/gorm"

	"universev/internal/model"
)

type SettingsRepo struct {
	db *gorm.DB
}

func NewSettingsRepo(db *gorm.DB) *SettingsRepo {
	return &SettingsRepo{db: db}
}

type AppSettingsDB struct {
	ID          uint   `gorm:"column:id;primaryKey;autoIncrement"`
	AppName     string `gorm:"column:app_name"`
	AppEnv      string `gorm:"column:app_env"`
	CompanyLogo string `gorm:"column:company_logo"`
	Theme       string `gorm:"column:theme"`
	Lang        string `gorm:"column:lang"`
	MenuVisJSON string `gorm:"column:menu_vis_json"`
}

func (AppSettingsDB) TableName() string { return "app_settings" }

func (r *SettingsRepo) GetAppSettings() (model.AppSettings, error) {
	var dbRow AppSettingsDB
	err := r.db.First(&dbRow).Error
	if err != nil {
		defaultVis := map[string]bool{
			"display": true, "roster": true, "employees": true, "ftw": true,
			"asset": true, "prestasi": true, "master": true, "users": true, "settings": true,
		}
		return model.AppSettings{
			AppName: "universev", AppEnv: "development", CompanyLogo: "",
			Theme: "dark", Lang: "id", MenuVis: defaultVis,
		}, nil
	}

	var menuVis map[string]bool
	_ = json.Unmarshal([]byte(dbRow.MenuVisJSON), &menuVis)

	return model.AppSettings{
		AppName:     dbRow.AppName,
		AppEnv:      dbRow.AppEnv,
		CompanyLogo: dbRow.CompanyLogo,
		Theme:       dbRow.Theme,
		Lang:        dbRow.Lang,
		MenuVis:     menuVis,
	}, nil
}

func (r *SettingsRepo) UpdateAppSettings(s model.AppSettings) error {
	visBytes, _ := json.Marshal(s.MenuVis)
	// Try to get existing record
	var existing AppSettingsDB
	if err := r.db.First(&existing).Error; err != nil {
		// No existing record, create new one
		dbRow := AppSettingsDB{
			AppName:     s.AppName,
			AppEnv:      s.AppEnv,
			CompanyLogo: s.CompanyLogo,
			Theme:       s.Theme,
			Lang:        s.Lang,
			MenuVisJSON: string(visBytes),
		}
		return r.db.Create(&dbRow).Error
	}
	// Update existing
	existing.AppName = s.AppName
	existing.AppEnv = s.AppEnv
	existing.CompanyLogo = s.CompanyLogo
	existing.Theme = s.Theme
	existing.Lang = s.Lang
	existing.MenuVisJSON = string(visBytes)
	return r.db.Save(&existing).Error
}

// Audio Schedules

func (r *SettingsRepo) GetAudioSchedules() ([]model.AudioSchedule, error) {
	var audios []model.AudioSchedule
	if err := r.db.Order("id ASC").Find(&audios).Error; err != nil {
		return nil, err
	}

	for i, a := range audios {
		var ads []model.AudioScheduleDisplay
		r.db.Where("audio_id = ?", a.ID).Find(&ads)
		var kinds []string
		for _, d := range ads {
			kinds = append(kinds, d.DisplayKind)
		}
		audios[i].Displays = kinds
	}

	return audios, nil
}

func (r *SettingsRepo) CreateAudioSchedule(a *model.AudioSchedule) error {
	if err := r.db.Create(a).Error; err != nil {
		return err
	}
	for _, kind := range a.Displays {
		r.db.Create(&model.AudioScheduleDisplay{AudioID: a.ID, DisplayKind: kind})
	}
	return nil
}

func (r *SettingsRepo) UpdateAudioSchedule(id string, a *model.AudioSchedule) error {
	aid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	if err := r.db.Model(&model.AudioSchedule{}).Where("id = ?", uint(aid)).
		Select("title", "trigger_time", "frequency", "file_name", "is_active").
		Updates(a).Error; err != nil {
		return err
	}
	r.db.Where("audio_id = ?", uint(aid)).Delete(&model.AudioScheduleDisplay{})
	for _, kind := range a.Displays {
		r.db.Create(&model.AudioScheduleDisplay{AudioID: uint(aid), DisplayKind: kind})
	}
	return nil
}

func (r *SettingsRepo) DeleteAudioSchedule(id string) error {
	aid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	r.db.Where("audio_id = ?", uint(aid)).Delete(&model.AudioScheduleDisplay{})
	return r.db.Where("id = ?", uint(aid)).Delete(&model.AudioSchedule{}).Error
}

// Display Devices

func (r *SettingsRepo) GetDisplays(kind string) ([]model.DisplayDevice, error) {
	var displays []model.DisplayDevice
	q := r.db.Model(&model.DisplayDevice{})
	if kind != "" {
		q = q.Where("content_kind = ?", kind)
	}
	err := q.Order("id ASC").Find(&displays).Error
	if err != nil {
		return nil, err
	}

	// Load fleet_ids for monitor displays
	for i := range displays {
		if displays[i].Content == "monitor" {
			var dfs []model.DisplayFleet
			r.db.Where("display_id = ?", displays[i].ID).Order("sort_order ASC").Find(&dfs)
			for _, df := range dfs {
				displays[i].FleetIDs = append(displays[i].FleetIDs, df.FleetID)
			}
		}
	}

	return displays, nil
}

func (r *SettingsRepo) CreateDisplay(d *model.DisplayDevice) error {
	if err := r.db.Create(d).Error; err != nil {
		return err
	}
	// Save pivot for monitor displays
	if d.Content == "monitor" {
		for i, fid := range d.FleetIDs {
			r.db.Create(&model.DisplayFleet{DisplayID: d.ID, FleetID: fid, SortOrder: i})
		}
	}
	return nil
}

func (r *SettingsRepo) UpdateDisplay(id string, d *model.DisplayDevice) error {
	did, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	if err := r.db.Model(&model.DisplayDevice{}).Where("id = ?", uint(did)).
		Select("code", "name", "location", "content_kind", "fleet_id", "rotate_sec", "running_text", "is_online", "is_active").
		Updates(d).Error; err != nil {
		return err
	}
	// Replace pivot for monitor displays
	if d.Content == "monitor" {
		r.db.Where("display_id = ?", uint(did)).Delete(&model.DisplayFleet{})
		for i, fid := range d.FleetIDs {
			r.db.Create(&model.DisplayFleet{DisplayID: uint(did), FleetID: fid, SortOrder: i})
		}
	}
	return nil
}

func (r *SettingsRepo) DeleteDisplay(id string) error {
	did, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	r.db.Where("display_id = ?", uint(did)).Delete(&model.DisplayFleet{})
	return r.db.Where("id = ?", uint(did)).Delete(&model.DisplayDevice{}).Error
}

func (r *SettingsRepo) UpdateHeartbeat(code string, hb string) error {
	return r.db.Model(&model.DisplayDevice{}).Where("code = ?", code).
		Updates(map[string]interface{}{"is_online": true, "last_heartbeat": hb}).Error
}

// GetDeviceByCode returns a display device by its code.
func (r *SettingsRepo) GetDeviceByCode(code string) (*model.DisplayDevice, error) {
	var dev model.DisplayDevice
	if err := r.db.Where("code = ?", code).First(&dev).Error; err != nil {
		return nil, err
	}
	return &dev, nil
}

// Business Rules

func (r *SettingsRepo) GetAllBusinessRules() ([]model.BusinessRule, error) {
	var rules []model.BusinessRule
	if err := r.db.Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

func (r *SettingsRepo) GetBusinessRuleByCategory(category string) (*model.BusinessRule, error) {
	var rule model.BusinessRule
	if err := r.db.Where("category = ?", category).First(&rule).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *SettingsRepo) UpsertBusinessRule(category string, rulesJSON string, updatedBy string) error {
	var existing model.BusinessRule
	err := r.db.Where("category = ?", category).First(&existing).Error

	if err != nil {
		// Create new
		rule := model.BusinessRule{
			Category:  category,
			Rules:     rulesJSON,
			UpdatedBy: updatedBy,
		}
		return r.db.Create(&rule).Error
	}

	// Update existing
	existing.Rules = rulesJSON
	existing.UpdatedBy = updatedBy
	return r.db.Save(&existing).Error
}
