package dto

// CreateEmployeeRequest represents the employee creation payload
type CreateEmployeeRequest struct {
	NIK     string `json:"nik"`
	Name    string `json:"name"`
	Dept    string `json:"dept"`
	Pos     string `json:"pos"`
	Simper  string `json:"simper"`
	Status  string `json:"status"`
	Company string `json:"company"`
	Equip   string `json:"equip"`
	HP      string `json:"hp"`
}

// UpdateEmployeeRequest represents the employee update payload
type UpdateEmployeeRequest struct {
	Name    string `json:"name"`
	Dept    string `json:"dept"`
	Pos     string `json:"pos"`
	Simper  string `json:"simper"`
	Status  string `json:"status"`
	Company string `json:"company"`
	Equip   string `json:"equip"`
	HP      string `json:"hp"`
}

// CompetencyRequest represents a single competency entry
type CompetencyRequest struct {
	Type     string `json:"type"`
	Category string `json:"category"`
	Expiry   string `json:"expiry"`
}
