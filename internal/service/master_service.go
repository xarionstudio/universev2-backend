package service

import (
	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
)

type MasterService struct {
	repo *repository.MasterRepo
}

func NewMasterService(repo *repository.MasterRepo) *MasterService {
	return &MasterService{repo: repo}
}

func (s *MasterService) GetByCategory(cat string) ([]model.MdEntry, error) {
	return s.repo.GetByCategory(cat)
}

func (s *MasterService) CreateEntry(entry *model.MdEntry) error {
	return s.repo.Create(entry)
}

func (s *MasterService) UpdateEntry(id string, entry *model.MdEntry) error {
	return s.repo.Update(id, entry)
}

func (s *MasterService) DeleteEntry(id string) error {
	return s.repo.Delete(id)
}
