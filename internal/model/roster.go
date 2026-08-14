package model

import "time"

type RosterMeta struct {
	ID        uint      `json:"key" gorm:"column:id;primaryKey;autoIncrement"`
	Label     string    `json:"label" gorm:"column:label"`
	Month     string    `json:"month" gorm:"column:month_period"`
	Dept      string    `json:"dept" gorm:"column:dept"`
	File      string    `json:"file" gorm:"column:filename"`
	Emp       int       `json:"emp" gorm:"column:total_employees"`
	Rows      string    `json:"rows" gorm:"column:total_rows"`
	By        string    `json:"by" gorm:"column:created_by"`
	DateISO   string    `json:"dateISO" gorm:"column:date_iso"`
	Status    string    `json:"status" gorm:"column:status"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (RosterMeta) TableName() string { return "roster_files" }

type RosterRevision struct {
	ID           int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	SubmissionID string    `json:"sid" gorm:"column:submission_id"`
	NIK          string    `json:"nik" gorm:"column:employee_nik"`
	Name         string    `json:"name" gorm:"-"`
	WhatId       string    `json:"whatId" gorm:"column:what_id"`
	WhatEn       string    `json:"whatEn" gorm:"column:what_en"`
	WhenId       string    `json:"whenId" gorm:"column:when_id"`
	WhenEn       string    `json:"whenEn" gorm:"column:when_en"`
	TargetDate   string    `json:"targetDate" gorm:"column:target_date"`
	Status       string    `json:"status" gorm:"column:status"`
	ById         string    `json:"byId,omitempty" gorm:"column:by_id"`
	ByEn         string    `json:"byEn,omitempty" gorm:"column:by_en"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (RosterRevision) TableName() string { return "roster_revisions" }

type RosterSchedule struct {
	ID           uint   `json:"-" gorm:"primaryKey;autoIncrement"`
	RosterFileID uint   `json:"-" gorm:"column:roster_file_id"`
	EmployeeNIK  string `json:"nik" gorm:"column:employee_nik"`
	ScheduleDate string `json:"date" gorm:"column:schedule_date"`
	ShiftCode    string `json:"code" gorm:"column:shift_code"`
}

func (RosterSchedule) TableName() string { return "roster_schedules" }

type AttendanceRow struct {
	ID   uint   `json:"-" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"-"`
	NIK  string `json:"nik" gorm:"column:employee_nik"`
	Dept string `json:"dept" gorm:"-"`
	Code string `json:"code" gorm:"column:shift_code"`
	In   string `json:"in" gorm:"column:check_in"`
	InM  string `json:"inM" gorm:"column:check_in_machine"`
	Out  string `json:"out" gorm:"column:check_out"`
	OutM string `json:"outM" gorm:"column:check_out_machine"`
	St   string `json:"st" gorm:"column:status"`
	Date string `json:"date,omitempty" gorm:"column:attendance_date"`
}

func (AttendanceRow) TableName() string { return "attendance_logs" }

type RosterExportRow struct {
	NIK       string         `json:"nik"`
	Name      string         `json:"name"`
	Dept      string         `json:"dept"`
	Pos       string         `json:"pos"`
	Schedules map[int]string `json:"schedules"`
}

type RosterPreviewRow struct {
	NIK   string         `json:"nik"`
	Name  string         `json:"name"`
	Codes map[int]string `json:"codes"`
}

type RosterError struct {
	Row          string `json:"row"`
	NIK          string `json:"nik"`
	Emp          string `json:"emp"`
	Issue        string `json:"issue"`
	IssueEn      string `json:"issueEn"`
	BadgeVariant string `json:"badgeVariant"` // "danger" | "warning"
	Badge        string `json:"badge"`
}

type RosterValidationResult struct {
	Preview    []RosterPreviewRow `json:"preview"`
	Days       []string           `json:"days"`
	Errors     []RosterError      `json:"errors"`
	ValidCount int                `json:"validCount"`
	DupCount   int                `json:"dupCount"`
	ErrCount   int                `json:"errCount"`
}
