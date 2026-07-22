package model

import "time"

type NotifTone string

const (
	ToneInfo    NotifTone = "info"
	ToneSuccess NotifTone = "success"
	ToneWarning NotifTone = "warning"
	ToneDanger  NotifTone = "danger"
)

type Notification struct {
	ID        string    `json:"id" gorm:"column:id;primaryKey"`
	UserID    string    `json:"userId,omitempty" gorm:"column:user_id"`
	Tone      NotifTone `json:"tone" gorm:"column:tone"`
	TextID    string    `json:"textId" gorm:"column:text_id"`
	TextEN    string    `json:"textEn" gorm:"column:text_en"`
	TimeID    string    `json:"timeId" gorm:"column:time_id"`
	TimeEN    string    `json:"timeEn" gorm:"column:time_en"`
	Read      bool      `json:"read" gorm:"column:is_read"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (Notification) TableName() string { return "notifications" }
