package model

import "time"

type FingerprintDevice struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Code      string     `gorm:"size:50;unique;not null" json:"code"`
	Name      string     `gorm:"size:150;not null" json:"name"`
	IPAddress string     `gorm:"column:ip_address;size:50;not null" json:"ipAddress"`
	Port      int        `gorm:"default:80" json:"port"`
	ComKey    int        `gorm:"column:com_key;default:0" json:"comKey"`
	Location  string     `gorm:"size:150" json:"location"`
	IsOnline  bool       `gorm:"column:is_online;default:true" json:"isOnline"`
	LastSync  *time.Time `gorm:"column:last_sync" json:"lastSync"`
	IsActive  bool       `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

func (FingerprintDevice) TableName() string {
	return "fingerprint_devices"
}
