package service

import (
	"fmt"
	"time"

	"universev2-backend/internal/dto"
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
		ID:           fmt.Sprintf("u-%d", now.UnixNano()),
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

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(id string, req dto.UpdateUserRequest) error {
	if internalpkg.IsTrimmedEmpty(id) {
		return fmt.Errorf("user ID is required")
	}

	existing, err := s.userRepo.GetByID(id)
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

	return s.userRepo.Update(id, existing)
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id string) error {
	if internalpkg.IsTrimmedEmpty(id) {
		return fmt.Errorf("user ID is required")
	}

	existing, err := s.userRepo.GetByID(id)
	if err != nil || existing == nil {
		return fmt.Errorf("user not found")
	}

	return s.userRepo.Delete(id)
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
	user, err := s.userRepo.GetByID(userID)
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
		user, err := s.userRepo.GetByID(userID)
		if err == nil && user != nil {
			user.Name = req.Name
			return s.userRepo.Update(userID, user)
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

	user, err := s.userRepo.GetByID(userID)
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

	return s.userRepo.UpdatePassword(userID, newHash, newSalt)
}
