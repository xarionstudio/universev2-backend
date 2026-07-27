package dto

// CreateEmployeeRequest represents the employee creation payload
type CreateEmployeeRequest struct {
	NIK       string `json:"nik"`
	Name      string `json:"name"`
	Dept      string `json:"dept"`
	Pos       string `json:"pos"`
	Simper    string `json:"simper"`
	SimperExp string `json:"simperExp"`
	Status    string `json:"status"`
	Company   string `json:"company"`
	Equip     string `json:"equip"`
	HP        string `json:"hp"`
	Join      string `json:"join"`
	Exp       string `json:"exp"`
	License   string `json:"license"`
	MCU       string `json:"mcu"`
	Medis     string `json:"medis"`
	Blood     string `json:"blood"`
	BPJS      string `json:"bpjs"`
	Mess      string `json:"mess"`
	Kamar     string `json:"kamar"`
	Emergency string `json:"emg"`
	Foto      string `json:"foto"`
}

// UpdateEmployeeRequest represents the employee update payload
type UpdateEmployeeRequest struct {
	Name      string `json:"name"`
	Dept      string `json:"dept"`
	Pos       string `json:"pos"`
	Simper    string `json:"simper"`
	SimperExp string `json:"simperExp"`
	Status    string `json:"status"`
	Company   string `json:"company"`
	Equip     string `json:"equip"`
	HP        string `json:"hp"`
	Join      string `json:"join"`
	Exp       string `json:"exp"`
	License   string `json:"license"`
	MCU       string `json:"mcu"`
	Medis     string `json:"medis"`
	Blood     string `json:"blood"`
	BPJS      string `json:"bpjs"`
	Mess      string `json:"mess"`
	Kamar     string `json:"kamar"`
	Emergency string `json:"emg"`
	Foto      string `json:"foto"`
}

// CompetencyRequest represents a single competency entry
type CompetencyRequest struct {
	Type     string `json:"type"`
	Category string `json:"category"`
	Expiry   string `json:"expiry"`
}
