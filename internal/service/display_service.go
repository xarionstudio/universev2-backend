package service

import (
	"fmt"
	"time"

	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
)

// DisplayService handles display TV business logic
type DisplayService struct {
	attRepo      *repository.AttendanceRepo
	ftwRepo      *repository.FTWRepo
	fleetRepo    *repository.FleetRepo
	empRepo      *repository.EmployeeRepo
	settingsRepo *repository.SettingsRepo
}

// NewDisplayService creates a new DisplayService
func NewDisplayService(
	attRepo *repository.AttendanceRepo,
	ftwRepo *repository.FTWRepo,
	fleetRepo *repository.FleetRepo,
	empRepo *repository.EmployeeRepo,
	settingsRepo *repository.SettingsRepo,
) *DisplayService {
	return &DisplayService{
		attRepo:      attRepo,
		ftwRepo:      ftwRepo,
		fleetRepo:    fleetRepo,
		empRepo:      empRepo,
		settingsRepo: settingsRepo,
	}
}

// DisplayAttRow represents attendance data for TV display
type DisplayAttRow struct {
	NIK   string `json:"nik"`
	Name  string `json:"name"`
	Pos   string `json:"pos"`
	Dept  string `json:"dept"`
	Shift string `json:"shift"`
	In    string `json:"in"`
	Out   string `json:"out"`
	InM   string `json:"inM"`
	OutM  string `json:"outM"`
	St    string `json:"st"`
	Tone  string `json:"tone"`
	Label string `json:"label"`
}

// GetDisplayAttendance returns attendance data for TV display
func (s *DisplayService) GetDisplayAttendance() ([]DisplayAttRow, error) {
	rows, err := s.attRepo.GetLogsByDate("")
	if err != nil {
		return nil, err
	}

	attTone := map[string]struct{ tone, label string }{
		"belum":     {"danger", "Belum absen"},
		"terlambat": {"warning", "Terlambat"},
		"hadir":     {"success", "Hadir"},
		"unfit":     {"success", "Hadir"},
		"off":       {"neutral", "Off"},
	}

	result := make([]DisplayAttRow, 0)
	for _, r := range rows {
		emp, _ := s.empRepo.GetByNIK(r.NIK)
		name := r.NIK
		pos := ""
		dept := r.Dept
		if emp != nil {
			name = emp.Name
			pos = emp.Pos
			dept = emp.Dept
		}

		tone := "neutral"
		label := r.St
		if t, ok := attTone[r.St]; ok {
			tone = t.tone
			label = t.label
		}

		result = append(result, DisplayAttRow{
			NIK:   r.NIK,
			Name:  name,
			Pos:   pos,
			Dept:  dept,
			Shift: r.Code,
			In:    r.In,
			Out:   r.Out,
			InM:   r.InM,
			OutM:  r.OutM,
			St:    r.St,
			Tone:  tone,
			Label: label,
		})
	}

	return result, nil
}

// DisplayFtwRow represents FTW data for TV display
type DisplayFtwRow struct {
	NIK   string `json:"nik"`
	Name  string `json:"name"`
	Pos   string `json:"pos"`
	Dept  string `json:"dept"`
	Sleep string `json:"sleep"`
	Shift string `json:"shift"`
	St    string `json:"st"`
	Tone  string `json:"tone"`
	Label string `json:"label"`
	Note  string `json:"note"`
}

// GetDisplayFTW returns fit-to-work data for TV display
func (s *DisplayService) GetDisplayFTW() ([]DisplayFtwRow, error) {
	today := time.Now().Format("2006-01-02")
	logs, err := s.ftwRepo.GetTodayLogs(today)
	if err != nil {
		return nil, err
	}

	ftwTone := map[string]struct{ tone, label string }{
		"pulang": {"danger", "Dipulangkan"},
		"spare":  {"warning", "Spare"},
		"belum":  {"warning", "Belum lapor"},
		"fit":    {"success", "Fit"},
	}

	result := make([]DisplayFtwRow, 0)
	for _, l := range logs {
		emp, _ := s.empRepo.GetByNIK(l.NIK)
		name := l.NIK
		pos := ""
		dept := l.Dept
		if emp != nil {
			name = emp.Name
			pos = emp.Pos
			dept = emp.Dept
		}

		tone := "neutral"
		label := string(l.St)
		if t, ok := ftwTone[string(l.St)]; ok {
			tone = t.tone
			label = t.label
		}

		note := ""
		switch l.St {
		case "pulang":
			note = "Tidur < 4 jam — dipulangkan, butuh penggantian"
		case "spare":
			note = "Spare — istirahat sebelum boleh bekerja"
		case "belum":
			note = "Belum mengirim log — hubungi sebelum shift"
		case "fit":
			note = "Lapor " + l.SendTime
		}

		result = append(result, DisplayFtwRow{
			NIK:   l.NIK,
			Name:  name,
			Pos:   pos,
			Dept:  dept,
			Sleep: l.Sleep,
			Shift: l.Shift,
			St:    string(l.St),
			Tone:  tone,
			Label: label,
			Note:  note,
		})
	}

	return result, nil
}

// FleetUnitCard represents a unit card in fleet display
type FleetUnitCard struct {
	Code     string `json:"code"`
	OpName   string `json:"opName"`
	OpNIK    string `json:"opNik"`
	Tone     string `json:"tone"`
	Label    string `json:"label"`
	IsDigger bool   `json:"isDigger"`
}

// FleetDisplayData represents fleet data for TV display
type FleetDisplayData struct {
	ID     string          `json:"id"`
	Digger string          `json:"digger"`
	Loc    string          `json:"loc"`
	Bus    string          `json:"bus"`
	Units  []FleetUnitCard `json:"units"`
}

// GetDisplayFleet returns fleet data for TV display
func (s *DisplayService) GetDisplayFleet(fleetID string) ([]FleetDisplayData, error) {
	settings, err := s.fleetRepo.GetFleetSettings()
	if err != nil {
		return nil, err
	}

	result := make([]FleetDisplayData, 0)
	for _, fs := range settings {
		if fleetID != "" && fmt.Sprintf("%d", fs.ID) != fleetID {
			continue
		}

		unitCodes := []string{fs.Digger}
		unitCodes = append(unitCodes, fs.Units...)

		cards := make([]FleetUnitCard, 0)
		for _, code := range unitCodes {
			statuses, _ := s.fleetRepo.GetUnitStatuses()
			var status *model.Unit
			for _, st := range statuses {
				if st.Code == code {
					status = &st
					break
				}
			}
			tone := "success"
			label := "Ready"
			if status != nil {
				switch status.Status {
				case model.UnitStatusBreakdown:
					tone = "danger"
					label = "Breakdown"
				case model.UnitStatusStandby:
					tone = "neutral"
					label = "Standby"
				}
			}

			cards = append(cards, FleetUnitCard{
				Code:     code,
				OpName:   "",
				OpNIK:    "",
				Tone:     tone,
				Label:    label,
				IsDigger: code == fs.Digger,
			})
		}

		result = append(result, FleetDisplayData{
			ID:     fmt.Sprintf("%d", fs.ID),
			Digger: fs.Digger,
			Loc:    fs.Loc,
			Bus:    fs.Bus,
			Units:  cards,
		})
	}

	return result, nil
}

// DisplayFpDevice represents fingerprint device for TV display
type DisplayFpDevice struct {
	ID     string `json:"id"`
	Loc    string `json:"loc"`
	Online bool   `json:"online"`
	Meta   string `json:"meta"`
}

// GetDisplayFingerprint returns fingerprint device status for TV display
func (s *DisplayService) GetDisplayFingerprint() ([]DisplayFpDevice, error) {
	displays, err := s.settingsRepo.GetDisplays("finger")
	if err != nil {
		return nil, err
	}

	result := make([]DisplayFpDevice, 0)
	for _, d := range displays {
		result = append(result, DisplayFpDevice{
			ID:     fmt.Sprintf("%d", d.ID),
			Loc:    d.Loc,
			Online: d.Online,
			Meta:   d.Heartbeat,
		})
	}

	return result, nil
}
