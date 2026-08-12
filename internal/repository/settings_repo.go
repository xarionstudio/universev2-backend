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
	AppDesc     string `gorm:"column:app_desc"`
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
			AppName: "universev", AppDesc: "", AppEnv: "development", CompanyLogo: "",
			Theme: "dark", Lang: "id", MenuVis: defaultVis,
		}, nil
	}

	var menuVis map[string]bool
	_ = json.Unmarshal([]byte(dbRow.MenuVisJSON), &menuVis)

	return model.AppSettings{
		AppName:     dbRow.AppName,
		AppDesc:     dbRow.AppDesc,
		AppEnv:      dbRow.AppEnv,
		CompanyLogo: dbRow.CompanyLogo,
		Theme:       dbRow.Theme,
		Lang:        dbRow.Lang,
		MenuVis:     menuVis,
	}, nil
}

func (r *SettingsRepo) UpdateAppSettings(s model.AppSettings) error {
	// Try to get existing record
	var existing AppSettingsDB
	if err := r.db.First(&existing).Error; err != nil {
		// No existing record, create new one with defaults
		visBytes, _ := json.Marshal(s.MenuVis)
		dbRow := AppSettingsDB{
			AppName:     s.AppName,
			AppDesc:     s.AppDesc,
			AppEnv:      s.AppEnv,
			CompanyLogo: s.CompanyLogo,
			Theme:       s.Theme,
			Lang:        s.Lang,
			MenuVisJSON: string(visBytes),
		}
		if dbRow.AppName == "" {
			dbRow.AppName = "universev"
		}
		if dbRow.AppEnv == "" {
			dbRow.AppEnv = "development"
		}
		if dbRow.Theme == "" {
			dbRow.Theme = "dark"
		}
		if dbRow.Lang == "" {
			dbRow.Lang = "id"
		}
		if dbRow.MenuVisJSON == "" || dbRow.MenuVisJSON == "null" {
			defaultVis, _ := json.Marshal(map[string]bool{
				"display": true, "roster": true, "employees": true, "ftw": true,
				"asset": true, "prestasi": true, "master": true, "users": true, "settings": true,
			})
			dbRow.MenuVisJSON = string(defaultVis)
		}
		return r.db.Create(&dbRow).Error
	}

	// Partial update: only overwrite fields that are non-empty / provided
	updates := map[string]interface{}{}
	if s.AppName != "" {
		updates["app_name"] = s.AppName
	}
	if s.AppDesc != "" {
		updates["app_desc"] = s.AppDesc
	}
	if s.AppEnv != "" {
		updates["app_env"] = s.AppEnv
	}
	if s.CompanyLogo != "" {
		updates["company_logo"] = s.CompanyLogo
	}
	if s.Theme != "" {
		updates["theme"] = s.Theme
	}
	if s.Lang != "" {
		updates["lang"] = s.Lang
	}
	if s.MenuVis != nil {
		visBytes, _ := json.Marshal(s.MenuVis)
		updates["menu_vis_json"] = string(visBytes)
	}

	if len(updates) == 0 {
		return nil
	}
	return r.db.Model(&AppSettingsDB{}).Where("id = ?", existing.ID).Updates(updates).Error
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
