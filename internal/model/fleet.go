package model

import "time"

type UnitStatus string

const (
	UnitStatusReady     UnitStatus = "ready"
	UnitStatusBreakdown UnitStatus = "breakdown"
	UnitStatusStandby   UnitStatus = "standby"
)

// UnitHist — matches FE type: UnitHist = [string, string, string, UnitStatus]
type UnitHist [4]string

// UnitStatusRow — DB table unit_statuses
type UnitStatusRow struct {
	Code        string     `json:"code" gorm:"column:code;primaryKey"`
	TypeInfo    string     `json:"type" gorm:"column:type_info"`
	Status      UnitStatus `json:"status" gorm:"column:status"`
	Location    string     `json:"loc" gorm:"column:location"`
	UpdatedNote string     `json:"upd" gorm:"column:updated_note"`
	UpdatedAt   time.Time  `json:"-" gorm:"column:updated_at"`
}

func (UnitStatusRow) TableName() string { return "unit_statuses" }

// UnitHistoryRow — DB table unit_status_histories
type UnitHistoryRow struct {
	ID         uint   `json:"-" gorm:"primaryKey;autoIncrement"`
	UnitCode   string `json:"-" gorm:"column:unit_code"`
	HistWhen   string `json:"when" gorm:"column:hist_when"`
	HistWhat   string `json:"what" gorm:"column:hist_what"`
	HistWhy    string `json:"why" gorm:"column:hist_why"`
	HistStatus string `json:"status" gorm:"column:hist_status"`
}

func (UnitHistoryRow) TableName() string { return "unit_status_histories" }

// Unit — API response combining status + history
type Unit struct {
	Code   string     `json:"code"`
	Type   string     `json:"type"`
	Status UnitStatus `json:"status"`
	Loc    string     `json:"loc"`
	Upd    string     `json:"upd"`
	Hist   []UnitHist `json:"hist"`
}

type FleetSetting struct {
	ID     string   `json:"id" gorm:"column:id;primaryKey"`
	Digger string   `json:"digger" gorm:"column:digger_code"`
	Loc    string   `json:"loc" gorm:"column:location"`
	Bus    string   `json:"bus" gorm:"column:bus_code"`
	Active bool     `json:"active" gorm:"column:is_active"`
	Units  []string `json:"units" gorm:"-"`
}

func (FleetSetting) TableName() string { return "fleet_settings" }

type FleetSettingUnit struct {
	FleetSettingID string `gorm:"column:fleet_setting_id;primaryKey"`
	UnitCode       string `gorm:"column:unit_code;primaryKey"`
}

func (FleetSettingUnit) TableName() string { return "fleet_setting_units" }

type FleetAlloc struct {
	ID     string            `json:"id" gorm:"column:id;primaryKey"`
	Date   string            `json:"date" gorm:"column:alloc_date"`
	Shift  string            `json:"shift" gorm:"column:shift"`
	FlID   string            `json:"flId" gorm:"column:fleet_id"`
	Digger string            `json:"digger" gorm:"column:digger_code"`
	Loc    string            `json:"loc" gorm:"column:location"`
	Bus    string            `json:"bus" gorm:"column:bus_code"`
	Units  map[string]string `json:"units" gorm:"-"`
}

func (FleetAlloc) TableName() string { return "fleet_allocations" }

// FleetAllocResponse — matches FE FaAlloc format:
// { "2026-07-23": { "pagi": { "EX7001": "503264133", ... }, "malam": { ... } } }
type FleetAllocResponse map[string]map[string]map[string]string

type UnitDb struct {
	ID        string    `json:"uid" gorm:"column:id;primaryKey"`
	Code      string    `json:"code" gorm:"column:code;uniqueIndex"`
	EGI       string    `json:"egi" gorm:"column:egi"`
	Product   string    `json:"product" gorm:"column:product"`
	Cls       string    `json:"cls" gorm:"column:class_name"`
	Category  string    `json:"cat" gorm:"column:category"`
	Area      string    `json:"area" gorm:"column:work_area"`
	Active    bool      `json:"active" gorm:"column:is_active"`
	Standby   bool      `json:"standby" gorm:"column:is_standby"`
	Breakdown bool      `json:"breakdown" gorm:"column:is_breakdown"`
	Loc       string    `json:"loc" gorm:"column:location"`
	Upd       string    `json:"upd" gorm:"column:upd_date"`
	By        string    `json:"by" gorm:"column:upd_by"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (UnitDb) TableName() string { return "units_db" }
