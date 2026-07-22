package model

import "time"

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

type MdEntry struct {
	ID        string    `json:"id" gorm:"column:id;primaryKey"`
	Cat       MdCat     `json:"cat" gorm:"column:category_key"`
	Name      string    `json:"name" gorm:"column:name"`
	FieldA    string    `json:"a" gorm:"column:field_a"`
	FieldB    string    `json:"b" gorm:"column:field_b"`
	Active    bool      `json:"active" gorm:"column:is_active"`
	CreatedAt time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"-" gorm:"column:updated_at"`
}

func (MdEntry) TableName() string { return "master_entries" }
