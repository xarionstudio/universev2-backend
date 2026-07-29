package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"universev2-backend/internal/dto"
	"universev2-backend/internal/export"
	"universev2-backend/internal/model"
	internalpkg "universev2-backend/internal/pkg"
	"universev2-backend/internal/repository"
)

// UserService handles user management business logic
type UserService struct {
	userRepo *repository.UserRepo
	roleRepo *repository.RoleRepo
}

// NewUserService creates a new UserService
func NewUserService(userRepo *repository.UserRepo, roleRepo *repository.RoleRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// GetUsers returns all users
func (s *UserService) GetUsers() ([]model.User, error) {
	return s.userRepo.GetAll()
}

// CreateUser creates a new user
func (s *UserService) CreateUser(req dto.CreateUserRequest) (*model.User, error) {
	if internalpkg.IsTrimmedEmpty(req.Name) {
		return nil, fmt.Errorf("name is required")
	}
	if !internalpkg.IsValidEmail(req.Email) {
		return nil, fmt.Errorf("invalid email format")
	}
	if len(req.Roles) == 0 {
		return nil, fmt.Errorf("at least one role is required")
	}
	if s.userRepo.ExistsByEmail(req.Email) {
		return nil, fmt.Errorf("email is already in use")
	}

	salt := internalpkg.GenerateSalt()
	hash := internalpkg.HashPasswordFE(req.Password, salt)
	now := time.Now()

	nikPtr := req.NIK
	user := &model.User{
		Email:        req.Email,
		Name:         req.Name,
		NIK:          &nikPtr,
		PasswordHash: hash,
		PasswordSalt: salt,
		IsActive:     true,
		Roles:        req.Roles,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// parseUintID converts string ID to uint
func parseUintID(id string) (uint, error) {
	val, err := strconv.ParseUint(id, 10, 64)
	return uint(val), err
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(id string, req dto.UpdateUserRequest) error {
	if internalpkg.IsTrimmedEmpty(id) {
		return fmt.Errorf("user ID is required")
	}

	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	existing, err := s.userRepo.GetByID(uint(uid))
	if err != nil || existing == nil {
		return fmt.Errorf("user not found")
	}

	if internalpkg.IsTrimmedEmpty(req.Name) {
		return fmt.Errorf("name is required")
	}
	if !internalpkg.IsValidEmail(req.Email) {
		return fmt.Errorf("invalid email format")
	}
	if len(req.Roles) == 0 {
		return fmt.Errorf("at least one role is required")
	}

	existing.Name = req.Name
	existing.Email = req.Email
	if req.NIK != "" {
		existing.NIK = &req.NIK
	}
	existing.Roles = req.Roles

	return s.userRepo.Update(uint(uid), existing)
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id string) error {
	if internalpkg.IsTrimmedEmpty(id) {
		return fmt.Errorf("user ID is required")
	}

	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	existing, err := s.userRepo.GetByID(uint(uid))
	if err != nil || existing == nil {
		return fmt.Errorf("user not found")
	}

	return s.userRepo.Delete(uint(uid))
}

// ToggleUserStatus updates is_active flag for user
func (s *UserService) ToggleUserStatus(id string, active bool) error {
	if internalpkg.IsTrimmedEmpty(id) {
		return fmt.Errorf("user ID is required")
	}
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}
	existing, err := s.userRepo.GetByID(uint(uid))
	if err != nil || existing == nil {
		return fmt.Errorf("user not found")
	}
	return s.userRepo.ToggleStatus(uint(uid), active)
}

// ExportUsersCSV generates a CSV file of all users matching FE format: email;karyawan;nik;roles;status
func (s *UserService) ExportUsersCSV() ([]byte, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("email;karyawan;nik;roles;status\n")

	for _, u := range users {
		nikStr := ""
		if u.NIK != nil {
			nikStr = *u.NIK
		}
		statusStr := "nonaktif"
		if u.IsActive {
			statusStr = "aktif"
		}
		rolesStr := strings.Join(u.Roles, ",")

		sb.WriteString(fmt.Sprintf("%s;%s;%s;%s;%s\n",
			u.Email,
			u.Name,
			nikStr,
			rolesStr,
			statusStr,
		))
	}

	return []byte(sb.String()), nil
}

// ImportUsersFromExcel parses Excel file and creates user records
func (s *UserService) ImportUsersFromExcel(data []byte) (imported int, skipped int, err error) {
	rows, err := export.ParseUserExcel(data)
	if err != nil {
		return 0, 0, err
	}
	if len(rows) == 0 {
		return 0, 0, fmt.Errorf("no valid user records found in file")
	}

	for _, r := range rows {
		if !internalpkg.IsValidEmail(r.Email) {
			skipped++
			continue
		}
		if s.userRepo.ExistsByEmail(r.Email) {
			skipped++
			continue
		}

		defaultRole := "3" // Viewer role
		if r.Role != "" {
			defaultRole = r.Role
		}

		nikVal := r.NIK
		var nikPtr *string
		if nikVal != "" {
			nikPtr = &nikVal
		}

		salt := internalpkg.GenerateSalt()
		hash := internalpkg.HashPasswordFE("Password123!", salt) // Default initial password
		now := time.Now()

		u := &model.User{
			Email:        r.Email,
			Name:         r.Name,
			NIK:          nikPtr,
			PasswordHash: hash,
			PasswordSalt: salt,
			IsActive:     true,
			Roles:        []string{defaultRole},
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		if err := s.userRepo.Create(u); err == nil {
			imported++
		} else {
			skipped++
		}
	}

	return imported, skipped, nil
}

// ProfileService handles profile-related business logic
type ProfileService struct {
	userRepo *repository.UserRepo
}

// NewProfileService creates a new ProfileService
func NewProfileService(userRepo *repository.UserRepo) *ProfileService {
	return &ProfileService{userRepo: userRepo}
}

// GetProfile returns a user's profile by ID
func (s *ProfileService) GetProfile(userID string) (*model.User, error) {
	if s.userRepo == nil {
		return nil, fmt.Errorf("user repository unavailable")
	}
	uid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}
	user, err := s.userRepo.GetByID(uint(uid))
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// UpdateProfile updates a user's profile name
func (s *ProfileService) UpdateProfile(userID string, req dto.UpdateProfileRequest) error {
	if internalpkg.IsTrimmedEmpty(req.Name) {
		return fmt.Errorf("name is required")
	}

	if s.userRepo != nil {
		uid, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid user ID")
		}
		user, err := s.userRepo.GetByID(uint(uid))
		if err == nil && user != nil {
			user.Name = req.Name
			return s.userRepo.Update(uint(uid), user)
		}
	}
	return nil
}

// UpdatePassword changes a user's password
func (s *ProfileService) UpdatePassword(userID string, req dto.UpdatePasswordRequest) error {
	if req.OldPassword == "" {
		return fmt.Errorf("old password is required")
	}
	if !internalpkg.IsPasswordStrong(req.NewPassword) {
		return fmt.Errorf("password must be at least 8 characters and contain letters and numbers")
	}
	if req.NewPassword != req.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	if s.userRepo == nil {
		return fmt.Errorf("user repository unavailable")
	}

	uid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	user, err := s.userRepo.GetByID(uint(uid))
	if err != nil || user == nil {
		return fmt.Errorf("user not found")
	}

	// Verify old password
	oldHashFE := internalpkg.HashPasswordFE(req.OldPassword, user.PasswordSalt)
	oldHashLegacy := internalpkg.HashPasswordLegacy(req.OldPassword, user.PasswordSalt)
	if oldHashFE != user.PasswordHash && oldHashLegacy != user.PasswordHash {
		return fmt.Errorf("incorrect old password")
	}

	newSalt := internalpkg.GenerateSalt()
	newHash := internalpkg.HashPasswordFE(req.NewPassword, newSalt)

	return s.userRepo.UpdatePassword(uint(uid), newHash, newSalt)
}
