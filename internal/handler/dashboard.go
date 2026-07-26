package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/service"
	"universev2-backend/pkg/response"
)

type DashboardHandler struct {
	dashSvc *service.DashboardService
}

func NewDashboardHandler(dashSvc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashSvc: dashSvc}
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

// GetDashboardSummary returns aggregated data for dashboard
func (h *DashboardHandler) GetDashboardSummary(c fiber.Ctx) error {
	today := c.Query("date")
	summary, err := h.dashSvc.GetSummary(today)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch dashboard summary: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success", summary)
}
