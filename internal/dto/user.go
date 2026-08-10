package dto

// CreateUserRequest represents the user creation payload
type CreateUserRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	NIK      string   `json:"nik"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

// UpdateUserRequest represents the user update payload
type UpdateUserRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	NIK      string   `json:"nik"`
	Password string   `json:"password,omitempty"`
	Roles    []string `json:"roles"`
}

// CreateRoleRequest represents the role creation payload
type CreateRoleRequest struct {
	Name  string            `json:"name"`
	Perms map[string]string `json:"perms"`
}

// UpdateRoleRequest represents the role update payload
type UpdateRoleRequest struct {
	Name  string            `json:"name"`
	Perms map[string]string `json:"perms"`
}
