package dto

// UpdateProfileRequest represents the profile update payload
type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdatePasswordRequest represents the password change payload
type UpdatePasswordRequest struct {
	OldPassword     string `json:"oldPassword"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}
