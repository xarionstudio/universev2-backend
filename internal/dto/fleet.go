package dto

// CreateFleetSettingRequest represents the fleet setting creation payload
type CreateFleetSettingRequest struct {
	Digger string   `json:"digger"`
	Loc    string   `json:"loc"`
	Bus    string   `json:"bus"`
	Units  []string `json:"units"`
}

// UpdateFleetSettingRequest represents the fleet setting update payload
type UpdateFleetSettingRequest struct {
	Digger string   `json:"digger"`
	Loc    string   `json:"loc"`
	Bus    string   `json:"bus"`
	Units  []string `json:"units"`
}

// AutoAllocateRequest represents the auto allocation payload
type AutoAllocateRequest struct {
	Date  string `json:"date"`
	Shift string `json:"shift"`
}

// SaveAllocationRequest represents the manual allocation save payload
type SaveAllocationRequest struct {
	Date  string            `json:"date"`
	Shift string            `json:"shift"`
	Units map[string]string `json:"units"`
}

// UpdateUnitStatusRequest represents the unit status update payload
type UpdateUnitStatusRequest struct {
	Status string `json:"status"`
	Note   string `json:"note"`
}

// ReportBreakdownRequest represents the breakdown report payload
type ReportBreakdownRequest struct {
	Reason string `json:"reason"`
}

// CreateUnitDBRequest represents the unit DB creation payload
type CreateUnitDBRequest struct {
	Code      string `json:"code"`
	EGI       string `json:"egi"`
	Product   string `json:"product"`
	Cls       string `json:"cls"`
	Category  string `json:"cat"`
	Area      string `json:"area"`
	Active    bool   `json:"active"`
	Standby   bool   `json:"standby"`
	Breakdown bool   `json:"breakdown"`
	Loc       string `json:"loc"`
	Upd       string `json:"upd"`
	By        string `json:"by"`
}

// UpdateUnitDBRequest represents the unit DB update payload
type UpdateUnitDBRequest struct {
	Code      string `json:"code"`
	EGI       string `json:"egi"`
	Product   string `json:"product"`
	Cls       string `json:"cls"`
	Category  string `json:"cat"`
	Area      string `json:"area"`
	Active    bool   `json:"active"`
	Standby   bool   `json:"standby"`
	Breakdown bool   `json:"breakdown"`
	Loc       string `json:"loc"`
	Upd       string `json:"upd"`
	By        string `json:"by"`
}
