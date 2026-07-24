package dto

// CheckInOutRequest represents the check-in/check-out payload
type CheckInOutRequest struct {
	NIK     string `json:"nik"`
	Machine string `json:"machine"`
}
