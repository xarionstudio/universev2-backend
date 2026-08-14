package dto

// SubmitFTWLogRequest represents the FTW log submission payload
type SubmitFTWLogRequest struct {
	NIK      string `json:"nik"`
	Shift    string `json:"shift"`
	SleepMin *int   `json:"sleepMin"`
	Sleep    string `json:"sleep"`
	SendTime string `json:"sendTime"`
	Date     string `json:"date"`
}
