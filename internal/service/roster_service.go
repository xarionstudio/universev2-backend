package service

import (
	"universev/internal/model"
	"universev/internal/repository"
)

type RosterService struct {
	repo *repository.RosterRepo
}

func NewRosterService(repo *repository.RosterRepo) *RosterService {
	return &RosterService{repo: repo}
}

func (s *RosterService) GetRosterByID(id string) (*model.RosterMeta, error) {
	return s.repo.GetRosterByID(id)
}

func (s *RosterService) GetExportRosterData(fileId, deptFilter string) ([]model.RosterExportRow, error) {
	return s.repo.GetExportRosterData(fileId, deptFilter)
}

func (s *RosterService) GetRosters(dept string) ([]model.RosterMeta, error) {
	return s.repo.GetRosters(dept)
}

func (s *RosterService) GetRevisions(status string) ([]model.RosterRevision, error) {
	return s.repo.GetRevisions(status)
}

func (s *RosterService) CreateRoster(meta *model.RosterMeta) error {
	return s.repo.CreateRoster(meta)
}

func (s *RosterService) CreateRevision(rev *model.RosterRevision) error {
	return s.repo.CreateRevision(rev)
}

func (s *RosterService) DeleteRevision(id int) error {
	return s.repo.DeleteRevision(id)
}

func (s *RosterService) ApproveRevision(id int, byId, byEn string) error {
	return s.repo.ApproveRevision(id, byId, byEn)
}

func (s *RosterService) RejectRevision(id int, byId, byEn string) error {
	return s.repo.RejectRevision(id, byId, byEn)
}

func (s *RosterService) GetAttendance(date string) ([]model.AttendanceRow, error) {
	return s.repo.GetAttendance(date)
}

func (s *RosterService) GetSchedulesByFile(fileId string) ([]model.RosterSchedule, error) {
	return s.repo.GetSchedulesByFile(fileId)
}

func (s *RosterService) GetRosterDetail(key string) (map[string]interface{}, error) {
	result, err := s.repo.GetRosterDetail(key)
	if err != nil {
		return nil, err
	}
	return result, nil
}
