package service

import (
	"fmt"
	"strconv"
	"strings"

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
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid role ID")
	}
	existing, err := s.repo.GetByID(uint(uid))
	if err != nil || existing == nil {
		return fmt.Errorf("role not found")
	}
	if existing.IsLocked {
		return fmt.Errorf("system locked roles cannot be modified")
	}
	if internalpkg.IsTrimmedEmpty(role.Name) {
		return fmt.Errorf("role name is required")
	}
	return s.repo.Update(uint(uid), role)
}

func (s *RoleService) DeleteRole(id string) error {
	if internalpkg.IsTrimmedEmpty(id) {
		return fmt.Errorf("role ID is required")
	}
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid role ID")
	}
	existing, err := s.repo.GetByID(uint(uid))
	if err != nil || existing == nil {
		return fmt.Errorf("role not found")
	}
	if existing.IsLocked {
		return fmt.Errorf("system locked roles cannot be deleted")
	}
	userCount, err := s.repo.CountUsersByRoleID(uint(uid))
	if err == nil && userCount > 0 {
		return fmt.Errorf("cannot delete role assigned to active users")
	}
	return s.repo.Delete(uint(uid))
}

var umModules = []string{"dashboard", "display", "employees", "ftw", "roster", "asset", "prestasi", "master", "users", "settings"}

// ExportRolesCSV generates CSV matching FE format: role;deskripsi;employees;ftw;roster;asset;prestasi;master;users;settings
func (s *RoleService) ExportRolesCSV() ([]byte, error) {
	roles, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("role;deskripsi;" + strings.Join(umModules, ";") + "\n")

	for _, r := range roles {
		permsStr := make([]string, len(umModules))
		for i, m := range umModules {
			p := "none"
			if r.Perms != nil && r.Perms[m] != "" {
				p = r.Perms[m]
			}
			permsStr[i] = p
		}
		sb.WriteString(fmt.Sprintf("%s;%s;%s\n", r.Name, r.Description, strings.Join(permsStr, ";")))
	}

	return []byte(sb.String()), nil
}
