package dto

// LoginRequest represents the login payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest represents the registration payload
type RegisterRequest struct {
	Name     string `json:"name"`
	NIK      string `json:"nik"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Dept     string `json:"dept"`
	Pos      string `json:"pos"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	Token string `json:"token"`
}
