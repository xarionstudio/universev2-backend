package repository

import (
	"gorm.io/gorm"

	"universev2-backend/internal/model"
)

type EmployeeRepo struct {
	db *gorm.DB
}

func NewEmployeeRepo(db *gorm.DB) *EmployeeRepo {
	return &EmployeeRepo{db: db}
}

func (r *EmployeeRepo) List(dept, status, search string) ([]model.Employee, error) {
	var employees []model.Employee
	q := r.db.Preload("Komp")

	if dept != "" {
		q = q.Where("dept = ?", dept)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if search != "" {
		q = q.Where("name ILIKE ? OR nik ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	err := q.Order("name ASC").Find(&employees).Error
	return employees, err
}

func (r *EmployeeRepo) GetByNIK(nik string) (*model.Employee, error) {
	var emp model.Employee
	err := r.db.Preload("Komp").Where("nik = ?", nik).First(&emp).Error
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *EmployeeRepo) Create(emp *model.Employee) error {
	return r.db.Create(emp).Error
}

func (r *EmployeeRepo) Update(nik string, emp *model.Employee) error {
	return r.db.Model(&model.Employee{}).Where("nik = ?", nik).Updates(emp).Error
}

func (r *EmployeeRepo) Delete(nik string) error {
	return r.db.Where("nik = ?", nik).Delete(&model.Employee{}).Error
}
