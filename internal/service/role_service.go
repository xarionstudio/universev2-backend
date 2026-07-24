package service

import (
	"fmt"

	"universev2-backend/internal/model"
	internalpkg "universev2-backend/internal/pkg"
	"universev2-backend/internal/repository"
)

type RoleService struct {
	repo *repository.RoleRepo
}

func NewRoleService(repo *repository.RoleRepo) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) GetRoles() ([]model.Role, error) {
	return s.repo.GetAll()
}

func (s *RoleService) CreateRole(role *model.Role) error {
	if internalpkg.IsTrimmedEmpty(role.Name) {
		return fmt.Errorf("role name is required")
	}
	return s.repo.Create(role)
}

func (s *RoleService) UpdateRole(id string, role *model.Role) error {
	if internalpkg.IsTrimmedEmpty(id) {
		return fmt.Errorf("role ID is required")
	}
	existing, err := s.repo.GetByID(id)
	if err != nil || existing == nil {
		return fmt.Errorf("role not found")
	}
	if existing.IsLocked {
		return fmt.Errorf("system locked roles cannot be modified")
	}
	if internalpkg.IsTrimmedEmpty(role.Name) {
		return fmt.Errorf("role name is required")
	}
	return s.repo.Update(id, role)
}

func (s *RoleService) DeleteRole(id string) error {
	if internalpkg.IsTrimmedEmpty(id) {
		return fmt.Errorf("role ID is required")
	}
	existing, err := s.repo.GetByID(id)
	if err != nil || existing == nil {
		return fmt.Errorf("role not found")
	}
	if existing.IsLocked {
		return fmt.Errorf("system locked roles cannot be deleted")
	}
	userCount, err := s.repo.CountUsersByRoleID(id)
	if err == nil && userCount > 0 {
		return fmt.Errorf("cannot delete role assigned to active users")
	}
	return s.repo.Delete(id)
}
