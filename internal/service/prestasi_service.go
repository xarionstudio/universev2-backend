package service

import (
	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
)

type PrestasiService struct {
	repo *repository.PrestasiRepo
}

func NewPrestasiService(repo *repository.PrestasiRepo) *PrestasiService {
	return &PrestasiService{repo: repo}
}

func (s *PrestasiService) GetLeaderboard(periodDays int) ([]model.PrestasiRecord, error) {
	return s.repo.GetLeaderboard(periodDays)
}

func (s *PrestasiService) GetHistory(nik string, days int) ([]model.PrestasiHistoryEntry, error) {
	return s.repo.GetOperatorHistory(nik, days)
}
