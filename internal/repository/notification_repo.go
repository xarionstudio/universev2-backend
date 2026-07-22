package repository

import (
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
	var notifs []model.Notification
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifs).Error
	return notifs, err
}

func (r *NotificationRepo) MarkRead(id string) error {
	return r.db.Model(&model.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *NotificationRepo) MarkAllRead(userID string) error {
	return r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = false", userID).Update("is_read", true).Error
}
