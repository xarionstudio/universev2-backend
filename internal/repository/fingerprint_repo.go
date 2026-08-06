package repository

import (
	"time"

	"gorm.io/gorm"

	"universev/internal/model"
)

type FingerprintRepo struct {
	db *gorm.DB
}

func NewFingerprintRepo(db *gorm.DB) *FingerprintRepo {
	return &FingerprintRepo{db: db}
}

func (r *FingerprintRepo) GetAllDevices() ([]model.FingerprintDevice, error) {
	var devices []model.FingerprintDevice
	err := r.db.Order("id ASC").Find(&devices).Error
	return devices, err
}

func (r *FingerprintRepo) GetActiveDevices() ([]model.FingerprintDevice, error) {
	var devices []model.FingerprintDevice
	err := r.db.Where("is_active = ?", true).Order("id ASC").Find(&devices).Error
	return devices, err
}

func (r *FingerprintRepo) GetDeviceByID(id uint) (*model.FingerprintDevice, error) {
	var device model.FingerprintDevice
	err := r.db.First(&device, id).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *FingerprintRepo) CreateDevice(device *model.FingerprintDevice) error {
	return r.db.Create(device).Error
}

func (r *FingerprintRepo) UpdateDevice(device *model.FingerprintDevice) error {
	device.UpdatedAt = time.Now()
	return r.db.Save(device).Error
}

func (r *FingerprintRepo) DeleteDevice(id uint) error {
	return r.db.Delete(&model.FingerprintDevice{}, id).Error
}

func (r *FingerprintRepo) UpdateSyncStatus(id uint, isOnline bool) error {
	now := time.Now()
	return r.db.Model(&model.FingerprintDevice{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_online":  isOnline,
			"last_sync":  now,
			"updated_at": now,
		}).Error
}
