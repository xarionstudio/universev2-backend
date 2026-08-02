package service

import (
	"time"

	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
)

// DashboardService handles dashboard business logic
type DashboardService struct {
	attRepo    *repository.AttendanceRepo
	ftwRepo    *repository.FTWRepo
	fleetRepo  *repository.FleetRepo
	rosterRepo *repository.RosterRepo
	notifRepo  *repository.NotificationRepo
	empRepo    *repository.EmployeeRepo
}

// NewDashboardService creates a new DashboardService
func NewDashboardService(
	attRepo *repository.AttendanceRepo,
	ftwRepo *repository.FTWRepo,
	fleetRepo *repository.FleetRepo,
	rosterRepo *repository.RosterRepo,
	notifRepo *repository.NotificationRepo,
	empRepo *repository.EmployeeRepo,
) *DashboardService {
	return &DashboardService{
		attRepo:    attRepo,
		ftwRepo:    ftwRepo,
		fleetRepo:  fleetRepo,
		rosterRepo: rosterRepo,
		notifRepo:  notifRepo,
		empRepo:    empRepo,
	}
}

// DashboardSummary represents the aggregated dashboard data
type DashboardSummary struct {
	Attendance struct {
		Total     int `json:"total"`
		Hadir     int `json:"hadir"`
		Terlambat int `json:"terlambat"`
		Belum     int `json:"belum"`
		Off       int `json:"off"`
	} `json:"attendance"`
	FTW struct {
		Total  int `json:"total"`
		Fit    int `json:"fit"`
		Spare  int `json:"spare"`
		Pulang int `json:"pulang"`
		Belum  int `json:"belum"`
	} `json:"ftw"`
	Fleet struct {
		Total     int `json:"total"`
		Ready     int `json:"ready"`
		Breakdown int `json:"breakdown"`
		Standby   int `json:"standby"`
	} `json:"fleet"`
	Roster struct {
		PendingApproval int `json:"pendingApproval"`
	} `json:"roster"`
	Notifications struct {
		Unread int `json:"unread"`
	} `json:"notifications"`
	Employees struct {
		TotalActive int `json:"totalActive"`
	} `json:"employees"`
}

// GetSummary returns aggregated dashboard data
func (s *DashboardService) GetSummary(date string) (*DashboardSummary, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	summary := &DashboardSummary{}

	// Attendance summary
	attRows, err := s.attRepo.GetLogsByDate(date)
	if err == nil {
		for _, row := range attRows {
			summary.Attendance.Total++
			switch row.St {
			case "hadir", "unfit":
				summary.Attendance.Hadir++
			case "terlambat":
				summary.Attendance.Terlambat++
			case "belum":
				summary.Attendance.Belum++
			case "off":
				summary.Attendance.Off++
			}
		}
	}

	// FTW summary
	ftwLogs, err := s.ftwRepo.GetTodayLogs(date)
	if err == nil {
		for _, log := range ftwLogs {
			summary.FTW.Total++
			switch log.St {
			case model.FTWStatusFit:
				summary.FTW.Fit++
			case model.FTWStatusSpare:
				summary.FTW.Spare++
			case model.FTWStatusPulang:
				summary.FTW.Pulang++
			case model.FTWStatusBelum:
				summary.FTW.Belum++
			}
		}
	}

	// Fleet summary
	units, err := s.fleetRepo.GetUnitStatuses()
	if err == nil {
		for _, unit := range units {
			summary.Fleet.Total++
			switch unit.Status {
			case model.UnitStatusReady:
				summary.Fleet.Ready++
			case model.UnitStatusBreakdown:
				summary.Fleet.Breakdown++
			case model.UnitStatusStandby:
				summary.Fleet.Standby++
			}
		}
	}

	// Roster pending approval
	if count, err := s.rosterRepo.CountPendingRevisions(); err == nil {
		summary.Roster.PendingApproval = count
	}

	// Notifications unread
	notifs, err := s.notifRepo.GetByUser("")
	if err == nil {
		for _, n := range notifs {
			if !n.Read {
				summary.Notifications.Unread++
			}
		}
	}

	// Employees total active
	employees, err := s.empRepo.List("", "", "")
	if err == nil {
		for _, emp := range employees {
			if emp.Status == "aktif" {
				summary.Employees.TotalActive++
			}
		}
	}

	return summary, nil
}
