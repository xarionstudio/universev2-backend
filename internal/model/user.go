package model

import "time"

type User struct {
	ID                string     `json:"id" gorm:"column:id;primaryKey"`
	Email             string     `json:"email" gorm:"column:email;uniqueIndex"`
	Name              string     `json:"kar" gorm:"column:name"`
	NIK               *string    `json:"nik" gorm:"column:nik"`
	PasswordHash      string     `json:"-" gorm:"column:password_hash"`
	PasswordSalt      string     `json:"-" gorm:"column:password_salt"`
	IsActive          bool       `json:"on" gorm:"column:is_active"`
	Roles             []string   `json:"roles" gorm:"-"`
	PasswordChangedAt *time.Time `json:"pwAt,omitempty" gorm:"column:password_changed_at"`
	CreatedAt         time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt         time.Time  `json:"updatedAt" gorm:"column:updated_at"`
}

func (User) TableName() string { return "users" }

type UserRole struct {
	UserID string `json:"userId" gorm:"column:user_id;primaryKey"`
	RoleID string `json:"roleId" gorm:"column:role_id;primaryKey"`
}

func (UserRole) TableName() string { return "user_roles" }

type Role struct {
	ID          string            `json:"id" gorm:"column:id;primaryKey"`
	Name        string            `json:"name" gorm:"column:name"`
	Description string            `json:"desc" gorm:"column:description"`
	IsLocked    bool              `json:"locked" gorm:"column:is_locked"`
	Perms       map[string]string `json:"perms" gorm:"-"`
	CreatedAt   time.Time         `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time         `json:"updatedAt" gorm:"column:updated_at"`
}

func (Role) TableName() string { return "roles" }

type RolePermission struct {
	RoleID          string `json:"roleId" gorm:"column:role_id;primaryKey"`
	ModuleName      string `json:"module" gorm:"column:module_name;primaryKey"`
	PermissionLevel string `json:"perm" gorm:"column:permission_level"`
}

func (RolePermission) TableName() string { return "role_permissions" }
