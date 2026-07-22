package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/model"
	"universev2-backend/pkg/response"
)

type AttendanceHandler struct{}

func NewAttendanceHandler() *AttendanceHandler {
	return &AttendanceHandler{}
}

// GetAttendanceToday — GET /api/v1/attendance/today
// Returns today's attendance rows matching FE AttRow type
func (h *AttendanceHandler) GetAttendanceToday(c fiber.Ctx) error {
	rows := []model.AttendanceRow{
		{Name: "First Angel Paustine", NIK: "503264133", Dept: "Operation", Code: "D", In: "05:45", InM: "FP-02 · Gate utara",   Out: "17:32", OutM: "FP-02 · Gate utara",   St: "hadir"},
		{Name: "Rahmat Hidayat",       NIK: "503264134", Dept: "SDI",       Code: "D", In: "07:02", InM: "FP-01 · Office",      Out: "",      OutM: "",                    St: "hadir"},
		{Name: "Budi Santoso",         NIK: "503264135", Dept: "HRGA",      Code: "D", In: "",      InM: "",                    Out: "",      OutM: "",                    St: "unfit"},
		{Name: "Siti Nurhaliza",       NIK: "503264136", Dept: "Operation", Code: "D", In: "05:51", InM: "FP-03 · Gate selatan", Out: "",      OutM: "",                    St: "hadir"},
		{Name: "Andi Prasetyo",        NIK: "503264137", Dept: "Plant",     Code: "D", In: "06:31", InM: "FP-04 · Workshop",    Out: "",      OutM: "",                    St: "terlambat"},
		{Name: "Dewi Lestari",         NIK: "503264138", Dept: "SDI",       Code: "D", In: "06:58", InM: "FP-01 · Office",      Out: "",      OutM: "",                    St: "hadir"},
		{Name: "Joko Widodo S.",       NIK: "503264139", Dept: "Operation", Code: "D", In: "",      InM: "",                    Out: "",      OutM: "",                    St: "belum"},
		{Name: "Rina Marlina",         NIK: "503264140", Dept: "HRGA",      Code: "CR",In: "",      InM: "",                    Out: "",      OutM: "",                    St: "off"},
		{Name: "Agus Salim",           NIK: "503264141", Dept: "Plant",     Code: "D", In: "",      InM: "",                    Out: "",      OutM: "",                    St: "unfit"},
		{Name: "Maya Sari",            NIK: "503264142", Dept: "Operation", Code: "N", In: "",      InM: "",                    Out: "",      OutM: "",                    St: "off"},
	}

	meta := &response.Meta{Page: 1, Limit: 50, Total: len(rows), TotalPage: 1}
	return response.SuccessWithMeta(c, fiber.StatusOK, "Success fetch today attendance", rows, meta)
}

// GetAttendanceByDate — GET /api/v1/attendance?date=YYYY-MM-DD
// Returns attendance rows for a specific date
func (h *AttendanceHandler) GetAttendanceByDate(c fiber.Ctx) error {
	date := c.Query("date")
	if date == "" {
		return response.Error(c, fiber.StatusBadRequest, "Query parameter 'date' is required")
	}

	// Will be replaced with DB query on attendance_logs WHERE attendance_date = date
	return response.Success(c, fiber.StatusOK, "Success fetch attendance for date", []model.AttendanceRow{})
}

// GetAttendanceRange — GET /api/v1/attendance/range?from=YYYY-MM-DD&to=YYYY-MM-DD
// Used by FE roster/attendance page for multi-day view
func (h *AttendanceHandler) GetAttendanceRange(c fiber.Ctx) error {
	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		return response.Error(c, fiber.StatusBadRequest, "Query parameters 'from' and 'to' are required")
	}

	// Will be replaced with DB query on attendance_logs WHERE attendance_date BETWEEN from AND to
	return response.Success(c, fiber.StatusOK, "Success fetch attendance range", []model.AttendanceRow{})
}

// RecordCheckIn — POST /api/v1/attendance/checkin
// Called by fingerprint device when operator checks in
func (h *AttendanceHandler) RecordCheckIn(c fiber.Ctx) error {
	return response.Success(c, fiber.StatusCreated, "Check-in recorded successfully", nil)
}

// RecordCheckOut — POST /api/v1/attendance/checkout
// Called by fingerprint device when operator checks out
func (h *AttendanceHandler) RecordCheckOut(c fiber.Ctx) error {
	return response.Success(c, fiber.StatusCreated, "Check-out recorded successfully", nil)
}
