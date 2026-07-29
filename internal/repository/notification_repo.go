package repository

import (
	"strconv"

	"gorm.io/gorm"

	"universev2-backend/internal/model"
)

type NotificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepo(db *gorm.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) GetByUser(userID string) ([]model.Notification, error) {
	uid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return nil, err
	}
	var notifs []model.Notification
	err = r.db.Where("user_id = ?", uint(uid)).Order("created_at DESC").Find(&notifs).Error
	return notifs, err
}

func (r *NotificationRepo) MarkRead(id string) error {
	nid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	return r.db.Model(&model.Notification{}).Where("id = ?", uint(nid)).Update("is_read", true).Error
}

func (r *NotificationRepo) MarkAllRead(userID string) error {
	uid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return err
	}
	return r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = false", uint(uid)).Update("is_read", true).Error
}
