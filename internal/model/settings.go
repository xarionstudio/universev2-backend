package model

import "time"

type AudioSchedule struct {
	ID       uint     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Title    string   `json:"title" gorm:"column:title"`
	When     string   `json:"when" gorm:"column:trigger_time"`
	Freq     string   `json:"freq" gorm:"column:frequency"`
	File     string   `json:"file" gorm:"column:file_name"`
	Active   bool     `json:"active" gorm:"column:is_active"`
	Displays []string `json:"displays" gorm:"-"`
}

func (AudioSchedule) TableName() string { return "audio_schedules" }

type AudioScheduleDisplay struct {
	AudioID     uint   `gorm:"column:audio_id;primaryKey"`
	DisplayKind string `gorm:"column:display_kind;primaryKey"`
}

func (AudioScheduleDisplay) TableName() string { return "audio_schedule_displays" }

type DisplayDevice struct {
	ID        uint      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Code      string    `json:"code" gorm:"column:code;uniqueIndex"`
	Name      string    `json:"name" gorm:"column:name"`
	Loc       string    `json:"loc" gorm:"column:location"`
	Content   string    `json:"content" gorm:"column:content_kind"`
	FleetID   *uint     `json:"fleetId,omitempty" gorm:"column:fleet_id"`
	RunText   string    `json:"runtext" gorm:"column:running_text"`
	Online    bool      `json:"online" gorm:"column:is_online"`
	Heartbeat string    `json:"hb" gorm:"column:last_heartbeat"`
	Active    bool      `json:"active" gorm:"column:is_active"`
	CreatedAt time.Time `json:"-" gorm:"column:created_at"`
}

func (DisplayDevice) TableName() string { return "display_devices" }

type AppSettings struct {
	AppName     string          `json:"appName"`
	AppEnv      string          `json:"appEnv"`
	CompanyLogo string          `json:"companyLogo"`
	Theme       string          `json:"theme"`
	Lang        string          `json:"lang"`
	MenuVis     map[string]bool `json:"menuVis"`
}
