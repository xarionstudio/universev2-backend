package model

import "time"

type Competency struct {
	ID         uint   `json:"-" gorm:"primaryKey;autoIncrement"`
	EmployeeID uint   `json:"-" gorm:"column:employee_id"`
	Class      string `json:"cls" gorm:"column:class_name"`
	Simper     string `json:"simper" gorm:"column:simper_no"`
	Exp        string `json:"exp" gorm:"column:expiry_date"`
}

func (Competency) TableName() string { return "employee_competencies" }

type Employee struct {
	ID        uint         `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	NIK       string       `json:"nik" gorm:"column:nik;uniqueIndex"`
	Name      string       `json:"name" gorm:"column:name"`
	Dept      string       `json:"dept" gorm:"column:dept"`
	Pos       string       `json:"pos" gorm:"column:pos"`
	Simper    string       `json:"simper" gorm:"column:simper"`
	SimperExp string       `json:"simperExp" gorm:"column:simper_exp"`
	Status    string       `json:"status" gorm:"column:status"`
	Company   string       `json:"company" gorm:"column:company"`
	Equip     string       `json:"equip" gorm:"column:equip_type"`
	Join      string       `json:"join" gorm:"column:join_date"`
	Exp       string       `json:"exp" gorm:"column:exp_date"`
	License   string       `json:"license" gorm:"column:license_type"`
	MCU       string       `json:"mcu" gorm:"column:mcu_status"`
	Medis     string       `json:"medis" gorm:"column:med_history"`
	Blood     string       `json:"blood" gorm:"column:blood_type"`
	BPJS      string       `json:"bpjs" gorm:"column:bpjs_no"`
	Mess      string       `json:"mess" gorm:"column:mess_name"`
	Kamar     string       `json:"kamar" gorm:"column:room_no"`
	HP        string       `json:"hp" gorm:"column:phone"`
	Emergency string       `json:"emg" gorm:"column:emergency_contact"`
	Foto      string       `json:"foto,omitempty" gorm:"column:photo_url"`
	Komp      []Competency `json:"komp,omitempty" gorm:"foreignKey:EmployeeID;references:ID"`
	CreatedAt time.Time    `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time    `json:"updatedAt" gorm:"column:updated_at"`
}

func (Employee) TableName() string { return "employees" }
