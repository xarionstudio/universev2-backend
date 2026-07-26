package service

import (
	"fmt"
	"time"

	"universev2-backend/internal/config"
	"universev2-backend/internal/dto"
	"universev2-backend/internal/model"
	internalpkg "universev2-backend/internal/pkg"
	"universev2-backend/internal/repository"
)

// AuthService handles authentication business logic
type AuthService struct {
	cfg      *config.Config
	userRepo *repository.UserRepo
	roleRepo *repository.RoleRepo
}

// NewAuthService creates a new AuthService
func NewAuthService(cfg *config.Config, userRepo *repository.UserRepo, roleRepo *repository.RoleRepo) *AuthService {
	return &AuthService{
		cfg:      cfg,
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// LoginResult contains the login response data
type LoginResult struct {
	Token string            `json:"token"`
	User  *model.User       `json:"user"`
	Perms map[string]string `json:"perms"`
}

// Login authenticates a user and returns JWT token
func (s *AuthService) Login(req dto.LoginRequest) (*LoginResult, error) {
	if s.userRepo == nil {
		return nil, fmt.Errorf("user repository unavailable")
	}

	if internalpkg.IsTrimmedEmpty(req.Email) || req.Password == "" {
		return nil, fmt.Errorf("email and password are required")
	}
	if !internalpkg.IsValidEmail(req.Email) {
		return nil, fmt.Errorf("invalid email format")
	}

	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil || user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}

	// Verify hash with FE format or fallback to legacy format
	computedHash := internalpkg.HashPasswordFE(req.Password, user.PasswordSalt)
	legacyHash := internalpkg.HashPasswordLegacy(req.Password, user.PasswordSalt)

	if computedHash != user.PasswordHash && legacyHash != user.PasswordHash {
		return nil, fmt.Errorf("invalid credentials")
	}

	tokenStr, err := internalpkg.GenerateToken(user.ID, user.Email, user.Roles, s.cfg.JWTSecret, s.cfg.JWTExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}

	var perms map[string]string
	if s.roleRepo != nil {
		perms, _ = s.roleRepo.GetPermissionsForRoles(user.Roles)
	}

	return &LoginResult{
		Token: tokenStr,
		User:  user,
		Perms: perms,
	}, nil
}

// RegisterResult contains the register response data
type RegisterResult struct {
	User *model.User `json:"user"`
}

// Register creates a new user account
func (s *AuthService) Register(req dto.RegisterRequest) (*RegisterResult, error) {
	if internalpkg.IsTrimmedEmpty(req.Name) {
		return nil, fmt.Errorf("name is required")
	}
	if !internalpkg.IsValidNIK(req.NIK) {
		return nil, fmt.Errorf("NIK must be exactly 9 digits")
	}
	if internalpkg.IsTrimmedEmpty(req.Dept) {
		return nil, fmt.Errorf("department is required")
	}
	if internalpkg.IsTrimmedEmpty(req.Pos) {
		return nil, fmt.Errorf("position is required")
	}
	if !internalpkg.IsValidEmail(req.Email) {
		return nil, fmt.Errorf("invalid email format")
	}
	if !internalpkg.IsPasswordStrong(req.Password) {
		return nil, fmt.Errorf("password must be at least 8 characters and contain both letters and numbers")
	}
	if len(req.Password) > 72 {
		return nil, fmt.Errorf("password must not exceed 72 characters")
	}
	if len(req.Name) > 100 {
		return nil, fmt.Errorf("name must not exceed 100 characters")
	}
	if len(req.NIK) > 9 {
		return nil, fmt.Errorf("NIK must not exceed 9 characters")
	}

	if s.userRepo.ExistsByEmail(req.Email) {
		return nil, fmt.Errorf("email is already registered")
	}

	if s.userRepo.ExistsByNIK(req.NIK) {
		return nil, fmt.Errorf("NIK is already registered")
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
		Roles:        []string{"r3"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &RegisterResult{User: user}, nil
}

// RefreshTokenResult contains the refresh token response data
type RefreshTokenResult struct {
	Token string            `json:"token"`
	User  *model.User       `json:"user"`
	Perms map[string]string `json:"perms"`
}

// RefreshToken refreshes an existing JWT token
func (s *AuthService) RefreshToken(tokenStr string) (*RefreshTokenResult, error) {
	claims, err := internalpkg.ParseToken(tokenStr, s.cfg.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired token")
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}
	if !user.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}

	newToken, err := internalpkg.GenerateToken(user.ID, user.Email, user.Roles, s.cfg.JWTSecret, s.cfg.JWTExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new token")
	}

	var perms map[string]string
	if s.roleRepo != nil {
		perms, _ = s.roleRepo.GetPermissionsForRoles(user.Roles)
	}

	return &RefreshTokenResult{
		Token: newToken,
		User:  user,
		Perms: perms,
	}, nil
}
