package service

import (
	"universev/internal/model"
	"universev/internal/repository"
)

type NotificationService struct {
	repo *repository.NotificationRepo
}

func NewNotificationService(repo *repository.NotificationRepo) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) GetByUser(userID string) ([]model.Notification, error) {
	return s.repo.GetByUser(userID)
}

func (s *NotificationService) MarkRead(id string) error {
	return s.repo.MarkRead(id)
}

func (s *NotificationService) MarkAllRead(userID string) error {
	return s.repo.MarkAllRead(userID)
}
