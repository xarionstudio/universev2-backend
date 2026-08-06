package handler

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"universev/internal/dto"
	"universev/internal/repository"
	"universev/pkg/response"
)

type AttendanceHandler struct {
	repo *repository.AttendanceRepo
}

func NewAttendanceHandler(repo *repository.AttendanceRepo) *AttendanceHandler {
	return &AttendanceHandler{repo: repo}
}

// GetAttendanceToday — GET /api/attendance/today
func (h *AttendanceHandler) GetAttendanceToday(c fiber.Ctx) error {
	today := time.Now().Format("2006-01-02")
	rows, err := h.repo.GetLogsByDate(today)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch today attendance: "+err.Error())
	}

	meta := &response.Meta{Page: 1, Limit: 50, Total: len(rows), TotalPage: 1}
	return response.SuccessWithMeta(c, fiber.StatusOK, "Success fetch today attendance", rows, meta)
}

// GetAttendanceByDate — GET /api/attendance/date?date=YYYY-MM-DD
func (h *AttendanceHandler) GetAttendanceByDate(c fiber.Ctx) error {
	date := c.Query("date")
	if isTrimmedEmpty(date) {
		date = time.Now().Format("2006-01-02")
	}

	rows, err := h.repo.GetLogsByDate(date)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch attendance by date: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Success fetch attendance for date", rows)
}

// GetAttendanceRange — GET /api/attendance/range?from=YYYY-MM-DD&to=YYYY-MM-DD
func (h *AttendanceHandler) GetAttendanceRange(c fiber.Ctx) error {
	from := c.Query("from")
	to := c.Query("to")
	if isTrimmedEmpty(from) || isTrimmedEmpty(to) {
		return response.Error(c, fiber.StatusBadRequest, "Query parameters 'from' and 'to' are required")
	}

	rows, err := h.repo.GetLogsRange(from, to)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch attendance range: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Success fetch attendance range", rows)
}

// RecordCheckIn — POST /api/attendance/checkin
func (h *AttendanceHandler) RecordCheckIn(c fiber.Ctx) error {
	var req dto.CheckInOutRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(req.NIK) {
		return sendValidationError(c, "nik", "NIK is required")
	}

	if isTrimmedEmpty(req.Machine) {
		req.Machine = "FP-01"
	}

	row, err := h.repo.RecordCheckIn(req.NIK, req.Machine)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to record check-in: "+err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Check-in recorded successfully", row)
}

// RecordCheckOut — POST /api/attendance/checkout
func (h *AttendanceHandler) RecordCheckOut(c fiber.Ctx) error {
	var req dto.CheckInOutRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(req.NIK) {
		return sendValidationError(c, "nik", "NIK is required")
	}

	if isTrimmedEmpty(req.Machine) {
		req.Machine = "FP-01"
	}

	row, err := h.repo.RecordCheckOut(req.NIK, req.Machine)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to record check-out: "+err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Check-out recorded successfully", row)
}
