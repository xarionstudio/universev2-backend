package repository

import (
	"gorm.io/gorm"

	"universev/internal/model"
	"universev/pkg/filter"
	"universev/pkg/pagination"
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
	return r.db.Model(&model.Employee{}).Where("nik = ?", nik).
		Select("name", "dept", "pos", "simper", "simper_exp", "status", "company", "equip_type",
			"join_date", "exp_date", "license_type", "mcu_status", "med_history", "blood_type",
			"bpjs_no", "mess_name", "room_no", "phone", "emergency_contact", "photo_url").
		Updates(emp).Error
}

func (r *EmployeeRepo) Delete(nik string) error {
	return r.db.Where("nik = ?", nik).Delete(&model.Employee{}).Error
}

func (r *EmployeeRepo) GetCompetencies(nik string) ([]model.Competency, error) {
	var comps []model.Competency
	err := r.db.Joins("JOIN employees ON employees.id = employee_competencies.employee_id").
		Where("employees.nik = ?", nik).Find(&comps).Error
	return comps, err
}

func (r *EmployeeRepo) UpdateCompetencies(nik string, comps []model.Competency) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var emp model.Employee
		if err := tx.Where("nik = ?", nik).First(&emp).Error; err != nil {
			return err
		}
		if err := tx.Where("employee_id = ?", emp.ID).Delete(&model.Competency{}).Error; err != nil {
			return err
		}
		for i := range comps {
			comps[i].EmployeeID = emp.ID
		}
		if len(comps) > 0 {
			if err := tx.Create(&comps).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *EmployeeRepo) UpdatePhoto(nik string, photoURL string) error {
	return r.db.Model(&model.Employee{}).Where("nik = ?", nik).Update("photo_url", photoURL).Error
}

func (r *EmployeeRepo) ListPaginated(f filter.Params, p pagination.Params) ([]model.Employee, int64, error) {
	var employees []model.Employee
	var total int64

	q := r.db.Model(&model.Employee{}).Preload("Komp")
	q = filter.Apply(q, f, filter.Options{
		SearchColumns: []string{"employees.name", "employees.nik"},
		DateColumn:    "employees.join_date",
		StatusColumn:  "employees.status",
		DeptColumn:    "employees.dept",
	})

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := q.Order("name ASC").Limit(p.PerPage).Offset(p.Offset()).Find(&employees).Error
	return employees, total, err
}
