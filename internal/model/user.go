package model

import "time"

type User struct {
	ID        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	NIK       string    `json:"nik,omitempty" db:"nik"`
	FullName  string    `json:"fullName" db:"full_name"`
	RoleID    string    `json:"roleId" db:"role_id"`
	RoleName  string    `json:"roleName,omitempty" db:"-"`
	IsActive  bool      `json:"isActive" db:"is_active"`
	Avatar    string    `json:"avatar,omitempty" db:"avatar"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}