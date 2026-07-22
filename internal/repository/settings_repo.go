package repository

import (
	"encoding/json"

	"gorm.io/gorm"

	"universev2-backend/internal/model"
)

type SettingsRepo struct {
	db *gorm.DB
}

func NewSettingsRepo(db *gorm.DB) *SettingsRepo {
	return &SettingsRepo{db: db}
}

type AppSettingsDB struct {
	ID          string `gorm:"column:id;primaryKey"`
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
	err := r.db.Where("id = ?", "default").First(&dbRow).Error
	if err != nil {
		defaultVis := map[string]bool{
			"display": true, "roster": true, "employees": true, "ftw": true,
			"asset": true, "prestasi": true, "master": true, "users": true, "settings": true,
		}
		return model.AppSettings{
			AppName: "universev2-backend", AppEnv: "development", CompanyLogo: "",
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
	dbRow := AppSettingsDB{
		ID:          "default",
		AppName:     s.AppName,
		AppEnv:      s.AppEnv,
		CompanyLogo: s.CompanyLogo,
		Theme:       s.Theme,
		Lang:        s.Lang,
		MenuVisJSON: string(visBytes),
	}

	return r.db.Save(&dbRow).Error
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
	if err := r.db.Model(&model.AudioSchedule{}).Where("id = ?", id).Updates(a).Error; err != nil {
		return err
	}
	r.db.Where("audio_id = ?", id).Delete(&model.AudioScheduleDisplay{})
	for _, kind := range a.Displays {
		r.db.Create(&model.AudioScheduleDisplay{AudioID: id, DisplayKind: kind})
	}
	return nil
}

func (r *SettingsRepo) DeleteAudioSchedule(id string) error {
	r.db.Where("audio_id = ?", id).Delete(&model.AudioScheduleDisplay{})
	return r.db.Where("id = ?", id).Delete(&model.AudioSchedule{}).Error
}

// Display Devices

func (r *SettingsRepo) GetDisplays(kind string) ([]model.DisplayDevice, error) {
	var displays []model.DisplayDevice
	q := r.db.Model(&model.DisplayDevice{})
	if kind != "" {
		q = q.Where("content_kind = ?", kind)
	}
	err := q.Order("id ASC").Find(&displays).Error
	return displays, err
}

func (r *SettingsRepo) CreateDisplay(d *model.DisplayDevice) error {
	return r.db.Create(d).Error
}

func (r *SettingsRepo) UpdateDisplay(id string, d *model.DisplayDevice) error {
	return r.db.Model(&model.DisplayDevice{}).Where("id = ?", id).Updates(d).Error
}

func (r *SettingsRepo) DeleteDisplay(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.DisplayDevice{}).Error
}

func (r *SettingsRepo) UpdateHeartbeat(id string, hb string) error {
	return r.db.Model(&model.DisplayDevice{}).Where("id = ?", id).
		Updates(map[string]interface{}{"is_online": true, "last_heartbeat": hb}).Error
}
