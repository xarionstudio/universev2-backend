package service

import (
	"fmt"

	"universev2-backend/internal/model"
	internalpkg "universev2-backend/internal/pkg"
	"universev2-backend/internal/repository"
)

type SettingsService struct {
	repo *repository.SettingsRepo
}

func NewSettingsService(repo *repository.SettingsRepo) *SettingsService {
	return &SettingsService{repo: repo}
}

func (s *SettingsService) GetAppSettings() (*model.AppSettings, error) {
	settings, err := s.repo.GetAppSettings()
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (s *SettingsService) UpdateAppSettings(settings model.AppSettings) error {
	if internalpkg.IsTrimmedEmpty(settings.AppName) {
		settings.AppName = "universev2-backend"
	}
	return s.repo.UpdateAppSettings(settings)
}

func (s *SettingsService) GetAudioSchedules() ([]model.AudioSchedule, error) {
	return s.repo.GetAudioSchedules()
}

func (s *SettingsService) CreateAudioSchedule(a *model.AudioSchedule) error {
	if internalpkg.IsTrimmedEmpty(a.Title) {
		return fmt.Errorf("audio schedule title is required")
	}
	if internalpkg.IsTrimmedEmpty(a.When) {
		return fmt.Errorf("trigger time is required")
	}
	return s.repo.CreateAudioSchedule(a)
}

func (s *SettingsService) UpdateAudioSchedule(id string, a *model.AudioSchedule) error {
	if internalpkg.IsTrimmedEmpty(id) {
		return fmt.Errorf("audio schedule ID is required")
	}
	if internalpkg.IsTrimmedEmpty(a.Title) {
		return fmt.Errorf("audio schedule title is required")
	}
	if internalpkg.IsTrimmedEmpty(a.When) {
		return fmt.Errorf("trigger time is required")
	}
	return s.repo.UpdateAudioSchedule(id, a)
}

func (s *SettingsService) DeleteAudioSchedule(id string) error {
	if internalpkg.IsTrimmedEmpty(id) {
		return fmt.Errorf("audio schedule ID is required")
	}
	return s.repo.DeleteAudioSchedule(id)
}

func (s *SettingsService) GetDisplays(kind string) ([]model.DisplayDevice, error) {
	return s.repo.GetDisplays(kind)
}

func (s *SettingsService) CreateDisplay(d *model.DisplayDevice) error {
	if internalpkg.IsTrimmedEmpty(d.Name) {
		return fmt.Errorf("display name is required")
	}
	if internalpkg.IsTrimmedEmpty(d.Loc) {
		return fmt.Errorf("display location is required")
	}
	return s.repo.CreateDisplay(d)
}

func (s *SettingsService) UpdateDisplay(id string, d *model.DisplayDevice) error {
	if internalpkg.IsTrimmedEmpty(id) {
		return fmt.Errorf("display ID is required")
	}
	if internalpkg.IsTrimmedEmpty(d.Name) {
		return fmt.Errorf("display name is required")
	}
	if internalpkg.IsTrimmedEmpty(d.Loc) {
		return fmt.Errorf("display location is required")
	}
	return s.repo.UpdateDisplay(id, d)
}

func (s *SettingsService) DeleteDisplay(id string) error {
	if internalpkg.IsTrimmedEmpty(id) {
		return fmt.Errorf("display ID is required")
	}
	return s.repo.DeleteDisplay(id)
}
