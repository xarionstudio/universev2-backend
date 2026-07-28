package model

import "time"

// ── Per-category master structs with descriptive JSON field names ──

type MasterEGIType struct {
	ID        int       `json:"id" gorm:"column:id;primaryKey"`
	Code      string    `json:"code" gorm:"column:code;uniqueIndex"`
	Name      string    `json:"name" gorm:"column:name"`
	Active    bool      `json:"active" gorm:"column:is_active"`
	CreatedAt time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"-" gorm:"column:updated_at"`
}

func (MasterEGIType) TableName() string { return "master_egi_types" }

type MasterProduct struct {
	ID        int       `json:"id" gorm:"column:id;primaryKey"`
	Code      string    `json:"code" gorm:"column:code;uniqueIndex"`
	Name      string    `json:"name" gorm:"column:name"`
	Active    bool      `json:"active" gorm:"column:is_active"`
	CreatedAt time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"-" gorm:"column:updated_at"`
}

func (MasterProduct) TableName() string { return "master_products" }

type MasterEqClass struct {
	ID          int       `json:"id" gorm:"column:id;primaryKey"`
	Code        string    `json:"code" gorm:"column:code;uniqueIndex"`
	Name        string    `json:"name" gorm:"column:name"`
	Description string    `json:"description" gorm:"column:description"`
	Active      bool      `json:"active" gorm:"column:is_active"`
	CreatedAt   time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"-" gorm:"column:updated_at"`
}

func (MasterEqClass) TableName() string { return "master_eq_classes" }

type MasterArea struct {
	ID        int       `json:"id" gorm:"column:id;primaryKey"`
	Code      string    `json:"code" gorm:"column:code;uniqueIndex"`
	Name      string    `json:"name" gorm:"column:name"`
	Category  string    `json:"category" gorm:"column:category"`
	Active    bool      `json:"active" gorm:"column:is_active"`
	CreatedAt time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"-" gorm:"column:updated_at"`
}

func (MasterArea) TableName() string { return "master_areas" }

type MasterTempudo struct {
	ID         int       `json:"id" gorm:"column:id;primaryKey"`
	Code       string    `json:"code" gorm:"column:code;uniqueIndex"`
	Name       string    `json:"name" gorm:"column:name"`
	Location   string    `json:"location" gorm:"column:location"`
	PickupType string    `json:"pickupType" gorm:"column:pickup_type"`
	Active     bool      `json:"active" gorm:"column:is_active"`
	CreatedAt  time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"-" gorm:"column:updated_at"`
}

func (MasterTempudo) TableName() string { return "master_tempudo" }

type MasterBus struct {
	ID            int       `json:"id" gorm:"column:id;primaryKey"`
	Code          string    `json:"code" gorm:"column:code;uniqueIndex"`
	Name          string    `json:"name" gorm:"column:name"`
	EGIType       string    `json:"egiType" gorm:"column:egi_type"`
	DepartureTime string    `json:"departureTime" gorm:"column:departure_time"`
	Active        bool      `json:"active" gorm:"column:is_active"`
	CreatedAt     time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"-" gorm:"column:updated_at"`
}

func (MasterBus) TableName() string { return "master_buses" }

type MasterLocationEx struct {
	ID          int       `json:"id" gorm:"column:id;primaryKey"`
	Code        string    `json:"code" gorm:"column:code;uniqueIndex"`
	Name        string    `json:"name" gorm:"column:name"`
	BusCode     string    `json:"busCode" gorm:"column:bus_code"`
	TempudoCode string    `json:"tempudoCode" gorm:"column:tempudo_code"`
	Active      bool      `json:"active" gorm:"column:is_active"`
	CreatedAt   time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"-" gorm:"column:updated_at"`
}

func (MasterLocationEx) TableName() string { return "master_locations_ex" }

type MasterMess struct {
	ID        int       `json:"id" gorm:"column:id;primaryKey"`
	Code      string    `json:"code" gorm:"column:code;uniqueIndex"`
	Name      string    `json:"name" gorm:"column:name"`
	Block     string    `json:"block" gorm:"column:block"`
	Active    bool      `json:"active" gorm:"column:is_active"`
	CreatedAt time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"-" gorm:"column:updated_at"`
}

func (MasterMess) TableName() string { return "master_mess" }

type MasterRunningText struct {
	ID            int       `json:"id" gorm:"column:id;primaryKey"`
	Code          string    `json:"code" gorm:"column:code;uniqueIndex"`
	Name          string    `json:"name" gorm:"column:name"`
	TargetDisplay string    `json:"targetDisplay" gorm:"column:target_display"`
	TextColor     string    `json:"textColor" gorm:"column:text_color"`
	Active        bool      `json:"active" gorm:"column:is_active"`
	CreatedAt     time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"-" gorm:"column:updated_at"`
}

func (MasterRunningText) TableName() string { return "master_running_texts" }

// ── MdCat constants (kept for reference) ──
type MdCat string

const (
	MdCatEGI      MdCat = "egi"
	MdCatProduct  MdCat = "product"
	MdCatEqClass  MdCat = "eqclass"
	MdCatArea     MdCat = "area"
	MdCatTempudo  MdCat = "tempudo"
	MdCatBus      MdCat = "bus"
	MdCatLokasiEx MdCat = "lokasiex"
	MdCatMess     MdCat = "mess"
	MdCatRunText  MdCat = "runtext"
)
